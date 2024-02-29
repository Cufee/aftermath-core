package sessions

import (
	core "github.com/cufee/aftermath-core/internal/core/stats"
	"github.com/cufee/aftermath-core/internal/logic/stats"
	wg "github.com/cufee/am-wg-proxy-next/types"
)

func AccountStatsToSession(account wg.ExtendedAccount, accountVehicles []wg.VehicleStatsFrame) *core.SessionSnapshot {
	session := &core.SessionSnapshot{
		AccountID:      account.ID,
		LastBattleTime: account.LastBattleTime,
		Global:         stats.FrameToReducedStatsFrame(account.Statistics.All),
		Rating:         stats.FrameToReducedStatsFrame(account.Statistics.Rating),
		Vehicles:       make(map[int]*core.ReducedVehicleStats),
	}

	for _, vehicle := range accountVehicles {
		session.Vehicles[vehicle.TankID] = &core.ReducedVehicleStats{
			VehicleID:         vehicle.TankID,
			ReducedStatsFrame: stats.FrameToReducedStatsFrame(vehicle.Stats),
			MarkOfMastery:     vehicle.MarkOfMastery,
			LastBattleTime:    vehicle.LastBattleTime,
		}
	}
	return session
}
