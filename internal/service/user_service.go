package service

import (
	"errors"
	"strings"

	"github.com/bytebury/fun-banking/internal/domain"
	"github.com/bytebury/fun-banking/internal/infrastructure/mail"
	"github.com/bytebury/fun-banking/internal/infrastructure/pagination"
	"github.com/bytebury/fun-banking/internal/infrastructure/persistence"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService interface {
	Create(user *domain.User) error
	Update(id string, user *domain.User) error
	UpdatePassword(id, password string) error
	FindByID(id string, user *domain.User) error
	FindByEmail(email string, user *domain.User) error
	Search(search string, pagingInfo *pagination.PagingInfo[domain.User]) error
	FindPendingTransactions(id string, transactions *[]domain.Transaction) error
}

type userService struct{}

func NewUserService() UserService {
	return userService{}
}

func (s userService) Create(user *domain.User) error {
	if len(user.Password) < 6 {
		return errors.New("password too short")
	}

	passwordHash, err := s.hashPassword(user.Password)

	if err != nil {
		return err
	}

	user.Password = passwordHash

	if err := persistence.DB.Create(&user).Error; err != nil {
		return err
	}

	return mail.NewWelcomeMailer().Send(user.Email, *user)
}

func (s userService) FindByID(id string, user *domain.User) error {
	return persistence.DB.Preload("Banks").Preload("Subscriptions").First(&user, "id = ?", id).Error
}

func (s userService) FindByEmail(email string, user *domain.User) error {
	return persistence.DB.First(&user, "email = ?", strings.TrimSpace(strings.ToLower(email))).Error
}

func (s userService) Search(search string, pagingInfo *pagination.PagingInfo[domain.User]) error {
	var users []domain.User
	var usersCount int64

	search = strings.ToLower(search)

	// Count the users first
	persistence.DB.
		Where("username LIKE ?", "%"+search+"%").
		Or("email LIKE ?", "%"+search+"%").
		Or("first_name LIKE ?", "%"+search+"%").
		Or("last_name LIKE ?", "%"+search+"%").
		Count(&usersCount)

	// Find the users
	persistence.DB.
		Preload("Subscriptions", func(db *gorm.DB) *gorm.DB {
			// Order subscriptions in reverse
			return db.Order("created_at DESC")
		}).
		Order("created_at DESC").
		Where("username LIKE ?", "%"+search+"%").
		Or("email LIKE ?", "%"+search+"%").
		Or("first_name LIKE ?", "%"+search+"%").
		Or("last_name LIKE ?", "%"+search+"%").
		Offset((pagingInfo.PageNumber - 1) * pagingInfo.ItemsPerPage).
		Limit(pagingInfo.ItemsPerPage).
		Find(&users)

	pagingInfo.Items = users
	pagingInfo.TotalItems = usersCount

	return nil
}

func (s userService) Update(id string, user *domain.User) error {
	return persistence.DB.Where("id = ?", id).Updates(&user).Error
}

func (s userService) UpdatePassword(id, password string) error {
	var user domain.User
	if err := s.FindByID(id, &user); err != nil {
		return err
	}

	newPassword, err := s.hashPassword(password)

	if err != nil {
		return err
	}

	user.Password = newPassword

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
