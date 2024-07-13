package domain

type Bank struct {
	Audit
	Name        string `gorm:"not null"`
	Description string
	Slug        string `gorm:"not null; uniqueIndex:idx_user_slug"`
	UserID      uint   `gorm:"not null; uniqueIndex:idx_user_slug"`
	User        User   `gorm:"foreignKey:UserID; constraint:OnDelete:CASCADE"`
	Customers   []Customer
}
