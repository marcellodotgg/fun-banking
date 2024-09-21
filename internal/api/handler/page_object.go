package handler

import (
	"os"

	"github.com/gin-gonic/gin"
)

type pageObject struct {
	Form              FormData
	SignedIn          bool
	Theme             string
	Hash              string
	ModalType         string
	CurrentUserID     string
	GoogleAnalyticsID string
}

func (po *pageObject) Reset(c *gin.Context) {
	po.Form = GetForm(c)
	po.SignedIn = c.GetString("user_id") != ""
	po.CurrentUserID = c.GetString("user_id")
	po.Theme = c.GetString("theme")
	po.GoogleAnalyticsID = os.Getenv("GOOGLE_ANALYTICS_ID")

	if os.Getenv("GIN_MODE") == "release" {
		po.Hash = os.Getenv("BUILD_HASH")
	} else {
		po.Hash = "local"
	}

	po.ModalType = ""
}
