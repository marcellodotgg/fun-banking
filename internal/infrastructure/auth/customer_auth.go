package auth

import (
	"strconv"

	"github.com/bytebury/fun-banking/internal/domain"
	"github.com/bytebury/fun-banking/internal/infrastructure/persistence"
)

type customerAuth struct {
	jwtService JWTService
}

func NewCustomerAuth() customerAuth {
	return customerAuth{jwtService: JWTService{}}
}

func (a customerAuth) Login(customer domain.Customer) (string, error) {
	var dbCustomer domain.Customer
	if err := persistence.DB.First(&dbCustomer, "id = ?", customer.ID).Error; err != nil {
		return "", err
	}
	return a.jwtService.GenerateCustomerToken(strconv.Itoa(int(dbCustomer.ID)))
}
