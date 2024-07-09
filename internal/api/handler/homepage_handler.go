package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type homepageHandler struct {
	PageObject struct{ Times *int }
}

func NewHomePageHandler() homepageHandler {
	return homepageHandler{
		PageObject: struct{ Times *int }{Times: new(int)},
	}
}

func (h homepageHandler) Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", h.PageObject)
}

func (h homepageHandler) Click(c *gin.Context) {
	*h.PageObject.Times++
	c.HTML(http.StatusOK, "click", h.PageObject)
}
