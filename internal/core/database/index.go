package database

import (
	"context"
	"errors"

	"slices"
	"strings"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Index mongo.IndexModel

func (i *Index) Name() (string, error) {
	if i.Options == nil || i.Options.Name == nil {
		var keys []string

		switch i.Keys.(type) {
		case bson.M:
			for key := range i.Keys.(bson.M) {
				keys = append(keys, key)
			}
		case bson.D:
			for _, key := range i.Keys.(bson.D) {
				keys = append(keys, key.Key)
			}
		default:
			return "", errors.New("invalid index keys")
		}
		return strings.Join(keys, "-"), nil
	}

	return *i.Options.Name, nil
}

var collectionIndexes = make(map[string][]Index)

func addCollectionIndexes(collection collectionName, indexes []Index) {
	collectionIndexes[string(collection)] = indexes
}

func SyncIndexes(db *mongo.Database) error {
	log.Info().Msg("Syncing indexes")
	defer log.Info().Msg("Finished syncing indexes")

	indexesToDelete := make(map[string][]string)

	for collection, indexes := range collectionIndexes {
		currentIndexes, err := db.Collection(string(collection)).Indexes().ListSpecifications(context.Background())
		if err != nil {
			return err
		}
		var currentIndexNames []string
		for _, index := range currentIndexes {
			currentIndexNames = append(currentIndexNames, index.Name)
		}

		var desiredCollectionIndexNames []string = []string{"_id_"}
		var indexesToCreate []mongo.IndexModel

		for _, index := range indexes {
			name, err := index.Name()
			if err != nil {
				return err
			}
			desiredCollectionIndexNames = append(desiredCollectionIndexNames, name)

			if !slices.Contains(currentIndexNames, name) {
				indexesToCreate = append(indexesToCreate, mongo.IndexModel(index))
			}
		}

		for _, index := range currentIndexes {
			if !slices.Contains(desiredCollectionIndexNames, index.Name) {
				indexesToDelete[collection] = append(indexesToDelete[collection], index.Name)
			}
		}

		if len(indexesToCreate) == 0 {
			log.Debug().Msgf("No new indexes to create for %s", collection)
			continue
		}

		createdNames, err := db.Collection(collection).Indexes().CreateMany(context.Background(), indexesToCreate)
		if err != nil {
			return err
		}
		log.Debug().Msgf("Created indexes for %s: %v", collection, createdNames)
	}

	for collection, names := range indexesToDelete {
		log.Debug().Msgf("Deleting indexes for %s: %v", collection, names)
		for _, name := range names {
			_, err := db.Collection(string(collection)).Indexes().DropOne(context.Background(), name)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
