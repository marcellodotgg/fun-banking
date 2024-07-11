package auth

import (
	"errors"
	"strconv"
	"strings"

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

func (a userAuth) Login(request domain.UserSignInRequest) (string, domain.User, error) {
	request.UsernameOrEmail = strings.TrimSpace(request.UsernameOrEmail)

	var user domain.User
	if err := persistence.DB.First(&user, "username = ? OR email = ?", request.UsernameOrEmail, request.UsernameOrEmail).Error; err != nil {
		return "", user, err
	}

	if !a.verifyPassword(request.Password, user.Password) {
		return "", user, errors.New("INVALID password")
	}

	token, err := a.jwtService.GenerateToken(strconv.Itoa(int(user.ID)))

	return token, user, err
}

func (a userAuth) verifyPassword(password, passwordHash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)) == nil
}
