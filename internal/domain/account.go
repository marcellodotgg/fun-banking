package domain

import (
	"errors"
	"strings"

	"gorm.io/gorm"
)

type Account struct {
	Audit
	Name       string   `gorm:"not null; default:Checking"`
	Balance    float64  `gorm:"decimal(50,2); default:0.00"`
	CustomerID int      `gorm:"not null"`
	Customer   Customer `gorm:"foreignKey:CustomerID; constraint:OnDelete:CASCADE"`
	IsPrimary  bool     `gorm:"not null;default:0"`
}

func (a *Account) BeforeCreate(tx *gorm.DB) error {
	a.Name = strings.TrimSpace(a.Name)
	return a.validate()
}

func (a *Account) BeforeUpdate(tx *gorm.DB) error {
	a.Name = strings.TrimSpace(a.Name)
	return a.validate()
}

func (a Account) validate() error {
	if len(a.Name) > 25 {
		return errors.New("name too long")
	}
	return nil
}
