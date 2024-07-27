package middleware

import (
	"github.com/gin-gonic/gin"
)

func NoAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("user_id")
		customerID := c.GetString("customer_id")

		if userID != "" || customerID != "" {
			renderForbidden(c)
			c.Abort()
			return
		}
		c.Next()
	}
}
