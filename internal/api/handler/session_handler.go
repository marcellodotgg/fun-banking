package handler

import (
	"net/http"

	"github.com/bytebury/fun-banking/internal/domain"
	"github.com/bytebury/fun-banking/internal/infrastructure/auth"
	"github.com/bytebury/fun-banking/internal/infrastructure/persistence"
	"github.com/gin-gonic/gin"
)

type sessionHandler struct {
	FormData FormData
	SignedIn bool
}

func NewSessionHandle() sessionHandler {
	return sessionHandler{
		FormData: NewFormData(),
		SignedIn: false,
	}
}

func (h sessionHandler) SignIn(c *gin.Context) {
	c.HTML(http.StatusOK, "sessions/signin", h)
}

func (h sessionHandler) CreateSession(c *gin.Context) {
	c.Request.ParseForm()

	username := c.PostForm("email_or_username")
	password := c.PostForm("password")

	h.FormData = NewFormData()

	var user domain.User
	if err := persistence.DB.First(&user, "username = ? OR email = ?", username, username).Error; err != nil {
		h.FormData.Errors["email_or_username"] = "Unable to sign you in. Invalid credentials."
		c.HTML(http.StatusUnauthorized, "sessions/signin_form", h)
		return
	}

	token, err := auth.NewUserAuth().Login(username, password)

	if err != nil || token == "" {
		h.FormData.Errors["email_or_username"] = "Unable to sign you in. Invalid credentials."
		c.HTML(http.StatusUnauthorized, "sessions/signin_form", h)
		return
	}

	c.SetCookie("auth_token", token, 3_600*24*365, "/", "localhost:8080", true, true)
	c.Header("HX-Redirect", "/")
}

func (h sessionHandler) DestroySession(c *gin.Context) {
	c.SetCookie("auth_token", "", -1, "/", "localhost:8080", true, true)
	c.Header("HX-Redirect", "/")
}
