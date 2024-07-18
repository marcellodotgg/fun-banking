package service

import (
	"strconv"

	"github.com/bytebury/fun-banking/internal/domain"
	"github.com/bytebury/fun-banking/internal/infrastructure/pagination"
	"github.com/bytebury/fun-banking/internal/infrastructure/persistence"
	"gorm.io/gorm"
)

type TransactionService interface {
	Create(transaction *domain.Transaction) error
	FindAllByAccount(accountID string, transactions *[]domain.Transaction, pagingInfo pagination.PagingInfo[domain.Transaction]) error
	CountAllByAccount(accountID string, count *int64) error
}

type transactionService struct {
	accountService AccountService
}

func NewTransactionService() TransactionService {
	return transactionService{
		accountService: NewAccountService(),
	}
}

func (ts transactionService) Create(transaction *domain.Transaction) error {
	return persistence.DB.Transaction(func(tx *gorm.DB) error {
		var account domain.Account
		if err := ts.accountService.FindByID(strconv.Itoa(int(transaction.AccountID)), &account); err != nil {
			return err
		}

		transaction.Balance = account.Balance

		if account.Customer.Bank.UserID == *transaction.UserID {
			transaction.Status = domain.TransactionApproved
			account.Balance += transaction.Amount

			if err := ts.accountService.Update(&account); err != nil {
				return err
			}
		}

		return persistence.DB.Create(&transaction).Error
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
