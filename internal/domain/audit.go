package domain

import "time"

type Audit struct {
	ID        uint      `gorm:"primary_key"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
