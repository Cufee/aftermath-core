package users

import (
	"errors"

	"github.com/cufee/aftermath-core/internal/core/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrConnectionNotFound = errors.New("connection not found")
)

type connectionType string

const (
	ConnectionTypeWargaming = connectionType("wargaming")
)

type UserConnection struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"-"`

	UserID         string         `bson:"userID" json:"userID"`
	ExternalID     string         `bson:"connectionID" json:"connectionID"`
	ConnectionType connectionType `bson:"connectionType" json:"connectionType"`

	Metadata map[string]any `bson:"metadata" json:"metadata"`
}

func FindUserConnection(userId string, connectionType connectionType) (*UserConnection, error) {
	ctx, cancel := database.DefaultClient.Ctx()
	defer cancel()

	var connection UserConnection
	err := database.DefaultClient.Collection(database.CollectionUserConnections).FindOne(ctx, bson.M{"userID": userId, "connectionType": connectionType}).Decode(&connection)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrConnectionNotFound
		}
		return nil, err
	}

	return &connection, nil
}

func AddUserConnection(id string, connectionType connectionType, externalID string, metadata map[string]any) error {
	ctx, cancel := database.DefaultClient.Ctx()
	defer cancel()

	_, err := database.DefaultClient.Collection(database.CollectionUserConnections).InsertOne(ctx, UserConnection{
		UserID:         id,
		ConnectionType: connectionType,
		ExternalID:     externalID,
		Metadata:       metadata,
	})
	if err != nil {
		return err
	}

	return nil
}

func UpdateUserConnection(userId, externalId string, payload UserConnection, upsert bool) error {
	ctx, cancel := database.DefaultClient.Ctx()
	defer cancel()

	payload.UserID = userId
	payload.ExternalID = externalId

	opts := options.Update().SetUpsert(upsert)
	_, err := database.DefaultClient.Collection(database.CollectionUserConnections).UpdateOne(ctx, bson.M{"userID": userId, "connectionID": externalId}, bson.M{"$set": payload}, opts)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrConnectionNotFound
		}
		return err
	}

	return nil
}
