package database

import "go.mongodb.org/mongo-driver/mongo"

type collectionName string

const (
	CollectionSessions = collectionName("sessions")
	CollectionAccounts = collectionName("accounts")
	CollectionClans    = collectionName("clans")
)

/*
Any Indexes required for the application to work should be created/updated here.
This function is called on startup.
*/
func syncIndexes(db *mongo.Database) error {
	return nil
}
