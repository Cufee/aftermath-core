package cache

import (
	"github.com/cufee/aftermath-core/internal/core/database"
	"github.com/cufee/aftermath-core/internal/core/stats"
	"github.com/cufee/aftermath-core/internal/logic/external"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type TankAverages struct {
	ID                      int `json:"id" bson:"_id"`
	stats.ReducedStatsFrame `bson:",inline"`
}

func UpdateAveragesCache() error {
	averages, err := external.GetTankAverages()
	if err != nil {
		return err
	}

	var writes []mongo.WriteModel
	for id, average := range averages {
		data := TankAverages{
			ID:                id,
			ReducedStatsFrame: average,
		}

		model := mongo.NewUpdateOneModel()
		model.SetFilter(bson.M{"_id": id})
		model.SetUpdate(bson.M{"$set": data})
		model.SetUpsert(true)
		writes = append(writes, model)
	}

	ctx, cancel := database.DefaultClient.Ctx()
	defer cancel()

	_, err = database.DefaultClient.Collection(database.CollectionVehicleAverages).BulkWrite(ctx, writes)
	if err != nil {
		return err
	}

	return nil
}

// TODO: add in-memory cache

func GetVehicleAverages(vehicleIDs ...int) (map[int]*stats.ReducedStatsFrame, error) {
	ctx, cancel := database.DefaultClient.Ctx()
	defer cancel()

	var averages []TankAverages
	cur, err := database.DefaultClient.Collection(database.CollectionVehicleAverages).Find(ctx, bson.M{"_id": bson.M{"$in": vehicleIDs}})
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
