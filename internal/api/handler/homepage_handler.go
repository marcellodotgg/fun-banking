package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/bytebury/fun-banking/internal/domain"
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
	SiteInfo    siteInfo
	SignedIn    bool
	Bank        domain.Bank
	bankService service.BankService
	Form        FormData
}

func NewHomePageHandler() homepageHandler {
	return homepageHandler{
		SiteInfo: siteInfo{
			UserCount:        0,
			CustomerCount:    0,
			BankCount:        0,
			TransactionCount: 0,
		},
		SignedIn:    false,
		Bank:        domain.Bank{},
		bankService: service.NewBankService(),
		Form:        NewFormData(),
	}
}

func (h homepageHandler) Homepage(c *gin.Context) {
	persistence.DB.Model(&domain.User{}).Count(&h.SiteInfo.UserCount)
	persistence.DB.Model(&domain.Customer{}).Count(&h.SiteInfo.CustomerCount)
	persistence.DB.Model(&domain.Bank{}).Count(&h.SiteInfo.BankCount)
	persistence.DB.Model(&domain.Transaction{}).Count(&h.SiteInfo.TransactionCount)

	h.SignedIn = c.GetString("user_id") != ""

	if h.SignedIn {
		h.Dashboard(c)
	} else {
		c.HTML(http.StatusOK, "index.html", h)
	}
}

func (h homepageHandler) Dashboard(c *gin.Context) {
	c.HTML(http.StatusOK, "dashboard", h)
}

func (h homepageHandler) SignUp(c *gin.Context) {
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

	h.Form.Data["bank_id"] = strconv.Itoa(int(h.Bank.ID))

	c.HTML(http.StatusOK, "customer_signin", h)
}
