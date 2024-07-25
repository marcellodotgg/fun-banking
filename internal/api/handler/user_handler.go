package handler

import (
	"net/http"
	"strings"

	"github.com/bytebury/fun-banking/internal/domain"
	"github.com/bytebury/fun-banking/internal/infrastructure/mail"
	"github.com/bytebury/fun-banking/internal/service"
	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService service.UserService
	Form        FormData
	SignedIn    bool
	User        domain.User
}

func NewUserHandler() userHandler {
	return userHandler{
		userService: service.NewUserService(),
		Form:        NewFormData(),
		SignedIn:    false,
		User:        domain.User{},
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

	h.Form = NewFormData()
	h.Form.Data["success"] = "Successfully created your account. We have sent you an e-mail with instructions to activate your account. If you do not see an e-mail, please check your spam folder. If you still do not see an e-mail after 5 minutes, please send us an e-mail at bytebury@gmail.com"
	c.HTML(http.StatusAccepted, "users/signup_form", h)
}

func (h userHandler) Settings(c *gin.Context) {
	userID := c.GetString("user_id")

	if err := h.userService.FindByID(userID, &h.User); err != nil {
		c.HTML(http.StatusNotFound, "not-found", nil)
		return
	}

	h.Form = NewFormData()
	h.Form.Data["first_name"] = h.User.FirstName
	h.Form.Data["last_name"] = h.User.LastName
	h.Form.Data["username"] = h.User.Username
	h.Form.Data["image_url"] = h.User.ImageURL

	c.HTML(http.StatusOK, "user_settings", h)
}

func (h userHandler) Update(c *gin.Context) {
	userID := c.GetString("user_id")

	if err := h.userService.FindByID(userID, &h.User); err != nil {
		c.HTML(http.StatusNotFound, "not-found", nil)
		return
	}

	h.Form = GetForm(c)
	h.User.FirstName = h.Form.Data["first_name"]
	h.User.LastName = h.Form.Data["last_name"]
	h.User.Username = h.Form.Data["username"]
	h.User.ImageURL = h.Form.Data["image_url"]

	if err := h.userService.Update(userID, &h.User); err != nil {
		h.Form.Errors["general"] = "Something went wrong updating your account settings"
		c.HTML(http.StatusUnprocessableEntity, "user_settings_form", h)
		return
	}

	h.Form.Data["success"] = "Successfully updated your account settings"

	// TODO: Might need to be OOB
	c.HTML(http.StatusAccepted, "user_settings_form", h)
}

func (h userHandler) ForgotPassword(c *gin.Context) {
	h.Form = NewFormData()
	c.HTML(http.StatusOK, "forgot_password", h)
}

func (h userHandler) SendForgotPasswordEmail(c *gin.Context) {
	h.Form = GetForm(c)

	if err := h.userService.FindByEmail(h.Form.Data["email"], &h.User); err != nil {
		h.Form.Data["success"] = "Sent password reset instructions to that e-mail if it exists"
		c.HTML(http.StatusOK, "forgot_password_form", h.Form)
		return
	}

	if err := mail.NewPasswordResetMailer().Send(h.Form.Data["email"], h.User); err != nil {
		h.Form.Errors["general"] = "Our e-mail service is down, please try again later."
		c.HTML(http.StatusUnprocessableEntity, "forgot_password_form", h.Form)
		return
	}

	h.Form.Data["success"] = "Sent password reset instructions to that e-mail if it exists"
	c.HTML(http.StatusOK, "forgot_password_form", h.Form)
}

func (h userHandler) Notifications(c *gin.Context) {
	h.SignedIn = true
	c.HTML(http.StatusOK, "notifications", h)
}

func (h userHandler) PendingTransactions(c *gin.Context) {
	userID := c.GetString("user_id")

	var transactions []domain.Transaction
	h.userService.FindPendingTransactions(userID, &transactions)

	c.HTML(http.StatusOK, "notifications_list", transactions)
}

func (h userHandler) HasPendingTransactions(c *gin.Context) {
	userID := c.GetString("user_id")

	var transactions []domain.Transaction
	if err := h.userService.FindPendingTransactions(userID, &transactions); err != nil {
		c.HTML(http.StatusOK, "inbox_badge", false)
		return
	}

	c.HTML(http.StatusOK, "inbox_badge", len(transactions))
}
