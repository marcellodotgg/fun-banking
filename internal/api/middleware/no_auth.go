package middleware

import (
	"github.com/bytebury/fun-banking/internal/api/handler"
	"github.com/gin-gonic/gin"
)

func NoAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetString("id")

		if id != "" {
			homePageHandler := handler.NewHomePageHandler()
			homePageHandler.Dashboard(c)
			c.Abort()
			return
		}
		c.Next()
	}
}
