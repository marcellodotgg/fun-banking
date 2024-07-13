package service

import (
	"strings"

	"github.com/bytebury/fun-banking/internal/domain"
	"github.com/bytebury/fun-banking/internal/infrastructure/persistence"
)

type BankService interface {
	MyBanks(userID string, banks *[]domain.Bank) error
	Create(bank *domain.Bank) error
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
	bank.Slug = strings.ToLower(strings.ReplaceAll(bank.Name, " ", "-"))

	if err := persistence.DB.Create(&bank).Error; err != nil {
		return err
	}

	return persistence.DB.Preload("User").First(&bank).Error
}

func (s bankService) FindByID(id string, bank *domain.Bank) error {
	return persistence.DB.Preload("User").Preload("Customers").First(&bank, "id = ?", id).Error
}

func (s bankService) FindByUsernameAndSlug(username, slug string, bank *domain.Bank) error {
	return persistence.DB.
		Preload("User").
		Preload("Customers").
		Joins("JOIN users ON users.id = banks.user_id").
		Where("banks.slug = ? AND users.username = ?", slug, username).
		First(&bank).Error
}
