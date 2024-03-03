package types

import (
	"time"

	"github.com/cufee/aftermath-core/internal/core/database/models"
)

type SubscriptionPayload struct {
	UserID      string `bson:"userID" json:"userID"`
	ReferenceID string `bson:"referenceID" json:"referenceID"`

	Type       models.SubscriptionType `bson:"subscriptionType" json:"subscriptionType"`
	ExpiryDate time.Time               `bson:"expiryDate" json:"expiryDate"`
}

func (s SubscriptionPayload) Valid() bool {
	return s.UserID != "" && s.ReferenceID != "" && s.Type != "" && !s.ExpiryDate.IsZero() && s.Type.Valid()
}

func (s SubscriptionPayload) ToUserSubscription() (models.UserSubscription, bool) {
	if !s.Valid() {
		return models.UserSubscription{}, false
	}

	return models.UserSubscription{
		Permissions:  s.Type.GetPermissions().Encode(),
		ReferenceID:  s.ReferenceID,
		ExpiryDate:   s.ExpiryDate,
		UserID:       s.UserID,
		Type:         s.Type,
		CreationDate: time.Now(),
	}, true
}

func (s *SubscriptionPayload) ToSubscriptionUpdate() models.SubscriptionUpdate {
	var sub models.SubscriptionUpdate
	if s.UserID != "" {
		sub.UserID = &s.UserID
	}
	if s.ReferenceID != "" {
		sub.ReferenceID = &s.ReferenceID
	}
	if !s.ExpiryDate.IsZero() {
		sub.ExpiryDate = &s.ExpiryDate
	}
	return sub
}

type ForceUpdateConnectionPayload struct {
	UserID         string `json:"userId"`
	ConnectionID   string `json:"connectionId"`
	ConnectionType string `json:"connectionType"`

	Metadata map[string]any `json:"metadata"`
}

func (update *ForceUpdateConnectionPayload) Type() (models.ConnectionType, bool) {
	switch update.ConnectionType {
	case string(models.ConnectionTypeWargaming):
		return models.ConnectionTypeWargaming, true
	default:
		return "", false
	}
}
