package service

import (
	"strings"

	"github.com/bytebury/fun-banking/internal/domain"
	"github.com/bytebury/fun-banking/internal/infrastructure/persistence"
)

type BankService interface {
	MyBanks(userID string) []domain.Bank
	Create(bank *domain.Bank) error
}

type bankService struct{}

func NewBankService() BankService {
	return bankService{}
}

func (s bankService) MyBanks(userID string) []domain.Bank {
	var banks []domain.Bank
	if err := persistence.DB.Preload("User").Find(&banks, "user_id = ?", userID).Error; err != nil {
		return []domain.Bank{}
	}
	return banks
}

func (s bankService) Create(bank *domain.Bank) error {
	bank.Slug = strings.ToLower(strings.ReplaceAll(bank.Name, " ", "-"))
	return persistence.DB.Create(&bank).Error
}
