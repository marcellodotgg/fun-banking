package handler

import (
	"net/http"
	"os"
	"strings"

	"github.com/bytebury/fun-banking/internal/domain"
	"github.com/bytebury/fun-banking/internal/infrastructure/auth"
	"github.com/bytebury/fun-banking/internal/infrastructure/persistence"
	"github.com/bytebury/fun-banking/internal/service"
	"github.com/gin-gonic/gin"
)

type sessionHandler struct {
	pageObject
	Bank            domain.Bank
	customerService service.CustomerService
}

func NewSessionHandler() sessionHandler {
	return sessionHandler{
		Bank:            domain.Bank{},
		customerService: service.NewCustomerService(),
	}
}

func (h sessionHandler) SignIn(c *gin.Context) {
	h.Reset(c)
	c.HTML(http.StatusOK, "sessions/signin", h)
}

func (h sessionHandler) CreateSession(c *gin.Context) {
	h.Reset(c)

	username := strings.TrimSpace(strings.ToLower(h.Form.Data["email_or_username"]))
	password := h.Form.Data["password"]

	var user domain.User
	if err := persistence.DB.First(&user, "username = ? OR email = ?", username, username).Error; err != nil {
		h.Form.Errors["general"] = "Unable to sign you in. Invalid credentials."
		c.HTML(http.StatusUnauthorized, "sessions/signin_form", h)
		return
	}

	token, err := auth.NewUserAuth().Login(username, password)

	if err != nil && strings.Contains(err.Error(), "not verified") {
		h.Form.Errors["general"] = "You are not verified, please check your e-mail"
		c.HTML(http.StatusUnauthorized, "sessions/signin_form", h)
		return
	}

	if err != nil || token == "" {
		h.Form.Errors["general"] = "Unable to sign you in. Invalid credentials."
		c.HTML(http.StatusUnauthorized, "sessions/signin_form", h)
		return
	}

	c.SetCookie("auth_token", token, 3_600*24*365, "/", os.Getenv("COOKIE_URL"), true, true)
	c.Header("HX-Redirect", "/")
}

func (h sessionHandler) CreateCustomerSession(c *gin.Context) {
	h.Reset(c)

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

	c.SetCookie("customer_auth_token", token, 3_600*24*365, "/", os.Getenv("COOKIE_URL"), true, true)
	c.Header("HX-Redirect", "/")
}

func (h sessionHandler) DestroySession(c *gin.Context) {
	c.SetCookie("auth_token", "", -1, "/", os.Getenv("COOKIE_URL"), true, true)
	c.Header("HX-Redirect", "/")
}

func (h sessionHandler) DestroyCustomerSession(c *gin.Context) {
	c.SetCookie("customer_auth_token", "", -1, "/", os.Getenv("COOKIE_URL"), true, true)
	c.Header("HX-Redirect", "/")
}
