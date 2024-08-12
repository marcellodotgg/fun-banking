package service

import (
	"github.com/bytebury/fun-banking/internal/domain"
	"github.com/bytebury/fun-banking/internal/infrastructure/persistence"
)

type SubscriptionService interface {
	Create(subscription *domain.Subscription) error
	Update(subscriptionID string, subscription domain.Subscription) error
	Cancel(id string) error
	FindByID(id string, subscription *domain.Subscription) error
	FindBySubscriptionID(subscriptionID string, subscription *domain.Subscription) error
}

type subscriptionService struct {
	paypalService paypalService
}

func NewSubscriptionService() SubscriptionService {
	return subscriptionService{
		paypalService: NewPayPalService(),
	}
}

func (s subscriptionService) Create(subscription *domain.Subscription) error {
	return persistence.DB.Create(&subscription).Error
}

func (s subscriptionService) Update(subscriptionID string, subscription domain.Subscription) error {
	var currentSubscription domain.Subscription
	if err := s.FindBySubscriptionID(subscriptionID, &currentSubscription); err != nil {
		return err
	}
	return persistence.DB.Model(&currentSubscription).Updates(&subscription).Error
}

func (s subscriptionService) Cancel(id string) error {
	var subscription domain.Subscription
	if err := s.FindByID(id, &subscription); err != nil {
		return err
	}
	if err := s.paypalService.CancelSubscription(subscription.SubscriptionID); err != nil {
		return err
	}

	subscription.Status = "CANCELLED"

	return s.Update(subscription.SubscriptionID, subscription)
}

func (s subscriptionService) FindByID(subscriptionID string, subscription *domain.Subscription) error {
	return persistence.DB.First(&subscription, "id = ?", subscriptionID).Error
}

func (s subscriptionService) FindBySubscriptionID(subscriptionID string, subscription *domain.Subscription) error {
	return persistence.DB.First(&subscription, "subscription_id = ?", subscriptionID).Error
}
