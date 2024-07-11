package auth

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type AuditClaims struct {
	ID string
	jwt.StandardClaims
}

type JWTService struct{}

func (j *JWTService) GenerateToken(id string) (string, error) {
	claims := AuditClaims{
		ID: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * 365 * 100 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
		}}
	return generateToken(claims)
}

func generateToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}
