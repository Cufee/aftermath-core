package stats

import (
	"github.com/cufee/aftermath-core/internal/core/database"
	core "github.com/cufee/aftermath-core/internal/core/stats"
)

func GetVehicleAverages(vehicles map[int]*core.ReducedVehicleStats) (map[int]*core.ReducedStatsFrame, error) {
	var vehicleIDs []int
	for _, vehicle := range vehicles {
		vehicleIDs = append(vehicleIDs, vehicle.VehicleID)
	}
	return database.GetVehicleAverages(vehicleIDs...)
}
