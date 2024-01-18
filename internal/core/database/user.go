package database

import (
	"errors"

	"github.com/cufee/aftermath-core/internal/core/database/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

func FindUserByID(id string) (*models.User, error) {
	ctx, cancel := DefaultClient.Ctx()
	defer cancel()

	var user models.User
	err := DefaultClient.Collection(CollectionUsers).FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func FindUserByConnection(connectionType models.ConnectionType, externalID string) (*models.User, error) {
	ctx, cancel := DefaultClient.Ctx()
	defer cancel()

	var connection models.UserConnection
	err := DefaultClient.Collection(CollectionUsers).FindOne(ctx, bson.M{"connectionType": connectionType, "connectionID": externalID}).Decode(&connection)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return FindUserByID(connection.UserID)
}

func CreateUser(id string) (*models.User, error) {
	ctx, cancel := DefaultClient.Ctx()
	defer cancel()

	user := models.NewUser(id)

	_, err := DefaultClient.Collection(CollectionUsers).InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func UpdateUser(id string, update models.User) (*models.User, error) {
	update.ID = id

	ctx, cancel := DefaultClient.Ctx()
	defer cancel()

	_, err := DefaultClient.Collection(CollectionUsers).UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": update})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &update, nil
}
