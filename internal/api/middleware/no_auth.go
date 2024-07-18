package middleware

import (
	"github.com/bytebury/fun-banking/internal/api/handler"
	"github.com/gin-gonic/gin"
)

func NoAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("user_id")

		if userID != "" {
			homePageHandler := handler.NewHomePageHandler()
			homePageHandler.Dashboard(c)
			c.Abort()
			return
		}
		c.Next()
	}
}
