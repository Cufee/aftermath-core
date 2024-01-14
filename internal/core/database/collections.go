package database

import "go.mongodb.org/mongo-driver/mongo"

type collectionName string

const (
	CollectionUsers           = collectionName("users")
	CollectionUserConnections = collectionName("user-connections")

	CollectionClans    = collectionName("clans")
	CollectionAccounts = collectionName("accounts")
	CollectionSessions = collectionName("sessions")

	CollectionConfiguration = collectionName("configuration")
)

/*
Any Indexes required for the application to work should be created/updated here.
This function is called on startup.
*/
func syncIndexes(db *mongo.Database) error {
	return nil
}
