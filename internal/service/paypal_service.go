package service

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/plutov/paypal/v4"
)

type paypalService struct {
	ClientID string
	secret   string
}

func NewPayPalService() paypalService {
	return paypalService{
		ClientID: os.Getenv("PAYPAL_CLIENT_ID"),
		secret:   os.Getenv("PAYPAL_SECRET"),
	}
}

func (p paypalService) CancelSubscription(subscriptionID string) error {
	client, err := paypal.NewClient(p.ClientID, p.secret, paypal.APIBaseSandBox)
	if err != nil {
		return err
	}

	// Create an access token
	_, err2 := client.GetAccessToken(context.Background())
	if err2 != nil {
		return err
	}

	// Cancel the subscription
	err = client.CancelSubscription(context.Background(), subscriptionID, "Customer request cancellation")
	if err != nil {
		fmt.Println(err.Error())
		if strings.Contains(err.Error(), "SUBSCRIPTION_STATUS_INVALID") {
			return nil
		}
		return err
	}
	return nil
}
