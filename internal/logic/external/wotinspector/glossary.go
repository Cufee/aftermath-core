package wotinspector

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/cufee/aftermath-core/internal/core/database/models"
	"github.com/cufee/aftermath-core/internal/core/utils"
	"github.com/cufee/aftermath-core/internal/core/wargaming"
	"golang.org/x/text/language"
)

// Response from https://wotinspector.com/en/
type InspectorVehicle struct {
	NameEn  string `json:"en"`
	NameRu  string `json:"ru"`
	Tier    int    `json:"tier"`
	Type    int    `json:"type"`
	Premium int    `json:"premium"`
}

var inspectorVehiclesURL string

func init() {
	inspectorVehiclesURL = utils.MustGetEnv("WOT_INSPECTOR_TANK_DB_URL")
}

func GetInspectorVehicles() (map[int]InspectorVehicle, error) {
	re := regexp.MustCompile(`(\d{1,9}):`)
	tanks := make(map[int]InspectorVehicle)

	res, err := insecureClient.Get(inspectorVehiclesURL)
	if err != nil || res == nil || res.StatusCode != http.StatusOK {
		return tanks, fmt.Errorf("status code: %+v. error: %s", res, err)
	}
	defer res.Body.Close()
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return tanks, fmt.Errorf("status code: %+v. error: %s", res, err)
	}
	tanksString := strings.ReplaceAll(string(bodyBytes), "TANK_DB = ", "")
	tanksString = re.ReplaceAllString(tanksString, `"$1":`)
	split := strings.SplitAfter(tanksString, "},")
	if len(split) <= 2 {
		return tanks, fmt.Errorf("failed to split string")
	}
	fix := strings.ReplaceAll(split[len(split)-2], "},", "}")
	tanksString = strings.ReplaceAll(tanksString, split[len(split)-2], fix)
	return tanks, json.Unmarshal([]byte(tanksString), &tanks)
}

func GetCompleteVehicleGlossary() (map[int]models.Vehicle, error) {
	vehicles := make(map[int]models.Vehicle)
	glossaryLocales := []language.Tag{language.English, language.Russian, language.Polish}
	for _, locale := range glossaryLocales {
		glossary, err := wargaming.Clients.Cache.CompleteVehicleGlossary("EU", locale.String())
		if err != nil {
			return nil, err
		}

		for _, vehicle := range glossary {
			existingData, ok := vehicles[vehicle.TankID]
			if !ok {
				existingData = models.Vehicle{
					ID:   vehicle.TankID,
					Tier: vehicle.Tier,
					LocalizedNames: map[string]string{
						locale.String(): vehicle.Name,
					},
					Class: models.VehicleClassUnknown,
					Type:  models.VehicleTypeRegular,
				}
				// TODO: Detect classes and collectible vehicles
				if vehicle.IsPremium {
					existingData.Type = models.VehicleTypePremium
				}
			}

			if strings.HasPrefix(vehicle.Name, "#") {
				// TODO: Handle secret vehicles
				continue
			}
			existingData.LocalizedNames[locale.String()] = vehicle.Name
			existingData.Nation = vehicle.Nation
			existingData.Tier = vehicle.Tier

			vehicles[vehicle.TankID] = existingData
		}
	}

	inspectorData, err := GetInspectorVehicles()
	if err != nil {
		return nil, err
	}

	for id, vehicle := range inspectorData {
		if _, ok := vehicles[id]; ok {
			continue
		}
		names := make(map[string]string)

		if !strings.HasPrefix(vehicle.NameEn, "#") {
			names[language.English.String()] = vehicle.NameEn
		}
		if !strings.HasPrefix(vehicle.NameRu, "#") {
			names[language.Russian.String()] = vehicle.NameRu
		}
		vehicles[id] = models.Vehicle{
			ID:             id,
			Tier:           vehicle.Tier,
			Class:          models.VehicleClassUnknown,
			Type:           models.VehicleTypeEarlyAccess,
			LocalizedNames: names,
		}
	}

	return vehicles, nil
}
