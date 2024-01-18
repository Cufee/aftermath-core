package models

type featureFlag string

const (
	FeatureFlagCustomizationDisabled = featureFlag("customizationDisabled")
)

type User struct {
	ID string `bson:"_id" json:"id"`

	FeatureFlags []featureFlag `bson:"featureFlags" json:"featureFlags"`
}

func NewUser(id string) User {
	return User{
		ID:           id,
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
