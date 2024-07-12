package domain

type User struct {
	Audit
	Username  string `gorm:"unique; not null"`
	Email     string `gorm:"unique; not null"`
	FirstName string `gorm:"not null; size:20"`
	LastName  string `gorm:"not null; size:20"`
	Role      string `gorm:"not null; default:FREE"`
	Password  string `json:"-" gorm:"not null"`
}
