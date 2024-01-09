package cache

import (
	"errors"
	"time"

	"github.com/cufee/aftermath-core/internal/core/database"
	core "github.com/cufee/aftermath-core/internal/core/stats"
	"github.com/cufee/aftermath-core/internal/logic/sessions"
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

func RefreshSessions(sessionType SessionType, realm string, accountIDs ...int) error {
	sessions, err := sessions.GetSessionsWithClient(cacheClient, realm, accountIDs...)
	if err != nil {
		return err
	}

	var writes []mongo.WriteModel
	for _, session := range sessions {
		model := mongo.NewInsertOneModel()
		model.SetDocument(SessionDatabaseRecord{
			Type:      sessionType,
			CreatedAt: time.Now(),
			Session:   session.Session,
		})
		writes = append(writes, model)
	}

	ctx, cancel := database.DefaultClient.Ctx()
	defer cancel()

	_, err = database.DefaultClient.Collection(database.CollectionSessions).BulkWrite(ctx, writes)
	if err != nil {
		return err
	}

	return nil
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
		query["createdAt"] = bson.M{"$lt": opts.Before}
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
