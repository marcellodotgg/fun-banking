package middleware

import (
	"os"

	"github.com/bytebury/fun-banking/internal/infrastructure/auth"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func Audit() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, err := c.Cookie("auth_token")

		if err != nil {
			c.Next()
			return
		}

		token, err := jwt.ParseWithClaims(tokenStr, &auth.AuditClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil {
			c.Next()
			return
		}

		if claims, ok := token.Claims.(*auth.AuditClaims); ok && token.Valid {
			c.Set("id", claims.ID)
			c.Next()
		} else {
			c.Next()
		}
	}
}
