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
)

func ParseSubscriptionType(s string) SubscriptionType {
	switch s {
	case "aftermath-pro":
		return SubscriptionTypePro
	case "aftermath-pro-clan":
		return SubscriptionTypeProClan
	case "aftermath-plus":
		return SubscriptionTypePlus
	case "supporter":
		return SubscriptionTypeSupporter
	case "verified-clan":
		return SubscriptionTypeVerifiedClan
	default:
		return ""
	}
}

type UserSubscription struct {
	id          primitive.ObjectID      `bson:"_id,omitempty" json:"-"`
	UserID      string                  `bson:"userID" json:"userID"`
	ReferenceID string                  `bson:"referenceID" json:"referenceID"`
	Permissions permissions.Permissions `bson:"permissions" json:"permissions"`

	Type         SubscriptionType `bson:"subscriptionType" json:"subscriptionType"`
	ExpiryDate   time.Time        `bson:"expiryDate" json:"expiryDate"`
	CreationDate time.Time        `bson:"creationDate" json:"creationDate"`
}

func (s *UserSubscription) ID() primitive.ObjectID {
	return s.id
}

func (s *UserSubscription) IsExpired() bool {
	return s.ExpiryDate.Before(time.Now())
}
