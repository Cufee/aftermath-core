package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	UserContentTypeBackground = UserContentType("background-image")
)

type UserContentType string

func (t UserContentType) Valid() bool {
	switch t {
	case UserContentTypeBackground:
		return true
	default:
		return false
	}
}

type UserContent[T any] struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	UserID string             `bson:"userId"`

	Type      UserContentType `bson:"type"`
	CreatedAt time.Time       `bson:"createdAt"`

	Data     T              `bson:"data"`
	Metadata map[string]any `bson:"metadata"`
}
