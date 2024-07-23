package service

import (
	"github.com/bytebury/fun-banking/internal/domain"
	"github.com/bytebury/fun-banking/internal/infrastructure/persistence"
)

type AccountService interface {
	FindByID(id string, account *domain.Account) error
	Create(account *domain.Account) error
	Update(id string, account *domain.Account) error
}

type accountService struct{}

func NewAccountService() AccountService {
	return accountService{}
}

func (s accountService) FindByID(id string, account *domain.Account) error {
	return persistence.DB.
		Preload("Customer.Bank").
		Preload("Customer.Accounts").
		First(&account, "id = ?", id).Error
}

func (s accountService) Create(account *domain.Account) error {
	return persistence.DB.Create(&account).Error
}

func (s accountService) Update(id string, account *domain.Account) error {
	if err := persistence.DB.Where("id = ?", id).Updates(&account).Error; err != nil {
		return err
	}
	return s.FindByID(id, account)
}
