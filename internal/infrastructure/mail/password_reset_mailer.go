package mail

import (
	"os"
	"strconv"

	"github.com/bytebury/fun-banking/internal/domain"
	"github.com/bytebury/fun-banking/internal/infrastructure/auth"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type passwordResetMailer struct {
	mailer     *Mailer
	jwtService *auth.JWTService
}

func NewPasswordResetMailer() passwordResetMailer {
	return passwordResetMailer{
		mailer:     &Mailer{},
		jwtService: &auth.JWTService{},
	}
}

func (m passwordResetMailer) Send(to string, user domain.User) error {
	user.FirstName = cases.Title(language.AmericanEnglish).String(user.FirstName)
	token, _ := m.jwtService.GenerateTempToken(strconv.Itoa(user.ID))

	return m.mailer.Send(to, "Fun Banking - Reset Your Password", "reset_password", struct {
		User       domain.User
		Token      string
		WebsiteURL string
	}{
		User:       user,
		Token:      token,
		WebsiteURL: os.Getenv("WEBSITE_URL"),
	})
}
