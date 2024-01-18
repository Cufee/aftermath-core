package cache

import (
	"github.com/cufee/aftermath-core/internal/core/database"
	"github.com/cufee/aftermath-core/internal/logic/external"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func UpdateGlossaryCache() error {
	vehicles, err := external.GetCompleteVehicleGlossary()
	if err != nil {
		return err
	}

	var vehicleWrites []mongo.WriteModel
	for _, vehicle := range vehicles {
		model := mongo.NewUpdateOneModel()
		model.SetFilter(bson.M{"_id": vehicle.ID})
		model.SetUpdate(bson.M{"$set": vehicle})
		model.SetUpsert(true)
		vehicleWrites = append(vehicleWrites, model)
	}

	ctx, cancel := database.DefaultClient.Ctx()
	defer cancel()

	_, err = database.DefaultClient.Collection(database.CollectionVehicleGlossary).BulkWrite(ctx, vehicleWrites)
	if err != nil {
		return err
	}

	return nil
}
