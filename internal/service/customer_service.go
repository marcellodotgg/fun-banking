package service

import (
	"github.com/bytebury/fun-banking/internal/domain"
	"github.com/bytebury/fun-banking/internal/infrastructure/persistence"
)

type CustomerService interface {
	Create(customer *domain.Customer) error
}

type customerService struct{}

func NewCustomerService() customerService {
	return customerService{}
}

func (s customerService) Create(customer *domain.Customer) error {
	return persistence.DB.Create(&customer).Error
}
