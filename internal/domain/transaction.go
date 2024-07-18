package domain

const (
	TransactionPending  = "PENDING"
	TransactionApproved = "APPROVED"
	TransactionDeclined = "DECLINED"
)

type Transaction struct {
	Audit
	Description string  `gorm:"not null; size:255"`
	Balance     float64 `gorm:"not null; type:decimal(50,2)"`
	Amount      float64 `gorm:"not null; type:decimal(50,2)"`
	Status      string  `gorm:"not null; size:20; default:PENDING"`
	AccountID   uint    `gorm:"not null"`
	Account     Account `gorm:"foreignKey:AccountID; constraint:OnDelete:CASCADE"`
}
