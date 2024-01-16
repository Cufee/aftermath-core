package sessions

import (
	"github.com/cufee/aftermath-core/internal/core/stats"
	wg "github.com/cufee/am-wg-proxy-next/types"
)

func AccountStatsToSession(account wg.ExtendedAccount, accountVehicles []wg.VehicleStatsFrame) *stats.SessionSnapshot {
	session := &stats.SessionSnapshot{
		AccountID:      account.ID,
		LastBattleTime: account.LastBattleTime,
		Global:         frameToReducedStatsFrame(account.Statistics.All),
		Rating:         frameToReducedStatsFrame(account.Statistics.Rating),
		Vehicles:       make(map[int]*stats.ReducedVehicleStats),
	}

	for _, vehicle := range accountVehicles {
		session.Vehicles[vehicle.TankID] = &stats.ReducedVehicleStats{
			VehicleID:         vehicle.TankID,
			ReducedStatsFrame: frameToReducedStatsFrame(vehicle.Stats),
			MarkOfMastery:     vehicle.MarkOfMastery,
			LastBattleTime:    vehicle.LastBattleTime,
		}
	}
	return session
}

func frameToReducedStatsFrame(frame wg.StatsFrame) *stats.ReducedStatsFrame {
	return &stats.ReducedStatsFrame{
		Battles:              frame.Battles,
		BattlesWon:           frame.Wins,
		BattlesSurvived:      frame.SurvivedBattles,
		DamageDealt:          frame.DamageDealt,
		DamageReceived:       frame.DamageReceived,
		ShotsHit:             frame.Hits,
		ShotsFired:           frame.Shots,
		Frags:                frame.Frags,
		MaxFrags:             frame.MaxFrags,
		EnemiesSpotted:       frame.Spotted,
		CapturePoints:        frame.CapturePoints,
		DroppedCapturePoints: frame.DroppedCapturePoints,
	}
}
