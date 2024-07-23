package domain

import (
	"bytes"
	"errors"
	"strings"

	"gorm.io/gorm"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

type Announcement struct {
	Audit
	ID          string `gorm:"primary_key"`
	UserID      uint
	User        User   `gorm:"foreignKey:UserID; constraint:OnDelete:CASCADE"`
	Title       string `gorm:"not null"`
	Description string `gorm:"not null"`
}

func (a Announcement) HTML() string {
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM, extension.Typographer),
		goldmark.WithParserOptions(parser.WithAutoHeadingID()),
		goldmark.WithRendererOptions(html.WithHardWraps()),
	)

	var buf bytes.Buffer
	markdown := []byte(a.Description)

	if err := md.Convert(markdown, &buf); err != nil {
		panic(err)
	}

	return buf.String()
}

func (a *Announcement) BeforeCreate(tx *gorm.DB) error {
	a.Title = strings.TrimSpace(a.Title)
	a.Description = strings.TrimSpace(a.Description)
	return a.validate()
}

func (a *Announcement) BeforeUpdate(tx *gorm.DB) error {
	a.Title = strings.TrimSpace(a.Title)
	a.Description = strings.TrimSpace(a.Description)
	return a.validate()
}

func (a Announcement) validate() error {
	if len(a.Title) < 2 {
		return errors.New("title is too short")
	}

	if len(a.Description) < 10 {
		return errors.New("description is too short")
	}

	return nil
}
