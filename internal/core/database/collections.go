package database

import "go.mongodb.org/mongo-driver/mongo"

type collectionName string

const (
	CollectionUsers           = collectionName("users")
	CollectionUserConnections = collectionName("user-connections")

	CollectionClans    = collectionName("clans")
	CollectionAccounts = collectionName("accounts")
	CollectionSessions = collectionName("sessions")

	CollectionVehicleAverages     = collectionName("vehicle-averages")
	CollectionVehicleGlossary     = collectionName("glossary-vehicles")
	CollectionAchievementGlossary = collectionName("glossary-achievements")

	CollectionTasks         = collectionName("tasks")
	CollectionMessages      = collectionName("messages")
	CollectionConfiguration = collectionName("configuration")
)

/*
Any Indexes required for the application to work should be created/updated here.
This function is called on startup.
*/
func syncIndexes(db *mongo.Database) error {
	return nil
}
