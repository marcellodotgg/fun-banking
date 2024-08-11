package domain

type Subscription struct {
	Audit
	User           User   `gorm:"foreignKey:UserID; constraint:OnDeleteCASCADE"`
	UserID         string `gorm:"not null;uniqueIndex:idx_user_sub"`
	PlanID         string
	SubscriptionID string `gorm:"uniqueIndex:idx_user_sub"`
	Status         string
}
