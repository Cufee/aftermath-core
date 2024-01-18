package database

import (
	"github.com/cufee/aftermath-core/internal/core/database/models"
	"github.com/cufee/aftermath-core/internal/core/stats"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func UpdateAverages(averages map[int]stats.ReducedStatsFrame) error {
	var writes []mongo.WriteModel
	for id, average := range averages {
		data := models.TankAverages{
			ID:                id,
			ReducedStatsFrame: average,
		}

		model := mongo.NewUpdateOneModel()
		model.SetFilter(bson.M{"_id": id})
		model.SetUpdate(bson.M{"$set": data})
		model.SetUpsert(true)
		writes = append(writes, model)
	}

	ctx, cancel := DefaultClient.Ctx()
	defer cancel()

	_, err := DefaultClient.Collection(CollectionVehicleAverages).BulkWrite(ctx, writes)
	if err != nil {
		return err
	}

	return nil
}

func GetVehicleAverages(vehicleIDs ...int) (map[int]*stats.ReducedStatsFrame, error) {
	if len(vehicleIDs) == 0 {
		return nil, nil
	}

	ctx, cancel := DefaultClient.Ctx()
	defer cancel()

	var averages []models.TankAverages
	cur, err := DefaultClient.Collection(CollectionVehicleAverages).Find(ctx, bson.M{"_id": bson.M{"$in": vehicleIDs}})
	if err != nil {
		return nil, err
	}

	err = cur.All(ctx, &averages)
	if err != nil {
		return nil, err
	}

	averageMap := make(map[int]*stats.ReducedStatsFrame)
	for _, average := range averages {
		averageMap[average.ID] = &average.ReducedStatsFrame
	}

	return averageMap, nil
}

func GetGlossaryVehicles(vehicleIDs ...int) (map[int]models.Vehicle, error) {
	if len(vehicleIDs) == 0 {
		return nil, nil
	}

	ctx, cancel := DefaultClient.Ctx()
	defer cancel()

	var vehicles []models.Vehicle
	cur, err := DefaultClient.Collection(CollectionVehicleGlossary).Find(ctx, bson.M{"_id": bson.M{"$in": vehicleIDs}})
	if err != nil {
		return nil, err
	}

	err = cur.All(ctx, &vehicles)
	if err != nil {
		return nil, err
	}

	vehicleMap := make(map[int]models.Vehicle)
	for _, vehicle := range vehicles {
		vehicleMap[vehicle.ID] = vehicle
	}

	return vehicleMap, nil
}
