package service

import (
	"log"
	"os"

	"github.com/plutov/paypal/v4"
)

type PayPal struct{}

func (p PayPal) CreateClient() *paypal.Client {
	apiBase := paypal.APIBaseSandBox

	// if os.Getenv("GIN_MODE") == "release" {
	// 	apiBase = paypal.APIBaseLive
	// }

	client, err := paypal.NewClient(os.Getenv("PAYPAL_CLIENT_ID"), os.Getenv("PAYPAL_SECRET"), apiBase)

	if err != nil {
		log.Fatal(err)
	}
	return client
}

func (p PayPal) CreateSubscription(client *paypal.Client, planID string) (*paypal.Subscription, error) {
	return nil, nil
}
