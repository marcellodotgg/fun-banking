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
	}
}

func (h homepageHandler) Homepage(c *gin.Context) {
	h.SiteInfo.UserCount = h.userService.Count()
	c.HTML(http.StatusOK, "index.html", h)
}

func (h homepageHandler) SignIn(c *gin.Context) {
	c.HTML(http.StatusOK, "signin.html", h)
}

func (h homepageHandler) CreateSession(c *gin.Context) {
	c.Request.ParseForm()

	emailOrUsername := c.PostForm("email_or_username")
	password := c.PostForm("password")

	if emailOrUsername == "marcello" && password == "password" {
		c.Header("HX-Redirect", "/")
	}

	c.HTML(http.StatusUnauthorized, "index", nil)
}

func (h homepageHandler) SignUp(c *gin.Context) {
	c.HTML(http.StatusOK, "signup.html", h)
}
