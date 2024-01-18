package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type subscriptionType string

const (
	SubscriptionTypePro          = subscriptionType("aftermath-pro")
	SubscriptionTypeProClan      = subscriptionType("aftermath-pro-clan")
	SubscriptionTypePlus         = subscriptionType("aftermath-plus")
	SubscriptionTypeSupporter    = subscriptionType("supporter")
	SubscriptionTypeVerifiedClan = subscriptionType("verified-clan")
)

type UserSubscription struct {
	id          primitive.ObjectID `bson:"_id,omitempty" json:"-"`
	UserID      string             `bson:"userID" json:"userID"`
	ReferenceID *string            `bson:"referenceID" json:"referenceID"`

	Type         subscriptionType `bson:"subscriptionType" json:"subscriptionType"`
	ExpiryDate   time.Time        `bson:"expiryDate" json:"expiryDate"`
	CreationDate time.Time        `bson:"creationDate" json:"creationDate"`
}

func (s *UserSubscription) ID() primitive.ObjectID {
	return s.id
}

func (s *UserSubscription) IsExpired() bool {
	return s.ExpiryDate.Before(time.Now())
}
