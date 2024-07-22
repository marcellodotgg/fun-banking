package middleware

import (
	"github.com/bytebury/fun-banking/internal/domain"
	"github.com/bytebury/fun-banking/internal/service"
	"github.com/gin-gonic/gin"
)

func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetString("user_id")

		userService := service.NewUserService()

		var user domain.User
		if err := userService.FindByID(id, &user); err != nil {
			renderUnauthorized(c)
			c.Abort()
			return
		}

		if !user.IsAdmin() {
			renderUnauthorized(c)
			c.Abort()
			return
		}

		c.Next()
	}
}
