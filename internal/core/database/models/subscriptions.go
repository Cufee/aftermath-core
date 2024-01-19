package models

import (
	"time"

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

type UserSubscription struct {
	id          primitive.ObjectID `bson:"_id,omitempty" json:"-"`
	UserID      string             `bson:"userID" json:"userID"`
	ReferenceID *string            `bson:"referenceID" json:"referenceID"`

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
