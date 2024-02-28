package models

import (
	"slices"
	"time"

	"github.com/cufee/aftermath-core/permissions/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SubscriptionType string

func (s SubscriptionType) GetPermissions() permissions.Permissions {
	switch s {
	case SubscriptionTypePlus:
		return permissions.SubscriptionAftermathPlus
	case SubscriptionTypePro:
		return permissions.SubscriptionAftermathPro
	case SubscriptionTypeProClan:
		return permissions.SubscriptionAftermathPro
	default:
		return permissions.User
	}
}

// Paid
const SubscriptionTypePro = SubscriptionType("aftermath-pro")
const SubscriptionTypeProClan = SubscriptionType("aftermath-pro-clan")
const SubscriptionTypePlus = SubscriptionType("aftermath-plus")

// Misc
const SubscriptionTypeSupporter = SubscriptionType("supporter")
const SubscriptionTypeVerifiedClan = SubscriptionType("verified-clan")

// Moderators
const SubscriptionTypeServerModerator = SubscriptionType("server-moderator")
const SubscriptionTypeContentModerator = SubscriptionType("content-moderator")

// Special
const SubscriptionTypeDeveloper = SubscriptionType("developer")
const SubscriptionTypeServerBooster = SubscriptionType("server-booster")
const SubscriptionTypeContentTranslator = SubscriptionType("content-translator")

var AllSubscriptionTypes = []SubscriptionType{
	SubscriptionTypePro,
	SubscriptionTypeProClan,
	SubscriptionTypePlus,
	SubscriptionTypeSupporter,
	SubscriptionTypeVerifiedClan,
	SubscriptionTypeServerModerator,
	SubscriptionTypeContentModerator,
	SubscriptionTypeDeveloper,
	SubscriptionTypeServerBooster,
	SubscriptionTypeContentTranslator,
}

func (s SubscriptionType) Valid() bool {
	return slices.Contains(AllSubscriptionTypes, s)
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
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID      string             `bson:"userID" json:"userID"`
	ReferenceID string             `bson:"referenceID" json:"referenceID"`
	Permissions string             `bson:"permissions" json:"permissions"`

	Type         SubscriptionType `bson:"subscriptionType" json:"subscriptionType"`
	ExpiryDate   time.Time        `bson:"expiryDate" json:"expiryDate"`
	CreationDate time.Time        `bson:"creationDate" json:"creationDate"`
}

type SubscriptionUpdate struct {
	UserID      *string `bson:"userID,omitempty" json:"userID"`
	ReferenceID *string `bson:"referenceID,omitempty" json:"referenceID"`
	Permissions *string `bson:"permissions,omitempty" json:"permissions"`

	Type       *SubscriptionType `bson:"subscriptionType,omitempty" json:"subscriptionType"`
	ExpiryDate *time.Time        `bson:"expiryDate,omitempty" json:"expiryDate"`
}

func (s *UserSubscription) IsExpired() bool {
	return s.ExpiryDate.Before(time.Now())
}
