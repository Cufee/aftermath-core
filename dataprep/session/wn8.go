package session

import (
	"github.com/cufee/aftermath-core/internal/core/stats"
	wg "github.com/cufee/am-wg-proxy-next/v2/types"
)

type SessionStats struct {
	Realm      string     `json:"realm"`
	Locale     string     `json:"locale"`
	LastBattle int        `json:"last_battle"`
	Clan       wg.Clan    `json:"clan"`
	Account    wg.Account `json:"account"`
	Cards      Cards      `json:"cards"`
}

func calculateWeightedWN8(vehicles map[int]stats.ReducedVehicleStats, averages map[int]stats.ReducedStatsFrame) int {
	var wn8VehiclesTotal, wn8VehiclesBattles int
	for id, vehicle := range vehicles {
		if vehicle.Valid(vehicle.WN8(averages[id])) {
			wn8VehiclesTotal += vehicle.WN8(averages[id]) * vehicle.Battles
			wn8VehiclesBattles += vehicle.Battles
		}
	}
	if wn8VehiclesBattles > 0 {
		return wn8VehiclesTotal / wn8VehiclesBattles
	}
	return stats.InvalidValueInt
}
