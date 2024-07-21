package middleware

import (
	"github.com/gin-gonic/gin"
)

func CustomerAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetString("customer_id")

		if id == "" {
			renderUnauthorized(c)
			c.Abort()
			return
		}

		c.Next()
	}
}
