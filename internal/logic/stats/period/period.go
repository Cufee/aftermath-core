package period

import (
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/cufee/aftermath-core/internal/core/database"
	core "github.com/cufee/aftermath-core/internal/core/stats"
	"github.com/cufee/aftermath-core/internal/core/wargaming"
	"github.com/cufee/aftermath-core/internal/logic/external/blitzstars"
	"github.com/cufee/aftermath-core/internal/logic/stats"

	"github.com/cufee/am-wg-proxy-next/v2/types"

	"github.com/gorhill/cronexpr"
)

type PeriodStats struct {
	Account types.Account `json:"account"`
	Clan    types.Clan    `json:"clan"`

	Start time.Time `json:"start"`
	End   time.Time `json:"end"`

	Vehicles map[int]core.ReducedVehicleStats `json:"vehicles"`
	Stats    core.ReducedStatsFrame           `json:"stats"`
}

func (stats *PeriodStats) CareerWN8(averages map[int]core.ReducedStatsFrame) int {
	if v := stats.Stats.WN8(); v != core.InvalidValueInt {
		return v
	}

	var weightedWN8Total, wn8BattlesTotal int
	for id, vehicle := range stats.Vehicles {
		if vehicle.Battles < 1 {
			continue
		}
		if data, ok := averages[id]; ok {
			weightedWN8Total += vehicle.Battles * vehicle.WN8(data)
			wn8BattlesTotal += vehicle.Battles
		}
	}
	if wn8BattlesTotal < 1 {
		return core.InvalidValueInt
	}

	v := weightedWN8Total / wn8BattlesTotal
	stats.Stats.SetWN8(v)
	return v
}

const durationDay = time.Hour * 24

var sessionsCronNA = cronexpr.MustParse("0 9 * * *")
var sessionsCronEU = cronexpr.MustParse("0 1 * * *")
var sessionsCronAsia = cronexpr.MustParse("0 18 * * *")

func GetPlayerStats(accountId int, days int) (PeriodStats, error) {
	realm := wargaming.Clients.Live.RealmFromAccountID(strconv.Itoa(accountId))
	allStats, err := stats.GetCompleteStatsWithClient(wargaming.Clients.Live, realm, accountId)
	if err != nil {
		return PeriodStats{}, err
	}
	accountStats, ok := allStats[accountId]
	if !ok {
		return PeriodStats{}, stats.ErrBlankResponse
	}
	if accountStats.Err != nil {
		return PeriodStats{}, accountStats.Err
	}

	var vehicleIDs []int
	for _, vehicle := range accountStats.Data.Vehicles {
		vehicleIDs = append(vehicleIDs, vehicle.TankID)
	}

	tankAverages, err := database.GetVehicleAverages(vehicleIDs...)
	if err != nil {
		return PeriodStats{}, err
	}

	var periodStats = PeriodStats{
		Clan:     accountStats.Data.Clan.Clan,
		Account:  accountStats.Data.Account.Account,
		Vehicles: make(map[int]core.ReducedVehicleStats),
		End:      time.Unix(int64(accountStats.Data.Account.LastBattleTime), 0),
	}

	switch {
	case days <= 0:
		fallthrough
	case days > 90:
		// Return career stats
		for _, vehicle := range accountStats.Data.Vehicles {
			frame := stats.FrameToReducedStatsFrame(vehicle.Stats)
			stats := core.ReducedVehicleStats{
				ReducedStatsFrame: &frame,
				LastBattleTime:    vehicle.LastBattleTime,
				MarkOfMastery:     vehicle.MarkOfMastery,
				VehicleID:         vehicle.TankID,
			}
			stats.WN8(tankAverages[vehicle.TankID])
			periodStats.Vehicles[vehicle.TankID] = stats
		}

		// TODO: Get calculated career WN8

		periodStats.Start = time.Unix(int64(accountStats.Data.Account.CreatedAt), 0)
		periodStats.Stats = accountStats.Data.Session.Global

		periodStats.CareerWN8(tankAverages)
		return periodStats, nil

	default:
		// Get time specific stats
		periodStats.Start = daysToRealmTime(realm, days)
		if periodStats.End.Before(periodStats.Start) {
			periodStats.End = time.Now()
		}
	}

	tankHistory, err := blitzstars.GetPlayerTankHistories(accountId)
	if err != nil {
		return PeriodStats{}, err
	}

	var vehiclesMap = make(map[int]types.VehicleStatsFrame)
	for _, vehicle := range accountStats.Data.Vehicles {
		if vehicle.LastBattleTime < int(periodStats.Start.Unix()) {
			continue
		}
		vehiclesMap[vehicle.TankID] = vehicle
	}

	for _, vehicle := range accountStats.Data.Vehicles {
		if vehicle.LastBattleTime < int(periodStats.Start.Unix()) {
			continue
		}

		entries := tankHistory[vehicle.TankID]
		// Sort entries by number of battles in descending order
		slices.SortFunc(entries, func(i, j blitzstars.TankHistoryEntry) int {
			return j.Stats.Battles - i.Stats.Battles
		})

		var selectedEntry blitzstars.TankHistoryEntry
		for _, entry := range entries {
			if entry.LastBattleTime < int(periodStats.Start.Unix()) {
				selectedEntry = entry
				break
			}
		}

		if selectedEntry.Stats.Battles < vehicle.Stats.Battles {
			selectedFrame := stats.FrameToReducedStatsFrame(selectedEntry.Stats)
			compareToFrame := stats.FrameToReducedStatsFrame(vehicle.Stats)
			compareToFrame.Subtract(selectedFrame)

			frame := core.ReducedVehicleStats{
				ReducedStatsFrame: &compareToFrame,
				LastBattleTime:    vehicle.LastBattleTime,
				VehicleID:         vehicle.TankID,
			}
			frame.WN8(tankAverages[vehicle.TankID])

			periodStats.Vehicles[vehicle.TankID] = frame
			periodStats.Stats.Add(*frame.ReducedStatsFrame)
		}
	}

	periodStats.CareerWN8(tankAverages)
	return periodStats, nil
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
