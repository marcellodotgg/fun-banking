package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type SessionHandler struct {
}

func (h SessionHandler) SignIn(c *gin.Context) {
	c.HTML(http.StatusOK, "sessions/signin", nil)
}

func (h SessionHandler) CreateSession(c *gin.Context) {
	c.Request.ParseForm()

	emailOrUsername := c.PostForm("email_or_username")
	password := c.PostForm("password")

	if emailOrUsername == "marcello" && password == "password" {
		c.Header("HX-Redirect", "/")
	}

	formData := newFormData()
	formData.Errors["email_or_username"] = "Unable to sign you in. Invalid credentials."

	c.HTML(http.StatusUnauthorized, "sessions/signin_form", formData)
}
