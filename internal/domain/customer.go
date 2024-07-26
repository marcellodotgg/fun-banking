package domain

import (
	"errors"
	"fmt"
	"regexp"
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
	BankID    int    `gorm:"not null; uniqueIndex:idx_bank_pin;"`
	Accounts  []Account
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

func (c Customer) PrimaryAccount() (Account, error) {
	if len(c.Accounts) == 0 {
		return Account{}, errors.New("no accounts")
	}
	for _, account := range c.Accounts {
		if account.IsPrimary {
			return account, nil
		}
	}
	return Account{}, errors.New("not found")
}

func (c *Customer) BeforeCreate(tx *gorm.DB) error {
	c.FirstName = strings.ToLower(c.FirstName)
	c.LastName = strings.ToLower(c.LastName)
	return c.validate()
}

func (c *Customer) BeforeUpdate(tx *gorm.DB) error {
	c.FirstName = strings.ToLower(c.FirstName)
	c.LastName = strings.ToLower(c.LastName)
	return c.validate()
}

func (c Customer) validate() error {
	if len(c.FirstName) > 20 || len(c.LastName) > 20 {
		return errors.New("first or last name is too long")
	}

	re := regexp.MustCompile("^[0-9]{4,6}$")
	if len(c.PIN) > 0 && !re.MatchString(c.PIN) {
		return errors.New("invalid PIN")
	}
	return nil
}
