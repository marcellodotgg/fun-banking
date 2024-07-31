package service

import (
	"errors"
	"fmt"
	"strconv"

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
	SendMoney(fromAccount domain.Account, recipient domain.Customer, transaction *domain.Transaction) error
	BulkTransfer(customerIDs []string, transaction *domain.Transaction) error
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

func (ts transactionService) Create(transaction *domain.Transaction) error {
	return persistence.DB.Transaction(func(tx *gorm.DB) error {
		accountID := strconv.Itoa(transaction.AccountID)

		var account domain.Account
		if err := ts.accountService.FindByID(accountID, &account); err != nil {
			return err
		}

		if account.Customer.Bank.UserID == *transaction.UserID {
			transaction.Status = domain.TransactionApproved
			account.Balance = utils.SafelyAddDollars(account.Balance, transaction.Amount)
			transaction.Balance = account.Balance

			if err := ts.accountService.UpdateBalance(accountID, &account); err != nil {
				return err
			}
		}

		if transaction.Status == domain.TransactionPending {
			transaction.Balance = utils.SafelyAddDollars(account.Balance, transaction.Amount)
		}

		return persistence.DB.Create(&transaction).Error
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
			account.Balance += transaction.Amount
			return persistence.DB.Select("Balance").Updates(&account).Error
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
