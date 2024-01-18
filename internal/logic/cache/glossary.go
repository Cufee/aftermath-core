package cache

import (
	"github.com/cufee/aftermath-core/internal/core/database"
	"github.com/cufee/aftermath-core/internal/core/database/models"
	"github.com/cufee/aftermath-core/internal/logic/external"
)

func UpdateGlossaryCache() error {
	vehiclesMap, err := external.GetCompleteVehicleGlossary()
	if err != nil {
		return err
	}

	var vehicles []models.Vehicle
	for _, vehicle := range vehiclesMap {
		vehicles = append(vehicles, vehicle)
	}

	return database.UpdateGlossary(vehicles)
}
