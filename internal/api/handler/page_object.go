package handler

import (
	"os"

	"github.com/gin-gonic/gin"
)

type pageObject struct {
	Form      FormData
	SignedIn  bool
	Theme     string
	Hash      string
	ModalType string
}

func (po *pageObject) Reset(c *gin.Context) {
	po.Form = GetForm(c)
	po.SignedIn = c.GetString("user_id") != ""
	po.Theme = c.GetString("theme")
	po.Hash = os.Getenv("BUILD_HASH")
	po.ModalType = ""
}
