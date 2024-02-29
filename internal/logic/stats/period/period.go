package period

import (
	"slices"
	"strings"
	"time"

	core "github.com/cufee/aftermath-core/internal/core/stats"
	"github.com/cufee/aftermath-core/internal/core/wargaming"
	"github.com/cufee/aftermath-core/internal/logic/external/blitzstars"
	"github.com/cufee/aftermath-core/internal/logic/stats"
	"github.com/cufee/aftermath-core/utils"
	"github.com/cufee/am-wg-proxy-next/types"
	"github.com/gorhill/cronexpr"
)

type PeriodStats struct {
	Account types.Account `json:"account"`
	Clan    types.Clan    `json:"clan"`

	Start time.Time `json:"start"`
	End   time.Time `json:"end"`

	Stats    core.ReducedStatsFrame            `json:"stats"`
	Vehicles map[int]*core.ReducedVehicleStats `json:"vehicles"`
}

const durationDay = time.Hour * 24

var sessionsCronNA = cronexpr.MustParse("0 9 * * *")
var sessionsCronEU = cronexpr.MustParse("0 1 * * *")
var sessionsCronAsia = cronexpr.MustParse("0 18 * * *")

func GetPlayerStats(accountId int, days int) (*PeriodStats, error) {
	realm := utils.RealmFromAccountID(accountId)
	allStats, err := stats.GetCompleteStatsWithClient(wargaming.Clients.Live, realm, accountId)
	if err != nil {
		return nil, err
	}
	accountStats, ok := allStats[accountId]
	if !ok {
		return nil, stats.ErrBlankResponse
	}
	if accountStats.Err != nil {
		return nil, accountStats.Err
	}

	var periodStats = PeriodStats{
		End: time.Unix(int64(accountStats.Data.Account.LastBattleTime), 0),

		Clan:     accountStats.Data.Clan.Clan,
		Account:  accountStats.Data.Account.Account,
		Vehicles: make(map[int]*core.ReducedVehicleStats),
	}

	var cutoffTime time.Time
	switch {
	case days <= 0:
		fallthrough
	case days > 90:
		// Return career stats
		for _, vehicle := range accountStats.Data.Vehicles {
			periodStats.Vehicles[vehicle.TankID] = &core.ReducedVehicleStats{
				ReducedStatsFrame: stats.FrameToReducedStatsFrame(vehicle.Stats),
				LastBattleTime:    vehicle.LastBattleTime,
				MarkOfMastery:     vehicle.MarkOfMastery,
				VehicleID:         vehicle.TankID,
			}
		}
		periodStats.Start = time.Unix(int64(accountStats.Data.Account.CreatedAt), 0)
		periodStats.Stats = *accountStats.Data.Session.Global
		return &periodStats, nil

	default:
		// Get time specific stats
		periodStats.Start = daysToRealmTime(realm, days)
	}

	tankHistory, err := blitzstars.GetPlayerTankHistories(accountId)
	if err != nil {
		return nil, err
	}

	var vehiclesMap = make(map[int]types.VehicleStatsFrame)
	for _, vehicle := range accountStats.Data.Vehicles {
		if vehicle.LastBattleTime < int(periodStats.Start.Unix()) {
			continue
		}
		vehiclesMap[vehicle.TankID] = vehicle
	}

	for id, entries := range tankHistory {
		vehicle, ok := vehiclesMap[id]
		if !ok {
			continue
		}

		// Sort entries by number of battles in descending order
		slices.SortFunc(entries, func(i, j blitzstars.TankHistoryEntry) int {
			return j.Stats.Battles - i.Stats.Battles
		})

		var selectedEntry blitzstars.TankHistoryEntry
		for _, entry := range entries {
			if entry.LastBattleTime < int(cutoffTime.Unix()) {
				selectedEntry = entry
				break
			}
		}

		if selectedEntry.Stats.Battles < vehicle.Stats.Battles {
			selectedFrame := stats.FrameToReducedStatsFrame(selectedEntry.Stats)
			compareToFrame := stats.FrameToReducedStatsFrame(vehicle.Stats)
			compareToFrame.Subtract(selectedFrame)

			frame := core.ReducedVehicleStats{
				LastBattleTime:    vehicle.LastBattleTime,
				ReducedStatsFrame: compareToFrame,
				VehicleID:         id,
			}
			periodStats.Vehicles[id] = &frame
			periodStats.Stats.Add(frame.ReducedStatsFrame)
		}
	}

	return &periodStats, nil
}

func daysToRealmTime(realm string, days int) time.Time {
	duration := durationDay * time.Duration(days)

	switch strings.ToLower(realm) {
	case "na":
		return sessionsCronNA.Next(time.Now()).Add(durationDay * -1).Add(-duration)
	case "eu":
		return sessionsCronEU.Next(time.Now()).Add(durationDay * -1).Add(-duration)
	case "as":
		return sessionsCronAsia.Next(time.Now()).Add(durationDay * -1).Add(-duration)
	default:
		return time.Now()
	}
}
