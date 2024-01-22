package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Nonce struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	CreatedAt time.Time          `bson:"createdAt"`
	ExpiresAt time.Time          `bson:"expiresAt"`

	ReferenceID string         `bson:"referenceId"`
	Metadata    map[string]any `bson:"metadata"`
}
