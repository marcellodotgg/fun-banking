package service

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/bytebury/fun-banking/internal/domain"
	"github.com/bytebury/fun-banking/internal/infrastructure/pagination"
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
	FindAllByAccount(accountID string, transactions *[]domain.Transaction, pagingInfo pagination.PagingInfo[domain.Transaction]) error
	CountAllByAccount(accountID string, count *int64) error
	CashflowByAccount(accountID string, cashflow *Cashflow) error
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

		transaction.Balance = account.Balance

		if account.Customer.Bank.UserID == *transaction.UserID {
			transaction.Status = domain.TransactionApproved
			account.Balance += transaction.Amount

			if err := ts.accountService.UpdateBalance(accountID, &account); err != nil {
				return err
			}
		}

		return persistence.DB.Create(&transaction).Error
	})
}

func (s transactionService) SendMoney(fromAccount domain.Account, recipient domain.Customer, transaction *domain.Transaction) error {
	return persistence.DB.Transaction(func(tx *gorm.DB) error {
		if fromAccount.Balance < transaction.Amount {
			return errors.New("you do not have enough money")
		}

		amount := transaction.Amount
		description := transaction.Description

		transaction.Amount = amount * -1
		transaction.Balance = fromAccount.Balance
		transaction.AccountID = fromAccount.ID
		transaction.Description = fmt.Sprintf("Money sent to %s. Note: %s", recipient.FullName(), description)
		transaction.Status = domain.TransactionApproved

		fromAccount.Balance -= amount

		if err := s.accountService.UpdateBalance(strconv.Itoa(fromAccount.ID), &fromAccount); err != nil {
			return err
		}

		if err := persistence.DB.Create(&transaction).Error; err != nil {
			return err
		}

		toAccount := recipient.Accounts[0]

		secondTransaction := domain.Transaction{}

		secondTransaction.Amount = amount
		secondTransaction.Balance = toAccount.Balance
		secondTransaction.AccountID = toAccount.ID
		secondTransaction.Description = fmt.Sprintf("Money sent from %s. Note: %s", fromAccount.Customer.FullName(), description)
		secondTransaction.Status = domain.TransactionApproved

		toAccount.Balance += amount

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

		account.Balance += transaction.Amount

		return persistence.DB.Select("Balance").Updates(&account).Error
	})
}

func (ts transactionService) FindAllByAccount(accountID string, transactions *[]domain.Transaction, pagingInfo pagination.PagingInfo[domain.Transaction]) error {
	return persistence.DB.
		Offset((pagingInfo.PageNumber-1)*pagingInfo.ItemsPerPage).
		Limit(pagingInfo.ItemsPerPage).
		Order("created_at DESC").
		Find(&transactions, "account_id = ?", accountID).Error
}

func (ts transactionService) CountAllByAccount(accountID string, count *int64) error {
	return persistence.DB.Model(&domain.Transaction{}).Where("account_id = ?", accountID).Count(count).Error
}

func (ts transactionService) CashflowByAccount(accountID string, cashflow *Cashflow) error {
	month := utils.ConvertMonthToNumeric(time.Now().Month())

	var deposits float64
	if err := persistence.DB.
		Model(&domain.Transaction{}).
		Where("strftime('%m', created_at) = ? AND amount >= ? AND account_id = ?", month, 0, accountID).
		Select("sum(amount)").
		Row().
		Scan(&deposits); err != nil {
		deposits = 0
	}

	var withdrawals float64
	if err := persistence.DB.
		Model(&domain.Transaction{}).
		Where("strftime('%m', created_at) = ? AND amount <= ? AND account_id = ?", month, 0, accountID).
		Select("sum(amount)").
		Row().
		Scan(&withdrawals); err != nil {
		withdrawals = 0
	}

	cashflow.Deposits = deposits
	cashflow.Withdrawals = withdrawals

	return nil
}

func (s transactionService) BulkTransfer(customerIDs []string, transaction *domain.Transaction) error {
	return persistence.DB.Transaction(func(tx *gorm.DB) error {
		for _, customerID := range customerIDs {
			var customer domain.Customer
			if err := s.customerService.FindByID(customerID, &customer); err != nil {
				return err
			}

			newTransaction := domain.Transaction{
				AccountID:   customer.Accounts[0].ID,
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
