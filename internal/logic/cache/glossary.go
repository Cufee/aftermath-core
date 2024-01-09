package cache

import (
	"fmt"

	"github.com/cufee/aftermath-core/internal/core/localization"
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

	Tier           int                                       `json:"tier" bson:"tier"`
	Nation         string                                    `json:"nation" bson:"nation"`
	LocalizedNames map[localization.SupportedLanguage]string `json:"localized_names" bson:"localized_names"`

	Class vehicleClass `json:"class" bson:"class"`
	Type  vehicleType  `json:"type" bson:"type"`
}

func (v VehicleInfo) IsPremium() bool {
	return v.Type == VehicleTypePremium
}

func (v VehicleInfo) Name(lang localization.SupportedLanguage) string {
	if name, ok := v.LocalizedNames[lang]; ok {
		return name
	}
	if name, ok := v.LocalizedNames[localization.LanguageEN]; ok {
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
