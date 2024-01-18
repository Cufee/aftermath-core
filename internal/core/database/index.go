package database

import (
	"context"

	"slices"
	"strings"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
)

var indexHandlers = make(map[collectionName](func(collection *mongo.Database) ([]string, error)))

func addIndexHandler(collection collectionName, handler func(db *mongo.Collection) ([]string, error)) {
	indexHandlers[collection] = (func(db *mongo.Database) ([]string, error) {
		return handler(db.Collection(string(collection)))
	})

}

func SyncIndexes(db *mongo.Database) error {
	log.Info().Msg("Syncing indexes")
	defer log.Info().Msg("Finished syncing indexes")

	toDelete := make(map[collectionName][]string)
	for collection, handler := range indexHandlers {
		names, err := handler(db)
		if err != nil {
			if strings.Contains(err.Error(), "Index already exists") {
				continue
			}
			return err
		}
		defer log.Debug().Msgf("Synced indexes for %s: %v", collection, names)

		current, err := db.Collection(string(collection)).Indexes().ListSpecifications(context.Background())
		if err != nil {
			return err
		}

		for _, index := range current {
			if index.Name == "_id_" {
				continue
			}

			if !slices.Contains(names, index.Name) {
				toDelete[collection] = append(toDelete[collection], index.Name)
			}
		}
	}

	for collection, names := range toDelete {
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
