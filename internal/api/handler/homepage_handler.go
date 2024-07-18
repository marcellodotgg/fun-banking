package handler

import (
	"net/http"

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
	userService service.UserService
	SignedIn    bool
}

func NewHomePageHandler() homepageHandler {
	return homepageHandler{
		SiteInfo: siteInfo{
			UserCount:        0,
			CustomerCount:    0,
			BankCount:        0,
			TransactionCount: 0,
		},
		userService: service.NewUserService(),
		SignedIn:    false,
	}
}

func (h homepageHandler) Homepage(c *gin.Context) {
	h.SiteInfo.UserCount = h.userService.Count()
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
