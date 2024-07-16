package service

import (
	"github.com/bytebury/fun-banking/internal/domain"
	"github.com/bytebury/fun-banking/internal/infrastructure/persistence"
)

type AccountService interface {
	Create(account *domain.Account) error
}

type accountService struct{}

func NewAccountService() AccountService {
	return accountService{}
}

func (as accountService) Create(account *domain.Account) error {
	return persistence.DB.Create(&account).Error
}
