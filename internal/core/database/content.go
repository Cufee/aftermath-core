package database

import (
	"errors"
	"time"

	"github.com/cufee/aftermath-core/internal/core/database/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrUserContentNotFound = errors.New("content not found")
)

func UpdateUserContent[T any](userID, referenceID string, contentType models.UserContentType, data T, metadata map[string]any, upsert bool) error {
	var payload models.UserContent[T]
	payload.ReferenceID = referenceID
	payload.UpdatedAt = time.Now()
	payload.Metadata = metadata
	payload.UserID = userID
	payload.Type = contentType
	payload.Data = data

	ctx, cancel := DefaultClient.Ctx()
	defer cancel()

	opts := options.Update().SetUpsert(upsert)

	_, err := DefaultClient.Collection(CollectionUserContent).UpdateOne(ctx, bson.M{"userId": userID, "type": contentType}, bson.M{"$set": payload}, opts)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrUserContentNotFound
		}
		return err
	}
	return nil
}

func UpdateUserContentReferenceID[T any](userID string, contentType models.UserContentType, newReferenceID string) (*models.UserContent[T], error) {
	ctx, cancel := DefaultClient.Ctx()
	defer cancel()

	var content models.UserContent[T]
	err := DefaultClient.Collection(CollectionUserContent).FindOneAndUpdate(ctx, bson.M{"userId": userID, "type": contentType}, bson.M{"$set": bson.M{"referenceId": newReferenceID}}).Decode(&content)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrUserContentNotFound
		}
		return nil, err
	}
	return &content, nil
}

func GetUserContent[T any](userID string, contentType ...models.UserContentType) (*models.UserContent[T], error) {
	ctx, cancel := DefaultClient.Ctx()
	defer cancel()

	var content models.UserContent[T]
	err := DefaultClient.Collection(CollectionUserContent).FindOne(ctx, bson.M{"userId": userID, "type": contentType}).Decode(&content)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrUserContentNotFound
		}
		return nil, err
	}

	return &content, nil
}

func GetContentByReferenceIDs[T any](referenceIDs []string, contentType ...models.UserContentType) ([]models.UserContent[T], error) {
	if len(referenceIDs) == 0 {
		return nil, ErrUserContentNotFound
	}

	ctx, cancel := DefaultClient.Ctx()
	defer cancel()

	var content []models.UserContent[T]
	cur, err := DefaultClient.Collection(CollectionUserContent).Find(ctx, bson.M{"referenceId": bson.M{"$in": referenceIDs}, "type": contentType})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrUserContentNotFound
		}
		return nil, err
	}

	return content, cur.All(ctx, &content)
}
