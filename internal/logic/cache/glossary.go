package cache

import (
	"github.com/cufee/aftermath-core/internal/core/database"
	"github.com/cufee/aftermath-core/internal/core/database/models"
	"github.com/cufee/aftermath-core/internal/logic/external/wotinspector"
)

func UpdateGlossaryCache() error {
	vehiclesMap, err := wotinspector.GetCompleteVehicleGlossary()
	if err != nil {
		return err
	}

	var vehicles []models.Vehicle
	for _, vehicle := range vehiclesMap {
		vehicles = append(vehicles, vehicle)
	}

	return database.UpdateGlossary(vehicles)
}
