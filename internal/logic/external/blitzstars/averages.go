package blitzstars

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cufee/aftermath-core/internal/core/stats"
)

// Response from https://www.blitzstars.com/ average stats endpoint
type VehicleAverages struct {
	TankID  int `json:"tank_id"`
	Players int `json:"number_of_players"`
	All     struct {
		AvgBattles              float64 `json:"battles,omitempty"`
		AvgDroppedCapturePoints float64 `json:"dropped_capture_points,omitempty"`
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

func GetTankAverages() (map[int]stats.ReducedStatsFrame, error) {
	res, err := insecureClient.Get(starsStatsApiURL + "/tankaverages.json")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status code: %d", res.StatusCode)
	}

	var averages []VehicleAverages
	err = json.NewDecoder(res.Body).Decode(&averages)
	if err != nil {
		return nil, err
	}

	averagesMap := make(map[int]stats.ReducedStatsFrame)
	for _, average := range averages {
		battles := average.All.AvgBattles * float64(average.Players)

		averagesMap[average.TankID] = stats.ReducedStatsFrame{
			Battles:     int(battles),
			BattlesWon:  int(average.Special.Winrate * battles / 100),
			DamageDealt: int(average.Special.DamagePerBattle * battles),

			ShotsHit:   int(average.Special.HitsPerBattle * battles),
			ShotsFired: int((average.Special.HitsPerBattle * battles) / (average.Special.HitRate / 100)),

			Frags:                int(average.Special.KillsPerBattle * battles),
			EnemiesSpotted:       int(average.Special.SpotsPerBattle * battles),
			DroppedCapturePoints: int(average.All.AvgDroppedCapturePoints * battles),
		}
	}

	return averagesMap, nil
}
