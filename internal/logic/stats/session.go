package stats

import (
	"errors"
	"sync"

	core "github.com/cufee/aftermath-core/internal/core/stats"
	"github.com/cufee/aftermath-core/internal/core/utils"
	"github.com/cufee/aftermath-core/internal/logic/cache"
	"github.com/cufee/aftermath-core/internal/logic/sessions"
)

var (
	ErrBadLiveSession = errors.New("bad live session")
)

type Snapshot struct {
	Selected *core.SessionSnapshot
	Live     *core.SessionSnapshot
}

func GetCurrentPlayerSession(realm string, accountId int, options ...cache.SessionGetOptions) (*Snapshot, error) {
	liveSessionChan := make(chan utils.DataWithError[*core.SessionSnapshot], 1)
	lastSessionChan := make(chan utils.DataWithError[*core.SessionSnapshot], 1)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		liveSessions, err := sessions.GetLiveSessions(realm, accountId)
		if err != nil {
			liveSessionChan <- utils.DataWithError[*core.SessionSnapshot]{Err: err}
			return
		}
		liveSession, ok := liveSessions[accountId]
		if !ok {
			liveSessionChan <- utils.DataWithError[*core.SessionSnapshot]{Err: ErrBadLiveSession}
			return
		}

		liveSessionChan <- utils.DataWithError[*core.SessionSnapshot]{Data: liveSession.Session}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		lastSession, err := cache.GetPlayerSessionSnapshot(accountId, options...)
		if errors.Is(err, cache.ErrNoSessionCache) {
			// If there is no session cache, we need to refresh the cache
			err := cache.RefreshSessions(cache.SessionTypeDaily, realm, accountId)
			if err != nil {
				lastSessionChan <- utils.DataWithError[*core.SessionSnapshot]{Err: err}
				return
			}

			lastSession := core.EmptySession(accountId, 0)
			lastSessionChan <- utils.DataWithError[*core.SessionSnapshot]{Data: lastSession}
			return
		} else if err != nil {
			lastSessionChan <- utils.DataWithError[*core.SessionSnapshot]{Err: err}
			return
		}

		lastSessionChan <- utils.DataWithError[*core.SessionSnapshot]{Data: lastSession}
	}()

	wg.Wait()
	close(liveSessionChan)
	close(lastSessionChan)

	liveSession := <-liveSessionChan
	if liveSession.Err != nil {
		return nil, liveSession.Err
	}
	lastSession := <-lastSessionChan
	if lastSession.Err != nil {
		return nil, lastSession.Err
	}

	selectedSession, err := liveSession.Data.Diff(lastSession.Data)
	if err != nil {
		return nil, err
	}

	return &Snapshot{
		Selected: selectedSession,
		Live:     liveSession.Data,
	}, nil
}
