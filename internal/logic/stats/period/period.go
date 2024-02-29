package stats

import (
	"strings"
	"time"

	core "github.com/cufee/aftermath-core/internal/core/stats"
	"github.com/cufee/aftermath-core/internal/logic/external/blitzstars"
	"github.com/cufee/aftermath-core/internal/logic/stats"
	"github.com/cufee/aftermath-core/utils"
	"github.com/gorhill/cronexpr"
)

type PeriodStats struct {
	Account stats.AccountWithClan

	Start time.Time `json:"start"`
	End   time.Time `json:"end"`

	Stats    *core.ReducedStatsFrame
	Vehicles map[int]*core.ReducedVehicleStats
}

const durationDay = time.Hour * 24

var sessionsCronNA = cronexpr.MustParse("0 9 * * *")
var sessionsCronEU = cronexpr.MustParse("0 1 * * *")
var sessionsCronAsia = cronexpr.MustParse("0 18 * * *")

func GetPlayerStats(accountId int, days int) (*PeriodStats, error) {
	cutoffTime := daysToRealmTime(utils.RealmFromAccountID(accountId), days)
	tankHistory, err := blitzstars.GetPlayerTankHistories(accountId)
	if err != nil {
		return nil, err
	}

	var periodStats PeriodStats
	periodStats.End = time.Now()
	periodStats.Start = cutoffTime

	for id, entries := range tankHistory {
		var tankTotal blitzstars.TankHistoryEntry
		for _, entry := range entries {
			if entry.LastBattleTime > int(cutoffTime.Unix()) {
				tankTotal.Add(entry)
			}
		}

		periodStats.Vehicles[id] = &core.ReducedVehicleStats{
			ReducedStatsFrame: stats.FrameToReducedStatsFrame(tankTotal.Stats),
			LastBattleTime:    tankTotal.LastBattleTime,
			MarkOfMastery:     tankTotal.MarkOfMastery,
			VehicleID:         id,
		}
		periodStats.Stats.Add(periodStats.Vehicles[id].ReducedStatsFrame)
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
