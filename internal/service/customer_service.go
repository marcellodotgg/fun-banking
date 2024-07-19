package service

import (
	"github.com/bytebury/fun-banking/internal/domain"
	"github.com/bytebury/fun-banking/internal/infrastructure/persistence"
	"gorm.io/gorm"
)

type CustomerService interface {
	Create(customer *domain.Customer) error
	FindByID(id string, customer *domain.Customer) error
	Update(customer *domain.Customer) error
}

type customerService struct{}

func NewCustomerService() CustomerService {
	return customerService{}
}

func (s customerService) Create(customer *domain.Customer) error {
	return persistence.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&customer).Error; err != nil {
			return err
		}
		return tx.Create(&domain.Account{Name: "Checking", CustomerID: customer.ID}).Error
	})

}
func (s customerService) FindByID(id string, customer *domain.Customer) error {
	return persistence.DB.Preload("Bank.User").Preload("Accounts").First(&customer, "id = ?", id).Error
}

func (s customerService) Update(customer *domain.Customer) error {
	return persistence.DB.Updates(&customer).Error
}
