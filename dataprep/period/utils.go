package period

import (
	"github.com/cufee/aftermath-core/internal/core/database/models"
	"github.com/cufee/aftermath-core/internal/core/stats"
)

func calculateAvgTier(vehicles map[int]*stats.ReducedVehicleStats, glossary map[int]models.Vehicle) float32 {
	var battlesTotal int
	var weightedTierTotal float32

	for _, vehicle := range vehicles {
		if data, ok := glossary[vehicle.VehicleID]; ok && data.Tier > 0 {
			battlesTotal += vehicle.Battles
			weightedTierTotal += float32(vehicle.Battles * data.Tier)
		}
	}
	if battlesTotal == 0 {
		return stats.InvalidValueFloat32
	}

	return weightedTierTotal / float32(battlesTotal)
}
