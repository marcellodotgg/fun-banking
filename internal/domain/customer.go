package domain

import (
	"strings"

	"gorm.io/gorm"
)

type Customer struct {
	Audit
	FirstName string `gorm:"not null"`
	LastName  string `gorm:"not null"`
	PIN       string `gorm:"uniqueIndex:idx_bank_pin; not null"`
	Bank      Bank   `gorm:"foreignKey:BankID; constraint:OnDelete:CASCADE"`
	BankID    uint   `gorm:"not null; uniqueIndex:idx_bank_pin;"`
	Accounts  []Account
}

func (c *Customer) BeforeCreate(tx *gorm.DB) error {
	c.FirstName = strings.ToLower(c.FirstName)
	c.LastName = strings.ToLower(c.LastName)
	return nil
}

func (c *Customer) BeforeUpdate(tx *gorm.DB) error {
	c.FirstName = strings.ToLower(c.FirstName)
	c.LastName = strings.ToLower(c.LastName)
	return nil
}
