package stats

import (
	"errors"

	core "github.com/cufee/aftermath-core/internal/core/stats"
	"github.com/cufee/aftermath-core/internal/logic/cache"
	"github.com/cufee/aftermath-core/internal/logic/sessions"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrNoSessionCache = errors.New("no session found")
	ErrBadLiveSession = errors.New("bad live session")
)

func GetCurrentPlayerSession(realm string, accountId int) (*core.SessionSnapshot, error) {
	snapshot, err := cache.GetPlayerSessionSnapshot(accountId)
	if err != nil {
		if errors.Is(mongo.ErrNoDocuments, err) {
			return nil, ErrNoSessionCache
		}
		return nil, err
	}

	liveSessions, err := sessions.GetLiveSessions(realm, accountId)
	if err != nil {
		return nil, err
	}

	liveSession, ok := liveSessions[accountId]
	if !ok {
		return nil, ErrBadLiveSession
	}

	liveSession.Session.Subtract(snapshot)
	return liveSession.Session, nil
}
