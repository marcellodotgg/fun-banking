package middleware

import (
	"github.com/bytebury/fun-banking/internal/api/handler"
	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetString("id")

		if id == "" {
			renderUnauthorized(c)
			c.Abort()
			return
		}

		c.Next()
	}
}

func renderUnauthorized(c *gin.Context) {
	signInHandler := handler.NewSessionHandler()
	signInHandler.FormData.Errors["general"] = "You need to be signed in to access that resource"
	signInHandler.SignIn(c)
}
