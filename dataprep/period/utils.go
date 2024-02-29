package period

import "github.com/cufee/aftermath-core/internal/core/stats"

func calculateAvgTier(vehicles map[int]*stats.ReducedVehicleStats) float32 {
	var battlesTotal int
	var weightedTierTotal float32

	for _, vehicle := range vehicles {
		battlesTotal += vehicle.Battles
		weightedTierTotal += float32(vehicle.Battles) * 10.0
	}

	return weightedTierTotal / float32(battlesTotal)
}
