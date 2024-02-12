package types

import (
	"time"

	"github.com/cufee/aftermath-core/internal/core/database/models"
	"github.com/cufee/aftermath-core/permissions/v2"
)

type SubscriptionPayload struct {
	UserID      string                  `bson:"userID" json:"userID"`
	ReferenceID string                  `bson:"referenceID" json:"referenceID"`
	Permissions permissions.Permissions `bson:"permissions" json:"permissions"`

	Type       models.SubscriptionType `bson:"subscriptionType" json:"subscriptionType"`
	ExpiryDate time.Time               `bson:"expiryDate" json:"expiryDate"`
}

func (s SubscriptionPayload) Valid() bool {
	return s.UserID != "" && s.ReferenceID != "" && s.Type != "" && !s.ExpiryDate.IsZero() && s.Type.Valid()
}

func (s SubscriptionPayload) ToUserSubscription() models.UserSubscription {
	return models.UserSubscription{
		UserID:       s.UserID,
		ReferenceID:  s.ReferenceID,
		Permissions:  s.Permissions.Encode(),
		Type:         s.Type,
		ExpiryDate:   s.ExpiryDate,
		CreationDate: time.Now(),
	}
}

func (s *SubscriptionPayload) ToSubscriptionUpdate() models.SubscriptionUpdate {
	var sub models.SubscriptionUpdate
	if s.UserID != "" {
		sub.UserID = &s.UserID
	}
	if s.ReferenceID != "" {
		sub.ReferenceID = &s.ReferenceID
	}
	if s.Permissions != permissions.Blank {
		p := s.Permissions.Encode()
		sub.Permissions = &p
	}
	if s.Type != "" {
		sub.Type = &s.Type
	}
	if !s.ExpiryDate.IsZero() {
		sub.ExpiryDate = &s.ExpiryDate
	}
	return sub
}
