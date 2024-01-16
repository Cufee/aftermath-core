package cache

import (
	"fmt"
	"strings"

	"github.com/cufee/aftermath-core/internal/core/database"
	"github.com/cufee/aftermath-core/internal/core/localization"
	"github.com/cufee/aftermath-core/internal/core/wargaming"
	"github.com/cufee/aftermath-core/internal/logic/external"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type vehicleType string
type vehicleClass string

const (
	VehicleTypeUnknown     vehicleType = "unknown"
	VehicleTypeRegular     vehicleType = "regular"
	VehicleTypePremium     vehicleType = "premium"
	VehicleTypeCollectible vehicleType = "collectible"
	VehicleTypeEarlyAccess vehicleType = "earlyAccess"

	VehicleClassUnknown       vehicleClass = "unknown"
	VehicleClassLightTank     vehicleClass = "lightTank"
	VehicleClassHeavyTank     vehicleClass = "heavyTank"
	VehicleClassMediumTank    vehicleClass = "mediumTank"
	VehicleClassTankDestroyer vehicleClass = "tankDestroyer"
	VehicleClassArtillery     vehicleClass = "artillery"
)

type VehicleInfo struct {
	ID int `json:"id" bson:"_id"`

	Tier           int               `json:"tier" bson:"tier"`
	Nation         string            `json:"nation" bson:"nation"`
	LocalizedNames map[string]string `json:"localized_names" bson:"localized_names"`

	Class vehicleClass `json:"class" bson:"class"`
	Type  vehicleType  `json:"type" bson:"type"`
}

func (v VehicleInfo) IsPremium() bool {
	return v.Type == VehicleTypePremium
}

func (v VehicleInfo) Name(lang localization.SupportedLanguage) string {
	if name, ok := v.LocalizedNames[lang.WargamingCode]; ok {
		return name
	}
	if name, ok := v.LocalizedNames[localization.LanguageEN.WargamingCode]; ok {
		return name
	}
	return fmt.Sprintf("Secret Tank #%d", v.ID)
}

type AchievementInfo struct {
	ID             string                                    `json:"id" bson:"_id"`
	ImageURL       string                                    `json:"image" bson:"image"`
	Description    string                                    `json:"description" bson:"description"`
	LocalizedNames map[localization.SupportedLanguage]string `json:"localized_names" bson:"localized_names"`
}

func (a AchievementInfo) Name(lang localization.SupportedLanguage) string {
	if name, ok := a.LocalizedNames[lang]; ok {
		return name
	}
	if name, ok := a.LocalizedNames[localization.LanguageEN]; ok {
		return name
	}
	return fmt.Sprintf("Secret Achievement #%s", a.ID)
}

func GetCompleteVehicleGlossary() (map[int]VehicleInfo, error) {
	vehicles := make(map[int]VehicleInfo)
	glossaryLocales := []localization.SupportedLanguage{localization.LanguageEN, localization.LanguageRU}
	for _, locale := range glossaryLocales {
		glossary, err := wargaming.Clients.Cache.GetVehiclesGlossary(locale.WargamingCode)
		if err != nil {
			return nil, err
		}

		for _, vehicle := range glossary {
			existingData, ok := vehicles[vehicle.TankID]
			if !ok {
				existingData = VehicleInfo{
					ID:   vehicle.TankID,
					Tier: vehicle.Tier,
					LocalizedNames: map[string]string{
						locale.WargamingCode: vehicle.Name,
					},
					Class: VehicleClassUnknown,
					Type:  VehicleTypeRegular,
				}
				// TODO: Detect classes and collectible vehicles
				if vehicle.IsPremium {
					existingData.Type = VehicleTypePremium
				}
			}

			if strings.HasPrefix(vehicle.Name, "#") {
				// TODO: Handle secret vehicles
				continue
			}
			existingData.LocalizedNames[locale.WargamingCode] = vehicle.Name
			existingData.Nation = vehicle.Nation
			existingData.Tier = vehicle.Tier

			vehicles[vehicle.TankID] = existingData
		}
	}

	inspectorData, err := external.GetInspectorVehicles()
	if err != nil {
		return nil, err
	}

	for id, vehicle := range inspectorData {
		if _, ok := vehicles[id]; ok {
			continue
		}
		names := make(map[string]string)

		if !strings.HasPrefix(vehicle.NameEn, "#") {
			names[localization.LanguageEN.WargamingCode] = vehicle.NameEn
		}
		if !strings.HasPrefix(vehicle.NameRu, "#") {
			names[localization.LanguageRU.WargamingCode] = vehicle.NameRu
		}
		vehicles[id] = VehicleInfo{
			ID:             id,
			Tier:           vehicle.Tier,
			Class:          VehicleClassUnknown,
			Type:           VehicleTypeEarlyAccess,
			LocalizedNames: names,
		}
	}

	return vehicles, nil
}

func UpdateGlossaryCache() error {
	vehicles, err := GetCompleteVehicleGlossary()
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

// TODO: add in-memory cache

func GetGlossaryVehicles(vehicleIDs ...int) (map[int]VehicleInfo, error) {
	if len(vehicleIDs) == 0 {
		return nil, nil
	}

	ctx, cancel := database.DefaultClient.Ctx()
	defer cancel()

	var vehicles []VehicleInfo
	cur, err := database.DefaultClient.Collection(database.CollectionVehicleGlossary).Find(ctx, bson.M{"_id": bson.M{"$in": vehicleIDs}})
	if err != nil {
		return nil, err
	}

	err = cur.All(ctx, &vehicles)
	if err != nil {
		return nil, err
	}

	vehicleMap := make(map[int]VehicleInfo)
	for _, vehicle := range vehicles {
		vehicleMap[vehicle.ID] = vehicle
	}

	return vehicleMap, nil
}
