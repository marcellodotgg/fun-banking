package service

import (
	"strings"

	"github.com/bytebury/fun-banking/internal/domain"
	"github.com/bytebury/fun-banking/internal/infrastructure/persistence"
	"gorm.io/gorm"
)

type CustomerService interface {
	Create(customer *domain.Customer) error
	FindByID(id string, customer *domain.Customer) error
	Update(customer *domain.Customer) error
	FindAllByBankIDAndName(bankID, name string, limit int, customers *[]domain.Customer) error
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

func (s customerService) FindAllByBankIDAndName(bankID, name string, limit int, customers *[]domain.Customer) error {
	var firstName, lastName string
	var names = strings.Split(strings.ToLower(name), " ")

	if len(names) == 1 {
		firstName = names[0]
	}

	if len(names) >= 2 {
		firstName = names[0]
		lastName = names[1]
	}

	return persistence.DB.
		Preload("Accounts").
		Order("last_name ASC, first_name ASC").
		Limit(limit).
		Find(&customers, "bank_id = ? AND ((first_name LIKE ? AND last_name LIKE ?) OR (last_name LIKE ? AND first_name LIKE ?))", bankID, "%"+firstName+"%", "%"+lastName+"%", "%"+firstName+"%", "%"+lastName+"%").Error
}
