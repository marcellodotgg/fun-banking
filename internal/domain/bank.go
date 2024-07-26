package domain

import (
	"errors"
	"regexp"
	"strings"

	"gorm.io/gorm"
)

type Bank struct {
	Audit
	Name        string `gorm:"not null"`
	Description string
	Slug        string `gorm:"not null; uniqueIndex:idx_user_slug"`
	UserID      int    `gorm:"not null; uniqueIndex:idx_user_slug"`
	User        User   `gorm:"foreignKey:UserID; constraint:OnDelete:CASCADE"`
	Customers   []Customer
}

func (b *Bank) BeforeCreate(tx *gorm.DB) error {
	b.Slug = strings.ToLower(strings.ReplaceAll(b.Name, " ", "-"))
	return b.validate()
}

func (b *Bank) BeforeUpdate(tx *gorm.DB) error {
	b.Slug = strings.ToLower(strings.ReplaceAll(b.Name, " ", "-"))
	return b.validate()
}

func (b Bank) validate() error {
	if len(b.Name) > 25 {
		return errors.New("name too long")
	}

	if len(b.Description) > 500 {
		return errors.New("description too long")
	}

	re := regexp.MustCompile(`^[A-Za-z0-9](?:\s?[A-Za-z0-9])*$`)
	if !re.MatchString(b.Name) {
		return errors.New("invalid name")
	}

	return nil
}
