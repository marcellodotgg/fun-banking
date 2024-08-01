package middleware

import (
	"os"

	"github.com/bytebury/fun-banking/internal/domain"
	"github.com/bytebury/fun-banking/internal/infrastructure/auth"
	"github.com/bytebury/fun-banking/internal/infrastructure/persistence"
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

		token, err := jwt.ParseWithClaims(tokenStr, &auth.UserClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil {
			c.Next()
			return
		}

		if claims, ok := token.Claims.(*auth.UserClaims); ok && token.Valid {
			c.Set("user_id", claims.UserID)
			c.Next()
		} else {
			c.Next()
		}
	}
}

func CustomerAudit() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, err := c.Cookie("customer_auth_token")

		if err != nil {
			c.Next()
			return
		}

		token, err := jwt.ParseWithClaims(tokenStr, &auth.CustomerClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil {
			c.Next()
			return
		}

		if claims, ok := token.Claims.(*auth.CustomerClaims); ok && token.Valid {
			c.Set("customer_id", claims.CustomerID)
			c.Next()
		} else {
			c.Next()
		}
	}
}

func PreferencesAudit() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("user_id")
		customerID := c.GetString("customer_id")

		if userID != "" {
			var user domain.User
			if err := persistence.DB.First(&user, "id = ?", userID).Error; err != nil {
				c.Next()
				return
			}
			c.Set("theme", user.Theme)
			c.Next()
			return
		}

		if customerID != "" {
			var customer domain.Customer
			if err := persistence.DB.Preload("Bank.User").First(&customer, "id = ?", customerID).Error; err != nil {
				c.Next()
				return
			}
			c.Set("theme", customer.Bank.User.Theme)
			c.Next()
			return
		}

		c.Next()
	}
}
