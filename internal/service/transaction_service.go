package service

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/bytebury/fun-banking/internal/domain"
	"github.com/bytebury/fun-banking/internal/infrastructure/persistence"
	"github.com/bytebury/fun-banking/internal/utils"
	"gorm.io/gorm"
)

type Cashflow struct {
	Deposits    float64
	Withdrawals float64
}

type TransactionService interface {
	Create(transaction *domain.Transaction) error
	Update(id, userID, status string) error
	TransferMoney(from domain.Account, to domain.Account, amount float64) error
	SendMoney(fromAccount domain.Account, recipient domain.Customer, transaction *domain.Transaction) error
	BulkTransfer(customerIDs []string, transaction *domain.Transaction) error
	AutoPay(autoPay domain.AutoPay) error
}

type transactionService struct {
	accountService  AccountService
	customerService CustomerService
}

func NewTransactionService() TransactionService {
	return transactionService{
		accountService:  NewAccountService(),
		customerService: NewCustomerService(),
	}
}

func (ts transactionService) AutoPay(autoPay domain.AutoPay) error {
	return persistence.DB.Transaction(func(tx *gorm.DB) error {
		var account domain.Account
		if err := persistence.DB.
			Preload("Customer").
			Preload("Customer.Bank").
			First(&account, "id = ?", autoPay.AccountID).Error; err != nil {
			return err
		}

		transaction := domain.Transaction{
			AccountID:   autoPay.AccountID,
			Amount:      autoPay.Amount,
			Description: autoPay.Description,
			UserID:      &account.Customer.Bank.UserID,
		}

		if err := ts.Create(&transaction); err != nil {
			return err
		}

		nextRunDate := time.Now()

		switch autoPay.Cadence {
		case "day":
			nextRunDate = nextRunDate.AddDate(0, 0, 1)
		case "week":
			nextRunDate = nextRunDate.AddDate(0, 0, 7)
		case "month":
			nextRunDate = nextRunDate.AddDate(0, 1, 0)
		}

		return persistence.DB.Model(&autoPay).Select("NextRunDate").Updates(domain.AutoPay{NextRunDate: nextRunDate}).Error
	})
}

func (ts transactionService) Create(transaction *domain.Transaction) error {
	return persistence.DB.Transaction(func(tx *gorm.DB) error {
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
			}
		}()

		accountID := strconv.Itoa(transaction.AccountID)

		var account domain.Account
		if err := ts.accountService.FindByID(accountID, &account); err != nil {
			return err
		}

		if account.Customer.Bank.UserID == *transaction.UserID {
			transaction.Status = domain.TransactionApproved
			account.Balance = utils.SafelyAddDollars(account.Balance, transaction.Amount)
			transaction.Balance = account.Balance

			if err := tx.Where("id = ?", accountID).Select("balance").Updates(&account).Error; err != nil {
				tx.Rollback()
				return err
			}
		}

		if transaction.Status == domain.TransactionPending {
			transaction.Balance = account.Balance
		}

		if err := tx.Create(&transaction).Error; err != nil {
			tx.Rollback()
			return err
		}

		return nil
	})
}

func (s transactionService) TransferMoney(from domain.Account, to domain.Account, amount float64) error {
	return persistence.DB.Transaction(func(tx *gorm.DB) error {
		if amount <= 0 {
			return errors.New("amount must be greater than 0")
		}

		if from.Balance < amount {
			return errors.New("not enough money")
		}

		if from.ID == to.ID {
			return errors.New("cannot transfer to same account")
		}

		if from.CustomerID != to.CustomerID {
			return errors.New("cannot transfer to other customers accounts")
		}

		fromTransaction := domain.Transaction{}

		from.Balance = utils.SafelySubtractDollars(from.Balance, amount)

		fromTransaction.Amount = amount * -1
		fromTransaction.Balance = from.Balance
		fromTransaction.AccountID = from.ID
		fromTransaction.Description = fmt.Sprintf("Money transfer to %s", to.Name)
		fromTransaction.Status = domain.TransactionApproved

		if err := s.accountService.UpdateBalance(strconv.Itoa(from.ID), &from); err != nil {
			return err
		}

		if err := persistence.DB.Create(&fromTransaction).Error; err != nil {
			return err
		}

		toTransaction := domain.Transaction{}

		to.Balance = utils.SafelyAddDollars(to.Balance, amount)

		toTransaction.Amount = amount
		toTransaction.Balance = to.Balance
		toTransaction.AccountID = to.ID
		toTransaction.Description = fmt.Sprintf("Money transfer from %s", from.Name)
		toTransaction.Status = domain.TransactionApproved

		if err := s.accountService.UpdateBalance(strconv.Itoa(to.ID), &to); err != nil {
			return err
		}

		return persistence.DB.Create(&toTransaction).Error
	})
}

