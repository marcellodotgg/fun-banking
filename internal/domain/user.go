package domain

import (
	"fmt"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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
	Banks     []Bank
	ImageURL  string `gorm:"not null; default:https://static.vecteezy.com/system/resources/previews/009/292/244/large_2x/default-avatar-icon-of-social-media-user-vector.jpg"`
}

func (u User) FullName() string {
	return cases.Title(language.AmericanEnglish).String(fmt.Sprintf("%s %s", u.FirstName, u.LastName))
}

func (u User) IsAdmin() bool {
	return u.Role == UserRoleAdmin
}

func (u User) IsFree() bool {
	return u.Role == UserRoleFree
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
