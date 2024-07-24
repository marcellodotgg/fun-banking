package handler

import (
	"net/http"

	"github.com/bytebury/fun-banking/internal/domain"
	"github.com/bytebury/fun-banking/internal/infrastructure/auth"
	"github.com/bytebury/fun-banking/internal/infrastructure/persistence"
	"github.com/bytebury/fun-banking/internal/service"
	"github.com/gin-gonic/gin"
)

type sessionHandler struct {
	Form            FormData
	Bank            domain.Bank
	SignedIn        bool
	customerService service.CustomerService
}

func NewSessionHandler() sessionHandler {
	return sessionHandler{
		Form:            NewFormData(),
		Bank:            domain.Bank{},
		SignedIn:        false,
		customerService: service.NewCustomerService(),
	}
}

func (h sessionHandler) SignIn(c *gin.Context) {
	h.Form = NewFormData()
	c.HTML(http.StatusOK, "sessions/signin", h)
}

func (h sessionHandler) CreateSession(c *gin.Context) {
	c.Request.ParseForm()

	username := c.PostForm("email_or_username")
	password := c.PostForm("password")

	h.Form = NewFormData()

	var user domain.User
	if err := persistence.DB.First(&user, "username = ? OR email = ?", username, username).Error; err != nil {
		h.Form.Errors["email_or_username"] = "Unable to sign you in. Invalid credentials."
		c.HTML(http.StatusUnauthorized, "sessions/signin_form", h)
		return
	}

	token, err := auth.NewUserAuth().Login(username, password)

	if err != nil || token == "" {
		h.Form.Errors["email_or_username"] = "Unable to sign you in. Invalid credentials."
		c.HTML(http.StatusUnauthorized, "sessions/signin_form", h)
		return
	}

	c.SetCookie("auth_token", token, 3_600*24*365, "/", "localhost:8080", true, true)
	c.Header("HX-Redirect", "/")
}

func (h sessionHandler) CreateCustomerSession(c *gin.Context) {
	h.Form = GetForm(c)

	var customer domain.Customer
	if err := h.customerService.FindByBankIDAndPIN(h.Form.Data["bank_id"], h.Form.Data["pin"], &customer); err != nil {
		h.Form.Errors["general"] = "Unable to sign you in, invalid credentials"
		c.HTML(http.StatusUnauthorized, "customer_signin_form", h)
		return
	}

	token, err := auth.NewCustomerAuth().Login(customer)

	if err != nil || token == "" {
		h.Form.Errors["general"] = "Unable to sign you in, invalid credentials"
		c.HTML(http.StatusUnauthorized, "customer_signin_form", h)
		return
	}

	c.SetCookie("customer_auth_token", token, 3_600*24*365, "/", "localhost:8080", true, true)
	c.Header("HX-Redirect", "/")
}

func (h sessionHandler) DestroySession(c *gin.Context) {
	c.SetCookie("auth_token", "", -1, "/", "localhost:8080", true, true)
	c.Header("HX-Redirect", "/")
}

func (h sessionHandler) DestroyCustomerSession(c *gin.Context) {
	c.SetCookie("customer_auth_token", "", -1, "/", "localhost:8080", true, true)
	c.Header("HX-Redirect", "/")
}
