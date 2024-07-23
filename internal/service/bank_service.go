package service

import (
	"github.com/bytebury/fun-banking/internal/domain"
	"github.com/bytebury/fun-banking/internal/infrastructure/persistence"
	"gorm.io/gorm"
)

type BankService interface {
	MyBanks(userID string, banks *[]domain.Bank) error
	Create(bank *domain.Bank) error
	Update(id string, bank *domain.Bank) error
	FindByID(id string, bank *domain.Bank) error
	FindByUsernameAndSlug(username, slug string, bank *domain.Bank) error
}

type bankService struct{}

func NewBankService() BankService {
	return bankService{}
}

func (s bankService) MyBanks(userID string, banks *[]domain.Bank) error {
	return persistence.DB.Preload("User").Find(&banks, "user_id = ?", userID).Error
}

func (s bankService) Create(bank *domain.Bank) error {
	if err := persistence.DB.Create(&bank).Error; err != nil {
		return err
	}

	return persistence.DB.Preload("User").First(&bank).Error
}

func (s bankService) Update(id string, bank *domain.Bank) error {
	var oldBank domain.Bank
	if err := persistence.DB.First(&oldBank, "id = ?", id).Error; err != nil {
		return err
	}

	if err := persistence.DB.Model(&oldBank).Updates(&bank).Error; err != nil {
		return err
	}
	return s.FindByID(id, bank)
}

func (s bankService) FindByID(id string, bank *domain.Bank) error {
	return persistence.DB.
		Preload("User").
		Preload("Customers", func(db *gorm.DB) *gorm.DB {
			return db.Order("last_name ASC, first_name ASC")
		}).
		Preload("Customers.Accounts").
		First(&bank, "id = ?", id).Error
}

func (s bankService) FindByUsernameAndSlug(username, slug string, bank *domain.Bank) error {
	return persistence.DB.
		Preload("User").
		Joins("JOIN users ON users.id = banks.user_id").
		Where("users.username = ? AND banks.slug = ?", username, slug).
		First(&bank).Error
}
