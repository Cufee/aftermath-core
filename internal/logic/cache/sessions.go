package cache

import (
	"github.com/cufee/aftermath-core/internal/core/database"
	"github.com/cufee/aftermath-core/internal/core/database/models"
	"github.com/cufee/aftermath-core/internal/core/stats"
	"github.com/cufee/aftermath-core/internal/core/wargaming"
	"github.com/cufee/aftermath-core/internal/logic/sessions"
	wg "github.com/cufee/am-wg-proxy-next/types"
	"github.com/rs/zerolog/log"
)

func RefreshSessionsAndAccounts(sessionType models.SessionType, realm string, accountIDs ...int) (map[int]error, error) {
	sessions, err := sessions.GetSessionsWithClient(wargaming.Clients.Cache, realm, accountIDs...)
	if err != nil {
		return nil, err
	}

	var accounts []*wg.ExtendedAccount
	for _, session := range sessions {
		if session.Err != nil {
			continue
		}
		accounts = append(accounts, session.Data.Account)
	}
	err = UpdatePlayerAccountsFromWG(realm, accounts...)
	if err != nil {
		return nil, err
	}

	lastBattle, err := database.GetLastBattleTimes(sessionType, accountIDs...)
	if err != nil {
		return nil, err
	}

	updateErrors := make(map[int]error)
	var snapshots []*stats.SessionSnapshot
	for accountId, session := range sessions {
		if session.Err != nil {
			updateErrors[accountId] = session.Err
			continue
		}

		if lastBattle[session.Data.Account.ID] >= session.Data.Session.LastBattleTime {
			log.Debug().Msgf("%d played 0 battles since last session, skipping update", session.Data.Account.ID)
			continue
		}
		snapshots = append(snapshots, session.Data.Session)
	}
	if len(snapshots) == 0 {
		return updateErrors, nil
	}

	return updateErrors, database.InsertSession(models.SessionTypeDaily, snapshots...)
}
