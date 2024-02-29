package sessions

import (
	"errors"

	"github.com/cufee/aftermath-core/internal/core/database"
	core "github.com/cufee/aftermath-core/internal/core/stats"

	"github.com/cufee/aftermath-core/internal/logic/stats"
	"github.com/rs/zerolog/log"
)

func GetCurrentPlayerSession(realm string, accountId int, options ...database.SessionGetOptions) (Snapshot, error) {
	opts := database.SessionGetOptions{}
	if len(options) > 0 {
		opts = options[0]
	}

	var snapshot Snapshot

	liveSessions, err := GetLiveSessions(realm, accountId)
	if err != nil {
		log.Err(err).Msg("failed to get live sessions")
		return snapshot, err
	}
	liveSession, ok := liveSessions[accountId]
	if !ok {
		log.Err(ErrBadLiveSession).Msg("failed to get live session")
		return snapshot, ErrBadLiveSession
	}
	if liveSession.Err != nil {
		log.Err(liveSession.Err).Msg("failed to get live session")
		return snapshot, liveSession.Err
	}

	snapshot.Live = liveSession.Data.Session
	snapshot.Account = stats.AccountWithClan{
		ExtendedAccount: *liveSession.Data.Account,
		ClanMember:      *liveSession.Data.Clan,
	}

	if opts.LastBattleBefore == nil && liveSession.Data.Account.LastBattleTime > 0 {
		opts.LastBattleBefore = &liveSession.Data.Account.LastBattleTime
	}

	lastSession, err := database.GetPlayerSessionSnapshot(accountId, opts)
	if errors.Is(err, database.ErrNoSessionCache) {
		// There is no session cache, so the live session is the same as the last session and there is no diff
		snapshot.Diff = core.EmptySession(liveSession.Data.Account.ID, liveSession.Data.Account.LastBattleTime)
		snapshot.Selected = liveSession.Data.Session
		return snapshot, ErrNoSessionCached
	}
	// All other errors
	if err != nil {
		log.Err(err).Msg("failed to get last session")
		return snapshot, err
	}

	diffSession, err := liveSession.Data.Session.Diff(lastSession.Session)
	if err != nil {
		log.Err(err).Msg("failed to diff sessions")
		return snapshot, err
	}

	// Clean up vehicles with 0 battles
	for _, vehicle := range diffSession.Vehicles {
		if vehicle.Battles == 0 {
			delete(diffSession.Vehicles, vehicle.VehicleID)
		}
	}

	snapshot.Selected = lastSession.Session
	snapshot.Diff = diffSession
	return snapshot, nil
}
