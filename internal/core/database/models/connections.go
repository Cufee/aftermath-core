package models

import (
	"github.com/cufee/aftermath-core/permissions/v1"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ConnectionType string

const (
	ConnectionTypeWargaming = ConnectionType("wargaming")
)

type UserConnection struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`

	UserID         string                  `bson:"userID" json:"userID"`
	ExternalID     string                  `bson:"connectionID" json:"connectionID"`
	ConnectionType ConnectionType          `bson:"connectionType" json:"connectionType"`
	Permissions    permissions.Permissions `bson:"permissions" json:"permissions"`

	Metadata map[string]any `bson:"metadata" json:"metadata"`
}
