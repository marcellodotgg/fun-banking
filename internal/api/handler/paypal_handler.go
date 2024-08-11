package handler

import (
	"fmt"
	"net/http"

	"github.com/bytebury/fun-banking/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/plutov/paypal/v4"
)

type payPalHandler struct {
	client        *paypal.Client
	payPalService service.PayPal
}

func NewPayPalHandler() payPalHandler {
	payPalService := service.PayPal{}
	client := payPalService.CreateClient()

	return payPalHandler{
		client,
		payPalService,
	}
}

func (p payPalHandler) Subcribe(c *gin.Context) {
	subscription, err := service.PayPal{}.CreateSubscription(p.client, "P-6D041016744961740M235G4Y")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, subscription)
}

func (p payPalHandler) HandleWebhook(c *gin.Context) {
	var webhookEvent paypal.WebhookEventType
	if err := c.BindJSON(&webhookEvent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println(webhookEvent)

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
