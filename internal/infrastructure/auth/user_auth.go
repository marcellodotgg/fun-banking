package auth

import (
	"errors"

	"github.com/bytebury/fun-banking/internal/domain"
	"github.com/bytebury/fun-banking/internal/infrastructure/persistence"
	"golang.org/x/crypto/bcrypt"
)

type userAuth struct {
	jwtService JWTService
}

func NewUserAuth() userAuth {
	return userAuth{jwtService: JWTService{}}
}

func (a userAuth) Login(usernameOrEmail, password string) (string, error) {
	var user domain.User
	if err := persistence.DB.First(&user, "username = ? OR email = ?", usernameOrEmail, usernameOrEmail).Error; err != nil {
		return "", err
	}

	if !a.verifyPassword(password, user.Password) {
		return "", errors.New("INVALID password")
	}

	return a.jwtService.GenerateToken(user.ID)
}

func (a userAuth) verifyPassword(password, passwordHash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)) == nil
}
