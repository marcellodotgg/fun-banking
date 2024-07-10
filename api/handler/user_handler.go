package handler

import (
	"net/http"

	"github.com/bytebury/fun-banking/internal/service"
	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService service.UserService
}

func NewUserHandler() userHandler {
	return userHandler{
		userService: service.NewUserService(),
	}
}

func (h userHandler) Count(c *gin.Context) {
	c.HTML(http.StatusOK, "users_count", h.userService.Count())
}
