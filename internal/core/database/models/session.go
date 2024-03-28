package models

import (
	"time"

	"github.com/cufee/aftermath-core/internal/core/stats"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SessionType string

const (
	SessionTypeDaily = SessionType("daily")
	SessionTypeLive  = SessionType("live")
)

func ParseSessionType(input string) SessionType {
	switch input {
	case "live":
		return SessionTypeLive
	default:
		return SessionTypeDaily
	}
}

type Snapshot struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Type        SessionType        `bson:"type"`
	CreatedAt   time.Time          `bson:"createdAt"`
	ReferenceID string             `bson:"referenceId"`

	Session stats.SessionSnapshot `bson:",inline"`
}
