package domain

import (
	"strings"

	"gorm.io/gorm"
)

type Bank struct {
	Audit
	Name        string `gorm:"not null"`
	Description string
	Slug        string `gorm:"not null; uniqueIndex:idx_user_slug"`
	UserID      string `gorm:"not null; uniqueIndex:idx_user_slug"`
	User        User   `gorm:"foreignKey:UserID; constraint:OnDelete:CASCADE"`
	Customers   []Customer
}

func (b *Bank) BeforeCreate(tx *gorm.DB) error {
	b.Slug = strings.ToLower(strings.ReplaceAll(b.Name, " ", "-"))
	return nil
}

func (b *Bank) BeforeUpdate(tx *gorm.DB) error {
	b.Slug = strings.ToLower(strings.ReplaceAll(b.Name, " ", "-"))
	return nil
}
