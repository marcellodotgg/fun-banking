package service

import (
	"github.com/bytebury/fun-banking/internal/domain"
	"github.com/bytebury/fun-banking/internal/infrastructure/persistence"
)

type AccountService interface {
	FindByID(id string, account *domain.Account) error
	Create(account *domain.Account) error
	Update(account *domain.Account) error
}

type accountService struct{}

func NewAccountService() AccountService {
	return accountService{}
}

func (as accountService) FindByID(id string, account *domain.Account) error {
	return persistence.DB.
		Preload("Customer.Accounts").
		Preload("Customer.Bank").
		First(&account, "id = ?", id).Error
}

func (as accountService) Create(account *domain.Account) error {
	return persistence.DB.Create(&account).Error
}

func (as accountService) Update(account *domain.Account) error {
	return persistence.DB.Updates(&account).Error
}
