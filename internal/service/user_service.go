package service

import (
	"strconv"

	"github.com/bytebury/fun-banking/internal/domain"
	"github.com/bytebury/fun-banking/internal/infrastructure/persistence"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Create(user *domain.User) error
	Update(id string, user *domain.User) error
	FindByID(id string, user *domain.User) error
	FindPendingTransactions(id string, transactions *[]domain.Transaction) error
}

type userService struct{}

func NewUserService() UserService {
	return userService{}
}

func (s userService) Create(user *domain.User) error {
	passwordHash, err := s.hashPassword(user.Password)

	if err != nil {
		return err
	}

	user.Password = passwordHash

	return persistence.DB.Create(&user).Error
}

func (s userService) FindByID(id string, user *domain.User) error {
	return persistence.DB.First(&user, "id = ?", id).Error
}

func (s userService) Update(id string, user *domain.User) error {
	userID, err := strconv.Atoi(id)

	if err != nil {
		return err
	}

	user.ID = uint(userID)

	return persistence.DB.Updates(&user).Error
}

func (s userService) FindPendingTransactions(id string, transactions *[]domain.Transaction) error {
	return persistence.DB.
		Joins("JOIN users on users.id = banks.user_id").
		Joins("JOIN banks on banks.id = customers.bank_id").
		Joins("JOIN customers on customers.id = accounts.customer_id").
		Joins("JOIN accounts on accounts.id = transactions.account_id").
		Where("users.id = ?", id).
		Where("transactions.status = 'PENDING'").
		Preload("Account").
		Preload("Account.Customer").
		Select("transactions.*").
		Find(&transactions).Error
}

func (s userService) hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}
