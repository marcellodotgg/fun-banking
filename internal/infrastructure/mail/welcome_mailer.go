package mail

import (
	"strconv"

	"github.com/bytebury/fun-banking/internal/domain"
	"github.com/bytebury/fun-banking/internal/infrastructure/auth"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type welcomeMailer struct {
	mailer     *Mailer
	jwtService *auth.JWTService
}

func NewWelcomeMailer() welcomeMailer {
	return welcomeMailer{
		mailer:     &Mailer{},
		jwtService: &auth.JWTService{},
	}
}

func (m welcomeMailer) Send(to string, user domain.User) error {
	user.FirstName = cases.Title(language.AmericanEnglish).String(user.FirstName)
	token, _ := m.jwtService.GenerateTempToken(strconv.Itoa(user.ID))

	return m.mailer.Send(to, "Fun Banking - Welcome!", "welcome_email", struct {
		User  domain.User
		Token string
	}{
		User:  user,
		Token: token,
	})
}
