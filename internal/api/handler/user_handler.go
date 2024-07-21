package handler

import (
	"net/http"
	"strings"

	"github.com/bytebury/fun-banking/internal/domain"
	"github.com/bytebury/fun-banking/internal/service"
	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService service.UserService
	Form        FormData
	SignedIn    bool
}

func NewUserHandler() userHandler {
	return userHandler{
		userService: service.NewUserService(),
		Form:        NewFormData(),
		SignedIn:    false,
	}
}

func (h userHandler) SignUp(c *gin.Context) {
	c.HTML(http.StatusOK, "users/signup", h)
}

func (h userHandler) Create(c *gin.Context) {
	c.Request.ParseForm()

	h.Form = NewFormData()

	for key, values := range c.Request.PostForm {
		if len(values) > 0 {
			h.Form.Data[key] = values[0]
		}
	}

	if h.Form.Data["password"] != h.Form.Data["password_confirmation"] {
		h.Form.Errors["general"] = "Passwords provided do not match"
		h.Form.Errors["passwords_dont_match"] = "Passwords provided do not match"
		c.HTML(http.StatusUnprocessableEntity, "users/signup_form", h)
		return
	}

	user := domain.User{
		FirstName: h.Form.Data["first_name"],
		LastName:  h.Form.Data["last_name"],
		Email:     h.Form.Data["email"],
		Username:  h.Form.Data["username"],
		Password:  h.Form.Data["password"],
	}
	if err := h.userService.Create(&user); err != nil {
		if strings.Contains(err.Error(), "UNIQUE") {
			h.Form.Errors["general"] = "An account with that username or e-mail already exists"
		} else {
			h.Form.Errors["general"] = "Something went wrong creating your account. If the problem persists, please contact us."
		}
		c.HTML(http.StatusUnprocessableEntity, "users/signup_form", h)
		return
	}

	c.Header("HX-Redirect", "/signin")
}

func (h userHandler) Settings(c *gin.Context) {
	c.HTML(http.StatusOK, "user_settings", h)
}