func (s transactionService) SendMoney(fromAccount domain.Account, recipient domain.Customer, transaction *domain.Transaction) error {
	return persistence.DB.Transaction(func(tx *gorm.DB) error {
		if fromAccount.Balance < transaction.Amount {
			return errors.New("not enough money")
		}

		amount := transaction.Amount
		description := transaction.Description

		fromAccount.Balance = utils.SafelySubtractDollars(fromAccount.Balance, amount)

		transaction.Amount = amount * -1
		transaction.Balance = fromAccount.Balance
		transaction.AccountID = fromAccount.ID
		transaction.Description = fmt.Sprintf("Money sent to %s. Note: %s", recipient.FullName(), description)
		transaction.Status = domain.TransactionApproved

		if err := s.accountService.UpdateBalance(strconv.Itoa(fromAccount.ID), &fromAccount); err != nil {
			return err
		}

		if err := persistence.DB.Create(&transaction).Error; err != nil {
			return err
		}

		toAccount, err := recipient.PrimaryAccount()

		if err != nil {
			return err
		}

		secondTransaction := domain.Transaction{}

		toAccount.Balance = utils.SafelyAddDollars(toAccount.Balance, amount)

		secondTransaction.Amount = amount
		secondTransaction.Balance = toAccount.Balance
		secondTransaction.AccountID = toAccount.ID
		secondTransaction.Description = fmt.Sprintf("Money sent from %s. Note: %s", fromAccount.Customer.FullName(), description)
		secondTransaction.Status = domain.TransactionApproved

		if err := s.accountService.UpdateBalance(strconv.Itoa(toAccount.ID), &toAccount); err != nil {
			return err
		}

		return persistence.DB.Create(&secondTransaction).Error
	})
}

func (s transactionService) Update(id, userID, status string) error {
	return persistence.DB.Transaction(func(tx *gorm.DB) error {
		var transaction domain.Transaction
		if err := persistence.DB.First(&transaction, "id = ?", id).Error; err != nil {
			return err
		}

		var account domain.Account
		if err := persistence.DB.First(&account, "id = ?", transaction.AccountID).Error; err != nil {
			return err
		}

		userIDPtr, _ := utils.ConvertToIntPointer(userID)

		transaction.UserID = userIDPtr
		transaction.Status = status
		transaction.Balance = account.Balance

		if err := persistence.DB.Updates(&transaction).Error; err != nil {
			return err
		}

		if transaction.Status == domain.TransactionApproved {
			account.Balance = utils.SafelyAddDollars(account.Balance, transaction.Amount)

			if err := persistence.DB.Select("Balance").Updates(&account).Error; err != nil {
				return err
			}

			transaction.Balance = account.Balance
			return persistence.DB.Select("Balance").Updates(&transaction).Error
		}
		return persistence.DB.Select("Balance").First(&account).Error
	})
}

func (s transactionService) BulkTransfer(customerIDs []string, transaction *domain.Transaction) error {
	return persistence.DB.Transaction(func(tx *gorm.DB) error {
		for _, customerID := range customerIDs {
			var customer domain.Customer
			if err := s.customerService.FindByID(customerID, &customer); err != nil {
				return err
			}

			primaryAccount, err := customer.PrimaryAccount()

			if err != nil {
				return err
			}

			newTransaction := domain.Transaction{
				AccountID:   primaryAccount.ID,
				Amount:      transaction.Amount,
				Description: transaction.Description,
				UserID:      transaction.UserID,
			}
			if err := s.Create(&newTransaction); err != nil {
				return err
			}
		}
		return nil
	})
}
