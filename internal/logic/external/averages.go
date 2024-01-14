package external

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cufee/aftermath-core/internal/core/utils"
)

// Response from https://www.blitzstars.com/ average stats endpoint
type VehicleAverages struct {
	TankID int `json:"tank_id"`
	All    struct {
		Battles              float64 `json:"battles,omitempty"`
		DroppedCapturePoints float64 `json:"dropped_capture_points,omitempty"`
	} `json:",omitempty"`
	Special struct {
		Winrate         float64 `json:"winrate,omitempty"`
		DamageRatio     float64 `json:"damageRatio,omitempty"`
		Kdr             float64 `json:"kdr,omitempty"`
		DamagePerBattle float64 `json:"damagePerBattle,omitempty"`
		KillsPerBattle  float64 `json:"killsPerBattle,omitempty"`
		HitsPerBattle   float64 `json:"hitsPerBattle,omitempty"`
		SpotsPerBattle  float64 `json:"spotsPerBattle,omitempty"`
		Wpm             float64 `json:"wpm,omitempty"`
		Dpm             float64 `json:"dpm,omitempty"`
		Kpm             float64 `json:"kpm,omitempty"`
		HitRate         float64 `json:"hitRate,omitempty"`
		SurvivalRate    float64 `json:"survivalRate,omitempty"`
	} `json:"special,omitempty"`
}

var starsStatsAveragesURL = utils.MustGetEnv("BLITZ_STARS_AVERAGES_URL")

func GetTankAverages() (data []VehicleAverages, err error) {
	res, err := insecureClient.Get(starsStatsAveragesURL)
	if err != nil || res == nil || res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code: %+v. error: %s", res, err)
	}
	defer res.Body.Close()
	return nil, json.NewDecoder(res.Body).Decode(&data)
}
