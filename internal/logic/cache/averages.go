package cache

import (
	"github.com/cufee/aftermath-core/internal/core/database"
	"github.com/cufee/aftermath-core/internal/logic/external/blitzstars"
)

func UpdateAveragesCache() error {
	averages, err := blitzstars.GetTankAverages()
	if err != nil {
		return err
	}
	return database.UpdateAverages(averages)
}
