package service

import (
	"github.com/bytebury/fun-banking/internal/domain"
	"github.com/bytebury/fun-banking/internal/infrastructure/persistence"
)

type TransactionService interface {
	Create(transaction *domain.Transaction) error
}

type transactionService struct{}

func NewTransactionService() TransactionService {
	return transactionService{}
}

func (ts transactionService) Create(transaction *domain.Transaction) error {
	return persistence.DB.Create(&transaction).Error
}
