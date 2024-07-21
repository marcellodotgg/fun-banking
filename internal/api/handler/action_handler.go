package handler

import (
	"net/http"

	"github.com/bytebury/fun-banking/internal/domain"
	"github.com/bytebury/fun-banking/internal/service"
	"github.com/gin-gonic/gin"
)

type actionHandler struct {
	SignedIn    bool
	CurrentUser domain.User
	userService service.UserService
}

func NewActionHandler() actionHandler {
	return actionHandler{
		SignedIn:    false,
		CurrentUser: domain.User{},
		userService: service.NewUserService(),
	}
}

func (h actionHandler) OpenAppDrawer(c *gin.Context) {
	userID := c.GetString("user_id")

	h.SignedIn = userID != ""

	if h.SignedIn {
		h.userService.FindByID(userID, &h.CurrentUser)
	}

	c.HTML(http.StatusOK, "app_drawer", h)
}
