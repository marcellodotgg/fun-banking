package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type payPalHandler struct {
}

type payPalResource struct {
	CustomID string `json:"custom_id"`
}

type webhook struct {
	ID           string         `json:"id"`
	CreateTime   string         `json:"create_time"`
	ResourceType string         `json:"resource_type"`
	EventType    string         `json:"event_type"`
	Summary      string         `json:"summary"`
	Resource     payPalResource `json:"resource"`
}

func NewPayPalHandler() payPalHandler {
	return payPalHandler{}
}

func (p payPalHandler) HandleWebhook(c *gin.Context) {
	var webhook webhook
	if err := c.BindJSON(&webhook); err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println(webhook)

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
