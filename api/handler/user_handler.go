package handler

import (
	"net/http"

	"github.com/bytebury/fun-banking/internal/domain"
	"github.com/bytebury/fun-banking/internal/infrastructure/persistence"
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

func (h userHandler) SignUp(c *gin.Context) {
	c.HTML(http.StatusOK, "users/signup", nil)
}

func (h userHandler) Create(c *gin.Context) {
	c.Request.ParseForm()

	formData := newFormData()

	for key, values := range c.Request.PostForm {
		if len(values) > 0 {
			formData.Data[key] = values[0]
		}
	}

	if formData.Data["password"] != formData.Data["password_confirmation"] {
		formData.Errors["passwords_dont_match"] = "Passwords provided do not match"
		c.HTML(http.StatusUnprocessableEntity, "users/signup_form", formData)
		return
	}

	if err := persistence.DB.Create(&domain.User{
		FirstName: formData.Data["first_name"],
		LastName:  formData.Data["last_name"],
		Email:     formData.Data["email"],
		Username:  formData.Data["username"],
	}).Error; err != nil {
		formData.Errors["unable_to_create"] = "Something went wrong creating your account"
		c.HTML(http.StatusUnprocessableEntity, "users/signup_form", formData)
		return
	}

	c.Header("HX-Redirect", "/")
}
