package service

import (
	"strconv"

	"github.com/bytebury/fun-banking/internal/domain"
	"github.com/bytebury/fun-banking/internal/infrastructure/persistence"
)

type TransactionService interface {
	Create(transaction *domain.Transaction) error
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
	var account domain.Account
	if err := ts.accountService.FindByID(strconv.Itoa(int(transaction.AccountID)), &account); err != nil {
		return err
	}

	if account.Customer.Bank.UserID == *transaction.UserID {
		transaction.Status = domain.TransactionApproved
	}

	return persistence.DB.Create(&transaction).Error
}
