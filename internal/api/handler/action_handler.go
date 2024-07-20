package handler

import (
	"net/http"

	"github.com/bytebury/fun-banking/internal/domain"
	"github.com/gin-gonic/gin"
)

type actionHandler struct {
	SignedIn    bool
	CurrentUser domain.User
}

func NewActionHandler() actionHandler {
	return actionHandler{
		SignedIn:    false,
		CurrentUser: domain.User{},
	}
}

func (a actionHandler) OpenAppDrawer(c *gin.Context) {
	userID := c.GetString("user_id")

	a.SignedIn = userID != ""

	c.HTML(http.StatusOK, "app_drawer", a)
}
