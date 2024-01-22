package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	UserContentTypeClanBackground     = UserContentType("clan-background-image")
	UserContentTypePersonalBackground = UserContentType("personal-background-image")
)

type UserContentType string

func (t UserContentType) Valid() bool {
	switch t {
	case UserContentTypeClanBackground:
		return true
	case UserContentTypePersonalBackground:
		return true
	default:
		return false
	}
}

type UserContent[T any] struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	UserID      string             `bson:"userID"`
	ReferenceID string             `bson:"referenceId"`

	Type      UserContentType `bson:"type"`
	UpdatedAt time.Time       `bson:"createdAt"`

	Data     T              `bson:"data"`
	Metadata map[string]any `bson:"metadata"`
}
