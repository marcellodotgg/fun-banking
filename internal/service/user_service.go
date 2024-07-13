package service

import (
	"github.com/bytebury/fun-banking/internal/domain"
	"github.com/bytebury/fun-banking/internal/infrastructure/persistence"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Count() int64
	Create(user *domain.User) error
}

type userService struct{}

func NewUserService() UserService {
	return userService{}
}

func (s userService) Count() int64 {
	var count int64
	persistence.DB.Model(&domain.User{}).Count(&count)
	return count
}

func (s userService) Create(user *domain.User) error {
	passwordHash, err := s.hashPassword(user.Password)

	if err != nil {
		return err
	}

	user.Password = passwordHash

	return persistence.DB.Create(&user).Error
}

func (s userService) hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}
