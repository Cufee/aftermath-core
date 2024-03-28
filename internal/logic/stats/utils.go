package stats

import (
	"github.com/cufee/aftermath-core/internal/core/stats"
	"github.com/cufee/am-wg-proxy-next/types"
)

func FrameToReducedStatsFrame(frame types.StatsFrame) stats.ReducedStatsFrame {
	return stats.ReducedStatsFrame{
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
		RawRating:            frame.Rating,
	}
}

func CompleteStatsFromWargaming(account types.ExtendedAccount, accountVehicles []types.VehicleStatsFrame) stats.SessionSnapshot {
	session := stats.SessionSnapshot{
		AccountID:      account.ID,
		LastBattleTime: account.LastBattleTime,
		Global:         FrameToReducedStatsFrame(account.Statistics.All),
		Rating:         FrameToReducedStatsFrame(account.Statistics.Rating),
		Vehicles:       make(map[int]stats.ReducedVehicleStats),
	}

	for _, vehicle := range accountVehicles {
		frame := FrameToReducedStatsFrame(vehicle.Stats)
		session.Vehicles[vehicle.TankID] = stats.ReducedVehicleStats{
			VehicleID:         vehicle.TankID,
			ReducedStatsFrame: &frame,
			MarkOfMastery:     vehicle.MarkOfMastery,
			LastBattleTime:    vehicle.LastBattleTime,
		}
	}
	return session
}
