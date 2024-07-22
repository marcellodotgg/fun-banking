package handler

import (
	"net/http"

	"github.com/bytebury/fun-banking/internal/domain"
	"github.com/bytebury/fun-banking/internal/infrastructure/persistence"
	"github.com/gin-gonic/gin"
)

type controlPanelHandler struct {
	SignedIn bool
	SiteInfo siteInfo
}

func NewControlPanelHandler() controlPanelHandler {
	return controlPanelHandler{
		SignedIn: true,
		SiteInfo: siteInfo{},
	}
}

func (h controlPanelHandler) AppInsights(c *gin.Context) {
	persistence.DB.Model(&domain.User{}).Count(&h.SiteInfo.UserCount)
	persistence.DB.Model(&domain.Customer{}).Count(&h.SiteInfo.CustomerCount)
	persistence.DB.Model(&domain.Bank{}).Count(&h.SiteInfo.BankCount)
	persistence.DB.Model(&domain.Transaction{}).Count(&h.SiteInfo.TransactionCount)

	c.HTML(http.StatusOK, "control_panel_app_insights", h)
}
