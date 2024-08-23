package domain

import "time"

type AutoPay struct {
	Audit
	Cadence     string    `gorm:"not null; default:day"`
	StartDate   time.Time `gorm:"not null"`
	Amount      float64   `gorm:"not null; default:0.00; type:decimal(50,2)"`
	Description string    `gorm:"not null; size:255"`
	AccountID   int       `gorm:"not null"`
	Account     Account   `gorm:"foreignKey:AccountID; constraint:OnDelete:CASCADE"`
	Active      bool      `gorm:"not null;default:true"`
}
