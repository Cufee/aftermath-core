package models

import (
	"time"

	"github.com/cufee/aftermath-core/permissions/v1"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SubscriptionType string

const (
	SubscriptionTypePro          = SubscriptionType("aftermath-pro")
	SubscriptionTypeProClan      = SubscriptionType("aftermath-pro-clan")
	SubscriptionTypePlus         = SubscriptionType("aftermath-plus")
	SubscriptionTypeSupporter    = SubscriptionType("supporter")
	SubscriptionTypeVerifiedClan = SubscriptionType("verified-clan")

	SubscriptionTypeServerBooster    = SubscriptionType("server-booster")
	SubscriptionTypeServerModerator  = SubscriptionType("server-moderator")
	SubscriptionTypeContentModerator = SubscriptionType("content-moderator")
)

func (s SubscriptionType) Valid() bool {
	switch s {
	case SubscriptionTypePro, SubscriptionTypeProClan, SubscriptionTypePlus, SubscriptionTypeSupporter, SubscriptionTypeVerifiedClan, SubscriptionTypeServerBooster, SubscriptionTypeServerModerator, SubscriptionTypeContentModerator:
		return true
	default:
		return false
	}
}

func ParseSubscriptionType(s string) (SubscriptionType, bool) {
	switch s {
	case "aftermath-pro":
		return SubscriptionTypePro, true
	case "aftermath-pro-clan":
		return SubscriptionTypeProClan, true
	case "aftermath-plus":
		return SubscriptionTypePlus, true
	case "supporter":
		return SubscriptionTypeSupporter, true
	case "verified-clan":
		return SubscriptionTypeVerifiedClan, true
	case "server-booster":
		return SubscriptionTypeServerBooster, true
	case "server-moderator":
		return SubscriptionTypeServerModerator, true
	case "content-moderator":
		return SubscriptionTypeContentModerator, true
	default:
		return "", false
	}
}

type UserSubscription struct {
	ID          primitive.ObjectID      `bson:"_id,omitempty" json:"id"`
	UserID      string                  `bson:"userID" json:"userID"`
	ReferenceID string                  `bson:"referenceID" json:"referenceID"`
	Permissions permissions.Permissions `bson:"permissions" json:"permissions"`

	Type         SubscriptionType `bson:"subscriptionType" json:"subscriptionType"`
	ExpiryDate   time.Time        `bson:"expiryDate" json:"expiryDate"`
	CreationDate time.Time        `bson:"creationDate" json:"creationDate"`
}

type SubscriptionUpdate struct {
	UserID      *string                  `bson:"userID,omitempty" json:"userID"`
	ReferenceID *string                  `bson:"referenceID,omitempty" json:"referenceID"`
	Permissions *permissions.Permissions `bson:"permissions,omitempty" json:"permissions"`

	Type       *SubscriptionType `bson:"subscriptionType,omitempty" json:"subscriptionType"`
	ExpiryDate *time.Time        `bson:"expiryDate,omitempty" json:"expiryDate"`
}

func (s *UserSubscription) IsExpired() bool {
	return s.ExpiryDate.Before(time.Now())
}
