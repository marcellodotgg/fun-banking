package middleware

import (
	"net/http"

	"github.com/bytebury/fun-banking/internal/api/handler"
	"github.com/gin-gonic/gin"
)

func Authenticated() gin.HandlerFunc {
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
	formData := handler.NewFormData()
	formData.Errors["general"] = "You need to be signed in to access that resource"

	c.HTML(http.StatusUnauthorized, "sessions/signin", struct {
		FormData handler.FormData
		SignedIn bool
	}{FormData: formData, SignedIn: false})
}
