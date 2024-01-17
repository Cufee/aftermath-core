package cache

import (
	"errors"
	"time"

	"github.com/cufee/aftermath-core/internal/core/database"
	core "github.com/cufee/aftermath-core/internal/core/stats"
	"github.com/cufee/aftermath-core/internal/core/wargaming"
	"github.com/cufee/aftermath-core/internal/logic/sessions"
	wg "github.com/cufee/am-wg-proxy-next/types"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SessionType string

const (
	SessionTypeDaily = SessionType("daily")
)

var (
	ErrNoSessionCache = errors.New("no session found")
)

type SessionDatabaseRecord struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Type      SessionType        `bson:"type"`
	CreatedAt time.Time          `bson:"createdAt"`

	Session *core.SessionSnapshot `bson:",inline"`
}

func RefreshSessionsAndAccounts(sessionType SessionType, realm string, accountIDs ...int) (map[int]error, error) {
	sessions, err := sessions.GetSessionsWithClient(wargaming.Clients.Cache, realm, accountIDs...)
	if err != nil {
		return nil, err
	}

	var accounts []*wg.ExtendedAccount
	for _, session := range sessions {
		accounts = append(accounts, session.Data.Account)
	}
	err = UpdatePlayerAccountsFromWG(realm, accounts...)
	if err != nil {
		return nil, err
	}

	lastBattle, err := GetLastBattleTimes(sessionType, accountIDs...)
	if err != nil {
		return nil, err
	}

	updateErrors := make(map[int]error)
	var sessionInserts []mongo.WriteModel
	for accountId, session := range sessions {
		if session.Err != nil {
			updateErrors[accountId] = session.Err
			continue
		}

		if lastBattle[session.Data.Account.ID] >= session.Data.Session.LastBattleTime {
			log.Debug().Msgf("%d played 0 battles since last session, skipping update", session.Data.Account.ID)
			continue
		}
		model := mongo.NewInsertOneModel()
		model.SetDocument(SessionDatabaseRecord{
			Type:      sessionType,
			CreatedAt: time.Now(),
			Session:   session.Data.Session,
		})
		sessionInserts = append(sessionInserts, model)
	}
	if len(sessionInserts) == 0 {
		return updateErrors, nil
	}

	ctx, cancel := database.DefaultClient.Ctx()
	defer cancel()

	_, err = database.DefaultClient.Collection(database.CollectionSessions).BulkWrite(ctx, sessionInserts)
	if err != nil {
		return nil, err
	}

	return updateErrors, nil
}

type SessionGetOptions struct {
	LastBattleAfter *int
	Before          time.Time
	Type            SessionType
}

func GetPlayerSessionSnapshot(accountID int, o ...SessionGetOptions) (*core.SessionSnapshot, error) {
	opts := SessionGetOptions{Type: SessionTypeDaily}
	if len(o) > 0 {
		opts = o[0]
	}

	ctx, cancel := database.DefaultClient.Ctx()
	defer cancel()

	findOptions := options.FindOne()
	findOptions.SetSort(bson.M{"createdAt": -1})

	query := bson.M{"accountId": accountID}
	if !opts.Before.IsZero() {
		query["createdAt"] = bson.M{"$gt": opts.Before}
		findOptions.SetSort(bson.M{"createdAt": 1})
	}
	if opts.Type != "" {
		query["type"] = opts.Type
	}
	if opts.LastBattleAfter != nil {
		query["lastBattleTime"] = bson.M{"$gt": *opts.LastBattleAfter}
	}

	var session SessionDatabaseRecord
	err := database.DefaultClient.Collection(database.CollectionSessions).FindOne(ctx, query, findOptions).Decode(&session)
	if err != nil {
		if errors.Is(mongo.ErrNoDocuments, err) {
			return nil, ErrNoSessionCache
		}
		return nil, err
	}

	return session.Session, nil
}

func GetLastBattleTimes(sessionType SessionType, accountIDs ...int) (map[int]int, error) {
	if len(accountIDs) == 0 {
		return make(map[int]int), nil
	}

	ctx, cancel := database.DefaultClient.Ctx()
	defer cancel()

	var lastBattles map[int]int = make(map[int]int)
	var pipeline mongo.Pipeline
	pipeline = append(pipeline, bson.D{{Key: "$match", Value: bson.M{"accountId": bson.M{"$in": accountIDs}, "type": sessionType}}})
	pipeline = append(pipeline, bson.D{{Key: "$group", Value: bson.M{"_id": "$accountId", "lastBattleTime": bson.M{"$max": "$lastBattleTime"}}}})
	pipeline = append(pipeline, bson.D{{Key: "$project", Value: bson.M{"_id": 0, "accountId": "$_id", "lastBattleTime": 1}}})
	cur, err := database.DefaultClient.Collection(database.CollectionSessions).Aggregate(ctx, pipeline)
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
