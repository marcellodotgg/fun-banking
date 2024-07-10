package domain

type User struct {
	audit
	Username  string
	Email     string
	FirstName string
	LastName  string
}
