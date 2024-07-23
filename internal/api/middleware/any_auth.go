package middleware

import "github.com/gin-gonic/gin"

func AnyAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("user_id")
		customerID := c.GetString("customer_id")

		if userID == "" && customerID == "" {
			renderUnauthorized(c)
			c.Abort()
			return
		}

		c.Next()
	}
}