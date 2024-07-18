package domain

import (
	"fmt"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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

func (c Customer) FullName() string {
	return cases.Title(language.AmericanEnglish).String(fmt.Sprintf("%s %s", c.FirstName, c.LastName))
}

func (c Customer) NetWorth() float64 {
	netWorth := float64(0)
	for _, account := range c.Accounts {
		netWorth += account.Balance
	}
	return netWorth
}
