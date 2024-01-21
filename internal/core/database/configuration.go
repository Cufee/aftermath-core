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
	ErrConfigurationNotFound = errors.New("configuration not found")
)

func UpdateAppConfiguration[T any](key string, data T, metadata map[string]any, upsert bool) error {
	var payload models.AppConfiguration[T]
	payload.UpdatedAt = time.Now()
	payload.Metadata = metadata
	payload.Value = data
	payload.Key = key

	ctx, cancel := DefaultClient.Ctx()
	defer cancel()

	opts := options.Update().SetUpsert(upsert)

	_, err := DefaultClient.Collection(CollectionConfiguration).UpdateOne(ctx, bson.M{"key": key}, bson.M{"$set": payload}, opts)
	if err != nil {
		return err
	}
	return nil
}

func GetAppConfiguration[T any](key string) (*models.AppConfiguration[T], error) {
	ctx, cancel := DefaultClient.Ctx()
	defer cancel()

	var content models.AppConfiguration[T]
	err := DefaultClient.Collection(CollectionConfiguration).FindOne(ctx, bson.M{"key": key}).Decode(&content)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrConfigurationNotFound
		}
		return nil, err
	}
	return &content, nil
}
