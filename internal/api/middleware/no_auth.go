package middleware

import (
	"github.com/bytebury/fun-banking/internal/api/handler"
	"github.com/gin-gonic/gin"
)

func NoAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("user_id")
		customerID := c.GetString("customer_id")

		if userID != "" || customerID != "" {
			homePageHandler := handler.NewHomePageHandler()
			homePageHandler.Homepage(c)
			c.Abort()
			return
		}
		c.Next()
	}
}
