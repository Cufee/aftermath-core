package database

import (
	"errors"

	"github.com/cufee/aftermath-core/internal/core/database/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrConnectionNotFound = errors.New("connection not found")
)

func FindUserConnection(userId string, connectionType models.ConnectionType) (*models.UserConnection, error) {
	ctx, cancel := DefaultClient.Ctx()
	defer cancel()

	var connection models.UserConnection
	err := DefaultClient.Collection(CollectionUserConnections).FindOne(ctx, bson.M{"userID": userId, "connectionType": connectionType}).Decode(&connection)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrConnectionNotFound
		}
		return nil, err
	}

	return &connection, nil
}

func GetUserConnections(userId string) ([]models.UserConnection, error) {
	ctx, cancel := DefaultClient.Ctx()
	defer cancel()

	var connections []models.UserConnection
	cur, err := DefaultClient.Collection(CollectionUserConnections).Find(ctx, bson.M{"userID": userId})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrConnectionNotFound
		}
		return nil, err
	}

	return connections, cur.All(ctx, &connections)
}

func AddUserConnection(id string, connectionType models.ConnectionType, externalID string, metadata map[string]any) error {
	ctx, cancel := DefaultClient.Ctx()
	defer cancel()

	_, err := DefaultClient.Collection(CollectionUserConnections).InsertOne(ctx, models.UserConnection{
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

func UpdateUserConnection(userId string, connectionType models.ConnectionType, payload models.UserConnection, upsert bool) error {
	ctx, cancel := DefaultClient.Ctx()
	defer cancel()

	payload.UserID = userId
	opts := options.Update().SetUpsert(upsert)
	_, err := DefaultClient.Collection(CollectionUserConnections).UpdateOne(ctx, bson.M{"userID": userId, "connectionType": connectionType}, bson.M{"$set": payload}, opts)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrConnectionNotFound
		}
		return err
	}

	return nil
}
