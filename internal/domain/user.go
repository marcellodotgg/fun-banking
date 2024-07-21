package domain

import (
	"strings"

	"gorm.io/gorm"
)

const (
	UserRoleFree  = "FREE"
	UserRoleAdmin = "ADMIN"
)

type User struct {
	Audit
	Username  string `gorm:"unique; not null"`
	Email     string `gorm:"unique; not null"`
	FirstName string `gorm:"not null; size:20"`
	LastName  string `gorm:"not null; size:20"`
	Role      string `gorm:"not null; default:FREE"`
	Password  string `json:"-" gorm:"not null"`
	ImageURL  string `gorm:"not null; default:https://i.pinimg.com/736x/0d/64/98/0d64989794b1a4c9d89bff571d3d5842.jpg"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.Username = strings.TrimSpace(strings.ToLower(u.Username))
	u.Email = strings.TrimSpace(strings.ToLower(u.Email))
	u.FirstName = strings.TrimSpace(strings.ToLower(u.FirstName))
	u.LastName = strings.TrimSpace(strings.ToLower(u.LastName))
	return nil
}

func (u *User) BeforeUpdate(tx *gorm.DB) error {
	u.Username = strings.TrimSpace(strings.ToLower(u.Username))
	u.Email = strings.TrimSpace(strings.ToLower(u.Email))
	u.FirstName = strings.TrimSpace(strings.ToLower(u.FirstName))
	u.LastName = strings.TrimSpace(strings.ToLower(u.LastName))
	return nil
}
