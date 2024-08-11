package service

import (
	"github.com/bytebury/fun-banking/internal/domain"
	"github.com/bytebury/fun-banking/internal/infrastructure/persistence"
)

type SubscriptionService interface {
	Create(subscription *domain.Subscription) error
	Update(subscriptionID string, subscription domain.Subscription) error
	FindBySubscriptionID(subscriptionID string, subscription *domain.Subscription) error
}

type subscriptionService struct{}

func NewSubscriptionService() SubscriptionService {
	return subscriptionService{}
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

func (s subscriptionService) FindBySubscriptionID(subscriptionID string, subscription *domain.Subscription) error {
	return persistence.DB.First(&subscription, "subscription_id = ?", subscriptionID).Error
}
