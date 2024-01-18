package cache

import (
	"github.com/cufee/aftermath-core/internal/core/database"
	"github.com/cufee/aftermath-core/internal/logic/external"
)

func UpdateAveragesCache() error {
	averages, err := external.GetTankAverages()
	if err != nil {
		return err
	}
	return database.UpdateAverages(averages)
}
