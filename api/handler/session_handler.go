package handler

import (
	"net/http"

	"github.com/bytebury/fun-banking/internal/domain"
	"github.com/bytebury/fun-banking/internal/infrastructure/persistence"
	"github.com/gin-gonic/gin"
)

type SessionHandler struct {
}

func (h SessionHandler) SignIn(c *gin.Context) {
	c.HTML(http.StatusOK, "sessions/signin", nil)
}

func (h SessionHandler) CreateSession(c *gin.Context) {
	c.Request.ParseForm()

	username := c.PostForm("email_or_username")
	password := c.PostForm("password")

	formData := newFormData()

	var user domain.User
	if err := persistence.DB.First(&user, "username = ? OR email = ?", username, username).Error; err != nil {
		formData.Errors["email_or_username"] = "Unable to sign you in. Invalid credentials."
		c.HTML(http.StatusUnauthorized, "sessions/signin_form", formData)
		return
	}

	if user.Password != password {
		formData.Errors["email_or_username"] = "Unable to sign you in. Invalid credentials."
		c.HTML(http.StatusUnauthorized, "sessions/signin_form", formData)
		return
	}

	c.Header("HX-Redirect", "/")
}
