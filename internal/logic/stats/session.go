package stats

import (
	"errors"
	"sort"

	"github.com/cufee/aftermath-core/internal/core/database"
	"github.com/cufee/aftermath-core/internal/core/database/models"
	core "github.com/cufee/aftermath-core/internal/core/stats"
	"github.com/cufee/aftermath-core/internal/logic/cache"
	"github.com/cufee/aftermath-core/internal/logic/sessions"
	wg "github.com/cufee/am-wg-proxy-next/types"
	"github.com/rs/zerolog/log"
)

var (
	ErrBadLiveSession = errors.New("bad live session")
)

type SnapshotAccount struct {
	wg.ExtendedAccount
	wg.ClanMember
}

type Snapshot struct {
	Account  SnapshotAccount
	Selected *core.SessionSnapshot // The session that was selected from the database
	Live     *core.SessionSnapshot // The live session
	Diff     *core.SessionSnapshot // The difference between the selected and live sessions
}

func GetCurrentPlayerSession(realm string, accountId int, options ...database.SessionGetOptions) (*Snapshot, error) {
	opts := database.SessionGetOptions{}
	if len(options) > 0 {
		opts = options[0]
	}

	liveSessions, err := sessions.GetLiveSessions(realm, accountId)
	if err != nil {
		log.Err(err).Msg("failed to get live sessions")
		return nil, err
	}
	liveSession, ok := liveSessions[accountId]
	if !ok {
		log.Err(ErrBadLiveSession).Msg("failed to get live session")
		return nil, ErrBadLiveSession
	}
	if liveSession.Err != nil {
		log.Err(liveSession.Err).Msg("failed to get live session")
		return nil, liveSession.Err
	}

	if opts.LastBattleBefore == nil && liveSession.Data.Account.LastBattleTime > 0 {
		opts.LastBattleBefore = &liveSession.Data.Account.LastBattleTime
	}

	lastSession, err := database.GetPlayerSessionSnapshot(accountId, opts)
	if errors.Is(err, database.ErrNoSessionCache) {
		// Refresh the session cache in the background
		go func(realm string, accountId int) {
			accountErrs, err := cache.RefreshSessionsAndAccounts(models.SessionTypeDaily, realm, accountId)
			if err != nil || len(accountErrs) > 0 {
				log.Err(err).Msg("failed to refresh session cache")
			}
		}(realm, accountId)

		// There is no session cache, so the live session is the same as the last session and there is no diff
		return &Snapshot{
			Selected: liveSession.Data.Session,
			Account: SnapshotAccount{
				ExtendedAccount: *liveSession.Data.Account,
				ClanMember:      *liveSession.Data.Clan,
			},
			Live: liveSession.Data.Session,
			Diff: core.EmptySession(liveSession.Data.Account.ID, liveSession.Data.Account.LastBattleTime),
		}, nil
	}
	// All other errors
	if err != nil {
		log.Err(err).Msg("failed to get last session")
		return nil, err
	}

	diffSession, err := liveSession.Data.Session.Diff(lastSession.Session)
	if err != nil {
		log.Err(err).Msg("failed to diff sessions")
		return nil, err
	}

	// Clean up vehicles with 0 battles
	for _, vehicle := range diffSession.Vehicles {
		if vehicle.Battles == 0 {
			delete(diffSession.Vehicles, vehicle.VehicleID)
		}
	}

	return &Snapshot{
		Selected: lastSession.Session,
		Account: SnapshotAccount{
			ExtendedAccount: *liveSession.Data.Account,
			ClanMember:      *liveSession.Data.Clan,
		},
		Live: liveSession.Data.Session,
		Diff: diffSession,
	}, nil
}

type SortOptions struct {
	By    vehicleSortOptions
	Limit int
}
type vehicleSortOptions string

const (
	SortByBattlesDesc   = vehicleSortOptions("-battles")
	SortByBattlesAsc    = vehicleSortOptions("battles")
	SortByWinrateDesc   = vehicleSortOptions("-winrate")
	SortByWinrateAsc    = vehicleSortOptions("winrate")
	SortByWN8Desc       = vehicleSortOptions("-wn8")
	SortByWN8Asc        = vehicleSortOptions("wn8")
	SortByAvgDamageDesc = vehicleSortOptions("-avgDamage")
	SortByAvgDamageAsc  = vehicleSortOptions("avgDamage")
	SortByLastBattle    = vehicleSortOptions("lastBattleTime")
)

func ParseSortOptions(sort string) vehicleSortOptions {
	switch sort {
	case string(SortByBattlesDesc):
		return SortByBattlesDesc
	case string(SortByBattlesAsc):
		return SortByBattlesAsc
	case string(SortByWinrateDesc):
		return SortByWinrateDesc
	case string(SortByWinrateAsc):
		return SortByWinrateAsc
	case string(SortByWN8Desc):
		return SortByWN8Desc
	case string(SortByWN8Asc):
		return SortByWN8Asc
	case string(SortByAvgDamageDesc):
		return SortByAvgDamageDesc
	case string(SortByAvgDamageAsc):
		return SortByAvgDamageAsc
	case string(SortByLastBattle):
		return SortByLastBattle
	default:
		return SortByLastBattle
	}
}

func GetVehicleAverages(vehicles map[int]*core.ReducedVehicleStats) (map[int]*core.ReducedStatsFrame, error) {
	var vehicleIDs []int
	for _, vehicle := range vehicles {
		vehicleIDs = append(vehicleIDs, vehicle.VehicleID)
	}
	return database.GetVehicleAverages(vehicleIDs...)
}

func SortVehicles(vehicles map[int]*core.ReducedVehicleStats, averages map[int]*core.ReducedStatsFrame, options ...SortOptions) []*core.ReducedVehicleStats {
	opts := SortOptions{By: SortByLastBattle, Limit: 10}
	if len(options) > 0 {
		opts = options[0]
	}

	var sorted []*core.ReducedVehicleStats
	for _, vehicle := range vehicles {
		sorted = append(sorted, vehicle)
	}

	switch opts.By {
	case SortByBattlesDesc:
		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i].Battles > sorted[j].Battles
		})
	case SortByBattlesAsc:
		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i].Battles < sorted[j].Battles
		})
	case SortByWinrateDesc:
		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i].Winrate() > sorted[j].Winrate()
		})
	case SortByWinrateAsc:
		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i].Winrate() < sorted[j].Winrate()
		})
	case SortByWN8Desc:
		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i].WN8(averages[sorted[i].VehicleID]) > sorted[j].WN8(averages[sorted[j].VehicleID])
		})
	case SortByWN8Asc:
		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i].WN8(averages[sorted[i].VehicleID]) < sorted[j].WN8(averages[sorted[j].VehicleID])
		})
	case SortByAvgDamageDesc:
		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i].AvgDamage() > sorted[j].AvgDamage()
		})
	case SortByAvgDamageAsc:
		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i].AvgDamage() < sorted[j].AvgDamage()
		})
	case SortByLastBattle:
		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i].LastBattleTime > sorted[j].LastBattleTime
		})
	}

	if opts.Limit > 0 && opts.Limit < len(sorted) {
		sorted = sorted[:opts.Limit]
	}
	return sorted
}
