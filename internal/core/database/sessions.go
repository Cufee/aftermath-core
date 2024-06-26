package database

import (
	"errors"
	"time"

	"github.com/cufee/aftermath-core/internal/core/database/models"
	"github.com/cufee/aftermath-core/internal/core/stats"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrNoSessionCache = errors.New("no session found")
)

type SessionGetOptions struct {
	LastBattleBefore *int
	LastBattleAfter  *int
	ReferenceID      *string
	Type             models.SessionType
}

func GetPlayerSessionSnapshot(accountID int, o ...SessionGetOptions) (models.Snapshot, error) {
	opts := SessionGetOptions{Type: models.SessionTypeDaily}
	if len(o) > 0 {
		opts = o[0]
	}

	ctx, cancel := DefaultClient.Ctx()
	defer cancel()

	findOptions := options.FindOne()
	findOptions.SetSort(bson.M{"createdAt": -1})

	query := bson.M{"accountId": accountID}
	if opts.Type != "" {
		query["type"] = opts.Type
	}
	if opts.ReferenceID != nil {
		query["referenceId"] = opts.ReferenceID
	}
	if opts.LastBattleBefore != nil {
		query["lastBattleTime"] = bson.M{"$lt": *opts.LastBattleBefore}
	}
	if opts.LastBattleAfter != nil {
		query["lastBattleTime"] = bson.M{"$gt": *opts.LastBattleAfter}
	}

	var snapshot models.Snapshot
	err := DefaultClient.Collection(CollectionSessions).FindOne(ctx, query, findOptions).Decode(&snapshot)
	if err != nil {
		if errors.Is(mongo.ErrNoDocuments, err) {
			return snapshot, ErrNoSessionCache
		}
		return snapshot, err
	}

	return snapshot, nil
}

func GetLastBattleTimes(sessionType models.SessionType, referenceId *string, accountIDs ...int) (map[int]int, error) {
	if len(accountIDs) == 0 {
		return make(map[int]int), nil
	}

	ctx, cancel := DefaultClient.Ctx()
	defer cancel()

	filter := bson.M{"accountId": bson.M{"$in": accountIDs}, "type": sessionType}
	if referenceId != nil {
		filter["referenceId"] = *referenceId
	}
	var lastBattles map[int]int = make(map[int]int)
	var pipeline mongo.Pipeline
	pipeline = append(pipeline, bson.D{{Key: "$match", Value: filter}})
	pipeline = append(pipeline, bson.D{{Key: "$group", Value: bson.M{"_id": "$accountId", "lastBattleTime": bson.M{"$max": "$lastBattleTime"}}}})
	pipeline = append(pipeline, bson.D{{Key: "$project", Value: bson.M{"_id": 0, "accountId": "$_id", "lastBattleTime": 1}}})
	cur, err := DefaultClient.Collection(CollectionSessions).Aggregate(ctx, pipeline)
	if err != nil {
		if errors.Is(mongo.ErrNoDocuments, err) {
			return lastBattles, nil
		}
		return nil, err
	}

	var results []struct {
		AccountID      int `bson:"accountId"`
		LastBattleTime int `bson:"lastBattleTime"`
	}
	err = cur.All(ctx, &results)
	if err != nil {
		return nil, err
	}

	for _, result := range results {
		lastBattles[result.AccountID] = result.LastBattleTime
	}
	return lastBattles, nil
}

func InsertSession(sessionType models.SessionType, referenceId *string, sessions ...stats.SessionSnapshot) error {
	var sessionInserts []mongo.WriteModel
	for _, session := range sessions {
		model := mongo.NewInsertOneModel()
		snapshot := models.Snapshot{
			Type:      sessionType,
			CreatedAt: time.Now(),
			Session:   session,
		}
		if referenceId != nil {
			snapshot.ReferenceID = *referenceId
		}
		model.SetDocument(snapshot)
		sessionInserts = append(sessionInserts, model)
	}
	if len(sessionInserts) == 0 {
		return nil
	}

	ctx, cancel := DefaultClient.Ctx()
	defer cancel()

	_, err := DefaultClient.Collection(CollectionSessions).BulkWrite(ctx, sessionInserts)
	if err != nil {
		return err
	}

	return nil
}
