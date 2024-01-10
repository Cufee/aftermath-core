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
	Selected *core.SessionSnapshot // The session that was selected from the database
	Live     *core.SessionSnapshot // The live session
	Diff     *core.SessionSnapshot // The difference between the selected and live sessions
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
		lastSessionChan <- utils.DataWithError[*core.SessionSnapshot]{Data: lastSession, Err: err}
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
		if errors.Is(lastSession.Err, cache.ErrNoSessionCache) {
			go cache.RefreshSessions(cache.SessionTypeDaily, realm, accountId) // Refresh the session cache in the background
			// There is no session cache, so the live session is the same as the last session and there is no diff
			return &Snapshot{
				Selected: liveSession.Data,
				Live:     liveSession.Data,
				Diff:     core.EmptySession(liveSession.Data.AccountID, liveSession.Data.LastBattleTime),
			}, nil
		}
		return nil, lastSession.Err
	}

	diffSession, err := liveSession.Data.Diff(lastSession.Data)
	if err != nil {
		return nil, err
	}

	return &Snapshot{
		Selected: lastSession.Data,
		Live:     liveSession.Data,
		Diff:     diffSession,
	}, nil
}
