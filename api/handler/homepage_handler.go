package handler

import (
	"net/http"

	"github.com/bytebury/fun-banking/internal/service"
	"github.com/gin-gonic/gin"
)

type siteInfo struct {
	UserCount int64
}

type homepageHandler struct {
	SiteInfo    siteInfo
	userService service.UserService
}

func NewHomePageHandler() homepageHandler {
	return homepageHandler{
		SiteInfo: siteInfo{
			UserCount: 0,
		},
		userService: service.NewUserService(),
	}
}

func (h homepageHandler) Index(c *gin.Context) {
	h.SiteInfo.UserCount = h.userService.Count()
	c.HTML(http.StatusOK, "index.html", h)
}

func (h homepageHandler) Click(c *gin.Context) {
	c.HTML(http.StatusOK, "click", h)
}
