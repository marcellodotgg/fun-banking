package service

import (
	"errors"
	"os"

	"github.com/bytebury/fun-banking/internal/infrastructure/auth"
	"github.com/dgrijalva/jwt-go"
)

type TokenService struct {
	jwtService *auth.JWTService
}

func NewTokenService() TokenService {
	return TokenService{
		jwtService: &auth.JWTService{},
	}
}

func (s TokenService) GetUserIDFromToken(tokenStr string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &auth.UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(*auth.UserClaims); ok && token.Valid {
		return claims.UserID, nil
	}

	return "", errors.New("invalid token or the token is expired")
}

func (s TokenService) GenerateUserToken(userID string) (string, error) {
	return s.jwtService.GenerateToken(userID)
}

func (s TokenService) GenerateTempToken(userID string) (string, error) {
	return s.jwtService.GenerateTempToken(userID)
}
