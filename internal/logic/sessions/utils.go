package sessions

import (
	"github.com/cufee/aftermath-core/internal/core/stats"
	wg "github.com/cufee/am-wg-proxy-next/types"
)

func BatchAccountIDs(accountIDs []int, batchSize int) [][]int {
	batches := len(accountIDs) / batchSize
	if len(accountIDs)%100 != 0 {
		batches++
	}

	batchedAccountIDs := make([][]int, batches)
	for i := 0; i < batches; i++ {
		start := i * batchSize
		end := start + batchSize
		if end > len(accountIDs) {
			end = len(accountIDs)
		}

		batchedAccountIDs[i] = accountIDs[start:end]
	}

	return batchedAccountIDs
}

func AccountStatsToSession(account wg.ExtendedAccount, accountVehicles []wg.VehicleStatsFrame) *stats.SessionSnapshot {
	session := &stats.SessionSnapshot{
		AccountID:      account.ID,
		LastBattleTime: account.LastBattleTime,
		Global:         frameToReducedStatsFrame(account.Statistics.All),
		Rating:         frameToReducedStatsFrame(account.Statistics.Rating),
		Vehicles:       make(map[int]stats.ReducedVehicleStats),
	}

	for _, vehicle := range accountVehicles {
		session.Vehicles[vehicle.TankID] = stats.ReducedVehicleStats{
			VehicleID:         vehicle.TankID,
			ReducedStatsFrame: frameToReducedStatsFrame(vehicle.Stats),
			MarkOfMastery:     vehicle.MarkOfMastery,
			LastBattleTime:    vehicle.LastBattleTime,
		}
	}
	return session
}

func frameToReducedStatsFrame(frame wg.StatsFrame) stats.ReducedStatsFrame {
	return stats.ReducedStatsFrame{
		Battles:              frame.Battles,
		BattlesWon:           frame.Wins,
		DamageDealt:          frame.DamageDealt,
		ShotsHit:             frame.Hits,
		ShotsFired:           frame.Shots,
		Frags:                frame.Frags,
		EnemiesSpotted:       frame.Spotted,
		DroppedCapturePoints: frame.DroppedCapturePoints,
	}
}
