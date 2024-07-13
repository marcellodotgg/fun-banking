package domain

type Customer struct {
	Audit
	FirstName string `gorm:"not null"`
	LastName  string `gorm:"not null"`
	PIN       string `gorm:"uniqueIndex:idx_bank_pin; not null"`
	Bank      Bank   `gorm:"foreignKey:BankID; constraint:OnDelete:CASCADE"`
	BankID    uint   `gorm:"not null; uniqueIndex:idx_bank_pin; not null"`
}
