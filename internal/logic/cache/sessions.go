package cache

import (
	"errors"
	"time"

	"github.com/cufee/aftermath-core/internal/core/database"
	core "github.com/cufee/aftermath-core/internal/core/stats"
	"github.com/cufee/aftermath-core/internal/logic/sessions"
	"go.mongodb.org/mongo-driver/mongo"
)

type SessionType string

const (
	SessionTypeDaily = SessionType("daily")
)

type SessionDatabaseRecord struct {
	ID        int         `bson:"_id,omitempty"`
	Type      SessionType `bson:"type"`
	CreatedAt time.Time   `bson:"createdAt"`

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
	Before time.Time
	Type   SessionType
}

func GetPlayerSessionSnapshot(accountID int, options ...SessionGetOptions) (*core.SessionSnapshot, error) {
	opts := SessionGetOptions{Type: SessionTypeDaily}
	if len(options) > 0 {
		opts = options[0]
	}

	_ = opts

	// TODO: save session to database

	return nil, errors.New("not implemented")
}
