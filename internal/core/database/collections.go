package database

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type collectionName string

const (
	CollectionNonce = collectionName("nonce")

	CollectionUsers             = collectionName("users")
	CollectionUserContent       = collectionName("user-content")
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
	// Nonce
	addCollectionIndexes(CollectionNonce, []Index{
		{
			Keys:    bson.M{"referenceId": 1},
			Options: options.Index().SetName("referenceId"),
		},
		{
			Keys: bson.D{
				{Key: "referenceId", Value: 1},
				{Key: "expiresAt", Value: -1},
			},
			Options: options.Index().SetName("referenceId-expireAt"),
		},
		{
			Keys:    bson.M{"createdAt": 1},
			Options: options.Index().SetExpireAfterSeconds(604800).SetName("createdAt"),
		},
	})

	// Users
	addCollectionIndexes(CollectionUsers, []Index{
		{
			Keys:    bson.M{"featureFlags": 1},
			Options: options.Index().SetName("featureFlags"),
		},
	})
	addCollectionIndexes(CollectionUserContent, []Index{
		{
			Keys:    bson.M{"userID": 1},
			Options: options.Index().SetName("userID"),
		},
		{
			Keys:    bson.M{"referenceId": 1},
			Options: options.Index().SetName("referenceId"),
		},
		{
			Keys: bson.D{
				{Key: "type", Value: 1},
				{Key: "userID", Value: 1},
			},
			Options: options.Index().SetUnique(true).SetName("type-userID"),
		},
		{
			Keys: bson.D{
				{Key: "userID", Value: 1},
				{Key: "referenceId", Value: 1},
			},
			Options: options.Index().SetName("userID-referenceId"),
		},
		{
			Keys: bson.D{
				{Key: "type", Value: 1},
				{Key: "referenceId", Value: 1},
			},
			Options: options.Index().SetName("type-referenceId"),
		},
	})
	addCollectionIndexes(CollectionUserConnections, []Index{
		{
			Keys:    bson.M{"userID": 1},
			Options: options.Index().SetName("userID"),
		},
		{
			Keys: bson.D{
				{Key: "userID", Value: 1},
				{Key: "connectionType", Value: 1},
			},
			Options: options.Index().SetUnique(true).SetName("userID-connectionType"),
		},
		{
			Keys: bson.D{
				{Key: "connectionID", Value: 1},
				{Key: "connectionType", Value: 1},
			},
			Options: options.Index().SetName("connectionID-connectionType"),
		},
	})
	addCollectionIndexes(CollectionUserSubscriptions, []Index{
		{
			Keys:    bson.M{"userID": 1},
			Options: options.Index().SetName("userID"),
		},
		{
			Keys: bson.D{
				{Key: "userID", Value: 1},
				{Key: "creationDate", Value: -1},
			},
			Options: options.Index().SetName("userID-creationDate"),
		},
		{
			Keys:    bson.M{"referenceID": 1},
			Options: options.Index().SetName("referenceID"),
		},
		{
			Keys: bson.D{
				{Key: "referenceID", Value: 1},
				{Key: "creationDate", Value: -1},
			},
			Options: options.Index().SetName("referenceID-creationDate"),
		},
	})

	// Accounts, Clans, Sessions
	addCollectionIndexes(CollectionAccounts, []Index{
		{
			Keys:    bson.M{"realm": 1},
			Options: options.Index().SetName("realm"),
		},
		{
			Keys:    bson.M{"nickname": 1},
			Options: options.Index().SetName("nickname"),
		},
	})
	addCollectionIndexes(CollectionClans, []Index{
		{
			Keys:    bson.M{"tag": 1},
			Options: options.Index().SetName("tag"),
		},
		{
			Keys:    bson.M{"name": 1},
			Options: options.Index().SetName("name"),
		},
		{
			Keys:    bson.M{"members": 1},
			Options: options.Index().SetName("members"),
		},
	})
	addCollectionIndexes(CollectionSessions, []Index{
		{
			Keys: bson.D{
				{Key: "type", Value: 1},
				{Key: "accountId", Value: 1},
			},
			Options: options.Index().SetName("type-accountId"),
		},
		{
			Keys: bson.D{
				{Key: "type", Value: 1},
				{Key: "accountId", Value: 1},
				{Key: "createdAt", Value: -1},
			},
			Options: options.Index().SetName("type-accountId-createdAt"),
		},
		{
			Keys: bson.D{
				{Key: "type", Value: 1},
				{Key: "accountId", Value: 1},
				{Key: "lastBattleTime", Value: -1},
			},
			Options: options.Index().SetName("type-accountId-lastBattleTime"),
		},
		{
			Keys:    bson.M{"createdAt": 1},
			Options: options.Index().SetExpireAfterSeconds(2_592_864).SetName("createdAt"),
		},
	})

	// Glossary
	// this is all by ID for now, no need for index

	// Internal
	addCollectionIndexes(CollectionTasks, []Index{
		{
			Keys: bson.D{
				{Key: "kind", Value: 1},
				{Key: "createdAt", Value: -1},
			},
			Options: options.Index().SetName("kind-createdAt"),
		},
		{
			Keys: bson.D{
				{Key: "status", Value: 1},
				{Key: "last_attempt", Value: 1},
			},
			Options: options.Index().SetName("status-last_attempt"),
		},
		{
			Keys: bson.D{
				{Key: "status", Value: 1},
				{Key: "scheduled_after", Value: -1},
			},
			Options: options.Index().SetName("status-scheduled_after"),
		},
		{
			Keys:    bson.M{"createdAt": 1},
			Options: options.Index().SetExpireAfterSeconds(604800).SetName("createdAt"),
		},
	})
	addCollectionIndexes(CollectionMessages, []Index{
		{
			Keys:    bson.M{"userID": 1},
			Options: options.Index().SetName("userID"),
		},
		{
			Keys:    bson.M{"guildID": 1},
			Options: options.Index().SetName("guildID"),
		},
		{
			Keys: bson.D{
				{Key: "type", Value: 1},
				{Key: "guildID", Value: 1},
			},
			Options: options.Index().SetName("type-guildID"),
		},
		{
			Keys:    bson.M{"channelID": 1},
			Options: options.Index().SetName("channelID"),
		},
		{
			Keys: bson.D{
				{Key: "type", Value: 1},
				{Key: "channelID", Value: 1},
			},
			Options: options.Index().SetName("type-channelID"),
		},
		{
			Keys: bson.D{
				{Key: "guildID", Value: 1},
				{Key: "channelID", Value: 1},
			},
			Options: options.Index().SetName("guildID-channelID"),
		},
		{
			Keys:    bson.M{"createdAt": 1},
			Options: options.Index().SetExpireAfterSeconds(604800).SetName("createdAt"),
		},
	})
}
