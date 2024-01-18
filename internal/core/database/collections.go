package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type collectionName string

const (
	CollectionUsers             = collectionName("users")
	CollectionUserConnections   = collectionName("user-connections")
	CollectionUserSubscriptions = collectionName("user-subscriptions")

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
func init() {
	// Users
	addIndexHandler(CollectionUsers, func(coll *mongo.Collection) ([]string, error) {
		return coll.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
			{
				Keys: bson.M{"featureFlags": 1},
			},
		})
	})
	addIndexHandler(CollectionUserConnections, func(coll *mongo.Collection) ([]string, error) {
		return coll.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
			{
				Keys: bson.M{"userID": 1},
			},
			{
				Keys: bson.D{
					{Key: "userID", Value: 1},
					{Key: "connectionType", Value: 1},
				},
			},
			{
				Keys: bson.D{
					{Key: "connectionID", Value: 1},
					{Key: "connectionType", Value: 1},
				},
			},
		})
	})
	addIndexHandler(CollectionUserSubscriptions, func(coll *mongo.Collection) ([]string, error) {
		return coll.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
			{
				Keys: bson.M{"userID": 1},
			},
			{
				Keys: bson.D{
					{Key: "userID", Value: 1},
					{Key: "creationDate", Value: -1},
				},
			},
			{
				Keys: bson.M{"referenceID": 1},
			},
			{
				Keys: bson.D{
					{Key: "referenceID", Value: 1},
					{Key: "creationDate", Value: -1},
				},
			},
		})
	})

	// Accounts, Clans, Sessions
	addIndexHandler(CollectionAccounts, func(coll *mongo.Collection) ([]string, error) {
		return coll.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
			{
				Keys: bson.M{"realm": 1},
			},
			{
				Keys: bson.M{"nickname": 1},
			},
		})
	})
	addIndexHandler(CollectionClans, func(coll *mongo.Collection) ([]string, error) {
		return coll.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
			{
				Keys: bson.M{"tag": 1},
			},
			{
				Keys: bson.M{"name": 1},
			},
			{
				Keys: bson.M{"members": 1},
			},
		})
	})
	addIndexHandler(CollectionSessions, func(coll *mongo.Collection) ([]string, error) {
		return coll.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
			{
				Keys: bson.D{
					{Key: "type", Value: 1},
					{Key: "accountId", Value: 1},
				},
			},
			{
				Keys: bson.D{
					{Key: "type", Value: 1},
					{Key: "accountId", Value: 1},
					{Key: "createdAt", Value: -1},
				},
			},
			{
				Keys: bson.D{
					{Key: "type", Value: 1},
					{Key: "accountId", Value: 1},
					{Key: "lastBattleTime", Value: -1},
				},
			},
			{
				Keys:    bson.M{"createdAt": 1},
				Options: options.Index().SetExpireAfterSeconds(604800),
			},
		})
	})

	// Glossary
	// this is all by ID for now, no need for index

	// Internal
	addIndexHandler(CollectionTasks, func(coll *mongo.Collection) ([]string, error) {
		return coll.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
			{
				Keys: bson.D{
					{Key: "kind", Value: 1},
					{Key: "createdAt", Value: -1},
				},
			},
			{
				Keys: bson.D{
					{Key: "status", Value: 1},
					{Key: "last_attempt", Value: 1},
				},
			},
			{
				Keys: bson.D{
					{Key: "status", Value: 1},
					{Key: "scheduled_after", Value: -1},
				},
			},
			{
				Keys:    bson.M{"createdAt": 1},
				Options: options.Index().SetExpireAfterSeconds(604800),
			},
		})
	})
	addIndexHandler(CollectionMessages, func(coll *mongo.Collection) ([]string, error) {
		return coll.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
			{
				Keys: bson.M{"userID": 1},
			},
			{
				Keys: bson.M{"guildID": 1},
			},
			{
				Keys: bson.D{
					{Key: "type", Value: 1},
					{Key: "guildID", Value: 1},
				},
			},
			{
				Keys: bson.M{"channelID": 1},
			},
			{
				Keys: bson.D{
					{Key: "type", Value: 1},
					{Key: "channelID", Value: 1},
				},
			},
			{

				Keys: bson.D{
					{Key: "guildID", Value: 1},
					{Key: "channelID", Value: 1},
				},
			},
			{
				Keys:    bson.M{"createdAt": 1},
				Options: options.Index().SetExpireAfterSeconds(604800),
			},
		})
	})
}
