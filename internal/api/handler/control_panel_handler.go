package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type controlPanelHandler struct {
	SignedIn bool
}

func NewControlPanelHandler() controlPanelHandler {
	return controlPanelHandler{
		SignedIn: true,
	}
}

func (h controlPanelHandler) AppInsights(c *gin.Context) {
	c.HTML(http.StatusOK, "control_panel_app_insights", h)
}
