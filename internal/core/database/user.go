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

func GetOrCreateUserByID(id string) (models.CompleteUser, error) {
	user, err := GetUserByID(id)
	if err != nil {
		if !errors.Is(err, ErrUserNotFound) {
			return models.CompleteUser{}, err
		}
		partial, err := CreateUser(id)
		if err != nil {
			return models.CompleteUser{}, err
		}
		user = models.CompleteUser{User: partial}
	}
	return user, nil
}

func GetUserByID(id string) (models.CompleteUser, error) {
	ctx, cancel := DefaultClient.Ctx()
	defer cancel()

	var pipeline mongo.Pipeline
	pipeline = append(pipeline, bson.D{{Key: "$match", Value: bson.D{{Key: "_id", Value: id}}}})
	pipeline = append(pipeline, bson.D{{Key: "$limit", Value: 1}})
	pipeline = append(pipeline, bson.D{{Key: "$lookup", Value: bson.D{{Key: "from", Value: CollectionUserConnections}, {Key: "localField", Value: "_id"}, {Key: "foreignField", Value: "userID"}, {Key: "as", Value: "connections"}}}})
	pipeline = append(pipeline, bson.D{{Key: "$lookup", Value: bson.D{{Key: "from", Value: CollectionUserSubscriptions}, {Key: "localField", Value: "_id"}, {Key: "foreignField", Value: "userID"}, {Key: "as", Value: "subscriptions"}}}})

	cur, err := DefaultClient.Collection(CollectionUsers).Aggregate(ctx, pipeline)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.CompleteUser{}, ErrUserNotFound
		}
		return models.CompleteUser{}, err
	}

	var results []models.CompleteUser
	if err := cur.All(ctx, &results); err != nil {
		return models.CompleteUser{}, err
	}

	if len(results) == 0 {
		return models.CompleteUser{}, ErrUserNotFound
	}
	if len(results) > 1 {
		return models.CompleteUser{}, errors.New("multiple users found")
	}

	return results[0], nil
}

func FindUserByConnection(connectionType models.ConnectionType, externalID string) (models.CompleteUser, error) {
	ctx, cancel := DefaultClient.Ctx()
	defer cancel()

	var pipeline mongo.Pipeline
	pipeline = append(pipeline, bson.D{{Key: "$match", Value: bson.D{{Key: "connectionType", Value: connectionType}, {Key: "connectionID", Value: externalID}}}})
	pipeline = append(pipeline, bson.D{{Key: "$limit", Value: 1}})
	pipeline = append(pipeline, bson.D{{Key: "$lookup", Value: bson.D{{Key: "from", Value: CollectionUsers}, {Key: "localField", Value: "userID"}, {Key: "foreignField", Value: "_id"}, {Key: "as", Value: "user"}}}})
	pipeline = append(pipeline, bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$user"}}}})
	pipeline = append(pipeline, bson.D{{Key: "$replaceRoot", Value: bson.D{{Key: "newRoot", Value: "$user"}}}})
	pipeline = append(pipeline, bson.D{{Key: "$lookup", Value: bson.D{{Key: "from", Value: CollectionUserConnections}, {Key: "localField", Value: "_id"}, {Key: "foreignField", Value: "userID"}, {Key: "as", Value: "connections"}}}})
	pipeline = append(pipeline, bson.D{{Key: "$lookup", Value: bson.D{{Key: "from", Value: CollectionUserSubscriptions}, {Key: "localField", Value: "_id"}, {Key: "foreignField", Value: "userID"}, {Key: "as", Value: "subscriptions"}}}})

	cur, err := DefaultClient.Collection(CollectionUsers).Aggregate(ctx, pipeline)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.CompleteUser{}, ErrUserNotFound
		}
		return models.CompleteUser{}, err
	}

	var users []models.CompleteUser
	if err := cur.All(ctx, &users); err != nil {
		return models.CompleteUser{}, err
	}

	if len(users) == 0 {
		return models.CompleteUser{}, ErrUserNotFound
	}
	if len(users) > 1 {
		return models.CompleteUser{}, errors.New("multiple users found")
	}

	return users[0], nil
}

func CreateUser(id string) (models.User, error) {
	ctx, cancel := DefaultClient.Ctx()
	defer cancel()

	user := models.NewUser(id)

	_, err := DefaultClient.Collection(CollectionUsers).InsertOne(ctx, user)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func UpdateUser(id string, update models.User) (models.User, error) {
	update.ID = id

	ctx, cancel := DefaultClient.Ctx()
	defer cancel()

	_, err := DefaultClient.Collection(CollectionUsers).UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": update})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.User{}, ErrUserNotFound
		}
		return models.User{}, err
	}

	return update, nil
}
