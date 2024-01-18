package models

import (
	"time"

	"github.com/cufee/aftermath-core/internal/core/stats"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SessionType string

const (
	SessionTypeDaily = SessionType("daily")
)

type Snapshot struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Type      SessionType        `bson:"type"`
	CreatedAt time.Time          `bson:"createdAt"`

	Session *stats.SessionSnapshot `bson:",inline"`
}
