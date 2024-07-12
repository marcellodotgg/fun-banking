package middleware

import (
	"net/http"
	"os"

	"github.com/bytebury/fun-banking/internal/api/handler"
	"github.com/bytebury/fun-banking/internal/infrastructure/auth"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func Authenticated() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, err := c.Cookie("auth_token")

		if err != nil {
			renderUnauthorizedAndAbort(c)
			return
		}

		token, err := jwt.ParseWithClaims(tokenStr, &auth.AuditClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil {
			renderUnauthorizedAndAbort(c)
			return
		}

		if claims, ok := token.Claims.(*auth.AuditClaims); ok && token.Valid {
			c.Set("id", claims.ID)
			c.Next()
		} else {
			renderUnauthorizedAndAbort(c)
		}
	}
}

func renderUnauthorizedAndAbort(c *gin.Context) {
	formData := handler.NewFormData()
	formData.Errors["general"] = "You need to be signed in to access that resource"

	c.HTML(http.StatusUnauthorized, "sessions/signin", formData)
	c.Abort()
}
