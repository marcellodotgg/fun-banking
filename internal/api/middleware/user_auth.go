package middleware

import (
	"github.com/bytebury/fun-banking/internal/api/handler"
	"github.com/gin-gonic/gin"
)

func UserAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetString("user_id")

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
	signInHandler.Form.Errors["general"] = "You need to be signed in to access that resource"
	signInHandler.SignIn(c)
}
