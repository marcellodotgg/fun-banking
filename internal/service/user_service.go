package service

import (
	"github.com/bytebury/fun-banking/internal/domain"
	"github.com/bytebury/fun-banking/internal/infrastructure/persistence"
)

type UserService interface {
	Count() int64
}

type userService struct{}

func NewUserService() UserService {
	return userService{}
}

func (service userService) Count() int64 {
	var count int64
	persistence.DB.Model(&domain.User{}).Count(&count)
	return count
}
