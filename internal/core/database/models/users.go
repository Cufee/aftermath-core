package models

import "github.com/cufee/aftermath-core/permissions/v1"

type featureFlag string

const (
	FeatureFlagCustomizationDisabled = featureFlag("customizationDisabled")
)

type User struct {
	ID string `bson:"_id" json:"id"`

	FeatureFlags []featureFlag           `bson:"featureFlags" json:"featureFlags"`
	Permissions  permissions.Permissions `bson:"permissions" json:"permissions"`
}

func NewUser(id string) User {
	return User{
		ID:           id,
		Permissions:  permissions.User,
		FeatureFlags: []featureFlag{},
	}
}

func (u *User) HasFeatureFlag(flag featureFlag) bool {
	for _, f := range u.FeatureFlags {
		if f == flag {
			return true
		}
	}
	return false
}

type CompleteUser struct {
	User          `bson:",inline" json:",inline"`
	Subscriptions []UserSubscription `bson:"subscriptions" json:"subscriptions"`
	Connections   []UserConnection   `bson:"connections" json:"connections"`
}

func (u CompleteUser) Permissions() permissions.Permissions {
	perms := u.User.Permissions | permissions.User
	for _, c := range u.Connections {
		perms = perms.Add(c.Permissions)
	}
	for _, s := range u.Subscriptions {
		perms = perms.Add(s.Permissions)
	}
	return perms
}

func (u CompleteUser) Connection(connectionType ConnectionType) *UserConnection {
	for _, c := range u.Connections {
		if c.ConnectionType == connectionType {
			return &c
		}
	}
	return nil
}

func (u CompleteUser) Subscription(subscriptionType SubscriptionType) *UserSubscription {
	for _, s := range u.Subscriptions {
		if s.Type == subscriptionType {
			return &s
		}
	}
	return nil
}
