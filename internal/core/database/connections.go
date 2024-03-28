package database

import (
	"errors"

	"github.com/cufee/aftermath-core/internal/core/database/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrConnectionNotFound = errors.New("connection not found")
)

func FindUserConnection(userId string, connectionType models.ConnectionType) (models.UserConnection, error) {
	ctx, cancel := DefaultClient.Ctx()
	defer cancel()

	var connection models.UserConnection
	err := DefaultClient.Collection(CollectionUserConnections).FindOne(ctx, bson.M{"userID": userId, "connectionType": connectionType}).Decode(&connection)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return connection, ErrConnectionNotFound
		}
		return connection, err
	}

	return connection, nil
}

func FindConnectionsByReferenceID(referenceId string, connectionType models.ConnectionType) ([]models.UserConnection, error) {
	ctx, cancel := DefaultClient.Ctx()
	defer cancel()

	var connections []models.UserConnection
	cur, err := DefaultClient.Collection(CollectionUserConnections).Find(ctx, bson.M{"connectionID": referenceId, "connectionType": connectionType})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrConnectionNotFound
		}
		return nil, err
	}
	return connections, cur.All(ctx, &connections)
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

func GetUserConnection(userId string, connectionType models.ConnectionType) (models.UserConnection, error) {
	ctx, cancel := DefaultClient.Ctx()
	defer cancel()

	var connection models.UserConnection
	err := DefaultClient.Collection(CollectionUserConnections).FindOne(ctx, bson.M{"userID": userId, "connectionType": connectionType}).Decode(&connection)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return connection, ErrConnectionNotFound
		}
		return connection, err
	}

	return connection, nil
}

func AddUserConnection(userId string, connectionType models.ConnectionType, externalID string, metadata map[string]any) (models.UserConnection, error) {
	ctx, cancel := DefaultClient.Ctx()
	defer cancel()

	connection := models.UserConnection{
		UserID:         userId,
		ConnectionType: connectionType,
		ExternalID:     externalID,
		Metadata:       metadata,
	}

	res, err := DefaultClient.Collection(CollectionUserConnections).InsertOne(ctx, connection)
	if err != nil {
		return models.UserConnection{}, err
	}
	insertedId, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return models.UserConnection{}, errors.New("failed to cast inserted ID to primitive.ObjectID")
	}

	connection.ID = insertedId
	return connection, nil
}

func UpdateUserConnection(userId string, connectionType models.ConnectionType, payload models.ConnectionUpdate) (models.UserConnection, error) {
	ctx, cancel := DefaultClient.Ctx()
	defer cancel()

	_, err := DefaultClient.Collection(CollectionUserConnections).UpdateOne(ctx, bson.M{"userID": userId, "connectionType": connectionType}, bson.M{"$set": payload})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.UserConnection{}, ErrConnectionNotFound
		}
		return models.UserConnection{}, err
	}

	return GetUserConnection(userId, connectionType)
}

func UpdateManyConnectionsByReferenceID(referenceId string, connectionType models.ConnectionType, payload models.ConnectionUpdate) error {
	ctx, cancel := DefaultClient.Ctx()
	defer cancel()

	_, err := DefaultClient.Collection(CollectionUserConnections).UpdateMany(ctx, bson.M{"connectionID": referenceId, "connectionType": connectionType}, bson.M{"$set": payload})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrConnectionNotFound
		}
		return err
	}

	return nil
}
