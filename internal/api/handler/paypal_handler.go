package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/bytebury/fun-banking/internal/domain"
	"github.com/bytebury/fun-banking/internal/service"
	"github.com/gin-gonic/gin"
)

type payPalHandler struct {
	subscriptionService service.SubscriptionService
}

type payPalResource struct {
	ID       string `json:"id"`
	CustomID string `json:"custom_id"`
	PlanID   string `json:"plan_id"`
	Status   string `json:"status"`
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
	return payPalHandler{
		subscriptionService: service.NewSubscriptionService(),
	}
}

func (p payPalHandler) HandleWebhook(c *gin.Context) {
	var webhook webhook
	if err := c.BindJSON(&webhook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var subscription domain.Subscription
	if strings.Contains(webhook.EventType, "BILLING.SUBSCRIPTION") {
		subscription = domain.Subscription{
			UserID:         webhook.Resource.CustomID,
			PlanID:         webhook.Resource.PlanID,
			SubscriptionID: webhook.Resource.ID,
			Status:         webhook.Resource.Status,
		}
	}

	switch webhook.EventType {
	case "BILLING.SUBSCRIPTION.ACTIVATED":
		fmt.Println("user", webhook.Resource.CustomID, "activated their subscription")
		p.subscriptionService.Update(webhook.Resource.ID, subscription)
	case "BILLING.SUBSCRIPTION.CREATED":
		fmt.Println("user", webhook.Resource.CustomID, "created their subscription")
		p.subscriptionService.Create(&subscription)
	case "BILLING.SUBSCRIPTION.CANCELLED":
		fmt.Println("user", webhook.Resource.CustomID, "cancelled their subscription")
		p.subscriptionService.Update(webhook.Resource.ID, subscription)
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
