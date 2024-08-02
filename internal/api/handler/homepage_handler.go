package handler

import (
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/bytebury/fun-banking/internal/domain"
	"github.com/bytebury/fun-banking/internal/infrastructure/mail"
	"github.com/bytebury/fun-banking/internal/infrastructure/persistence"
	"github.com/bytebury/fun-banking/internal/service"
	"github.com/gin-gonic/gin"
)

type siteInfo struct {
	UserCount        int64
	CustomerCount    int64
	BankCount        int64
	TransactionCount int64
}

type homepageHandler struct {
	pageObject
	SiteInfo        siteInfo
	Bank            domain.Bank
	Customer        domain.Customer
	bankService     service.BankService
	customerService service.CustomerService
	userService     service.UserService
	tokenService    service.TokenService
}

func NewHomePageHandler() homepageHandler {
	return homepageHandler{
		SiteInfo: siteInfo{
			UserCount:        0,
			CustomerCount:    0,
			BankCount:        0,
			TransactionCount: 0,
		},
		Bank:            domain.Bank{},
		Customer:        domain.Customer{},
		bankService:     service.NewBankService(),
		customerService: service.NewCustomerService(),
		tokenService:    service.NewTokenService(),
		userService:     service.NewUserService(),
	}
}

func (h homepageHandler) Homepage(c *gin.Context) {
	h.Reset(c)

	persistence.DB.Model(&domain.User{}).Count(&h.SiteInfo.UserCount)
	persistence.DB.Model(&domain.Customer{}).Count(&h.SiteInfo.CustomerCount)
	persistence.DB.Model(&domain.Bank{}).Count(&h.SiteInfo.BankCount)
	persistence.DB.Model(&domain.Transaction{}).Count(&h.SiteInfo.TransactionCount)

	if c.GetString("customer_id") != "" {
		h.CustomerDashboard(c)
		return
	}

	if h.SignedIn {
		h.Dashboard(c)
	} else {
		c.HTML(http.StatusOK, "index.html", h)
	}
}

func (h homepageHandler) TermsOfService(c *gin.Context) {
	h.Reset(c)
	c.HTML(http.StatusOK, "terms", h)
}

func (h homepageHandler) PrivacyPolicy(c *gin.Context) {
	h.Reset(c)
	c.HTML(http.StatusOK, "privacy", h)
}

func (h homepageHandler) VerifyEmail(c *gin.Context) {
	h.Reset(c)
	userID, err := h.tokenService.GetUserIDFromToken(c.Query("token"))
	h.Form = NewFormData()

	if err != nil {
		h.Form.Errors["general"] = "Token was invalid or has been expired"
		c.HTML(http.StatusUnprocessableEntity, "resend_verification", h)
		return
	}

	if err := h.userService.Update(userID, &domain.User{Verified: true}); err != nil {
		h.Form.Errors["general"] = "We were unable to verify your account. Please try again later"
		c.HTML(http.StatusUnprocessableEntity, "resend_verification", h)
		return
	}

	token, _ := h.tokenService.GenerateUserToken(userID)

	c.SetCookie("auth_token", token, 3_600*24*365, "/", os.Getenv("COOKIE_URL"), true, true)
	c.HTML(http.StatusOK, "dashboard", h)
}

func (h homepageHandler) ResendVerifyEmail(c *gin.Context) {
	h.Reset(c)

	var user domain.User
	if err := h.userService.FindByEmail(h.Form.Data["email"], &user); err != nil {
		h.Form.Data["success"] = "We have sent an e-mail to that account if it exists."
		c.HTML(http.StatusUnprocessableEntity, "resend_verification_form", h.Form)
		return
	}

	if err := mail.NewWelcomeMailer().Send(user.Email, user); err != nil {
		h.Form.Errors["general"] = "Email service is down, unable to send. Please try again later."
		c.HTML(http.StatusUnprocessableEntity, "resend_verification_form", h.Form)
		return
	}

	h.Form.Data["success"] = "We have sent an e-mail to that account if it exists."
	c.HTML(http.StatusUnprocessableEntity, "resend_verification_form", h.Form)
}

func (h homepageHandler) Dashboard(c *gin.Context) {
	h.Reset(c)
	c.HTML(http.StatusOK, "dashboard", h)
}

func (h homepageHandler) CustomerDashboard(c *gin.Context) {
	h.Theme = c.GetString("theme")
	customerID := c.GetString("customer_id")

	if err := h.customerService.FindByID(customerID, &h.Customer); err != nil {
		c.HTML(http.StatusNotFound, "index.html", h)
		return
	}

	c.HTML(http.StatusOK, "customer", h)
}

func (h homepageHandler) SignUp(c *gin.Context) {
	h.Theme = c.GetString("theme")
	c.HTML(http.StatusOK, "signup.html", h)
}

func (h homepageHandler) BankSignIn(c *gin.Context) {
	h.Form = NewFormData()
	username := strings.ToLower(c.Param("username"))
	bankSlug := strings.ToLower(c.Param("slug"))

	if err := h.bankService.FindByUsernameAndSlug(username, bankSlug, &h.Bank); err != nil {
		c.HTML(http.StatusNotFound, "not-found", nil)
		return
	}

	h.Theme = h.Bank.User.Theme
	h.Form.Data["bank_id"] = strconv.Itoa(int(h.Bank.ID))

	c.HTML(http.StatusOK, "customer_signin", h)
}
