package service

import (
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
	FindAllByAccount(accountID string, transactions *[]domain.Transaction, pagingInfo pagination.PagingInfo[domain.Transaction]) error
	CountAllByAccount(accountID string, count *int64) error
	CashflowByAccount(accountID string, cashflow *Cashflow) error
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

func (ts transactionService) CashflowByAccount(accountID string, cashflow *Cashflow) error {
	month := utils.ConvertMonthToNumeric(time.Now().Month())

	var deposits float64
	if err := persistence.DB.
		Model(&domain.Transaction{}).
		Where("strftime('%m', created_at) = ? AND amount >= ?", month, 0).
		Select("sum(amount)").
		Row().
		Scan(&deposits); err != nil {
		deposits = 0
	}

	var withdrawals float64
	if err := persistence.DB.
		Model(&domain.Transaction{}).
		Where("strftime('%m', created_at) = ? AND amount <= ?", month, 0).
		Select("sum(amount)").
		Row().
		Scan(&withdrawals); err != nil {
		withdrawals = 0
	}

	cashflow.Deposits = deposits
	cashflow.Withdrawals = withdrawals

	return nil
}
