package stats

import "github.com/cufee/aftermath-core/internal/core/utils"

func EmptySession(accountID, lastBattle int) *SessionSnapshot {
	return &SessionSnapshot{
		AccountID:      accountID,
		LastBattleTime: lastBattle,
		Global:         &ReducedStatsFrame{},
		Rating:         &ReducedStatsFrame{},
		Vehicles:       make(map[int]*ReducedVehicleStats),
	}
}

type SessionSnapshot struct {
	AccountID      int `json:"accountId" bson:"accountId"`
	LastBattleTime int `json:"lastBattleTime" bson:"lastBattleTime"`

	Global *ReducedStatsFrame `json:"global" bson:"global"`
	Rating *ReducedStatsFrame `json:"rating" bson:"rating"`

	Vehicles map[int]*ReducedVehicleStats `json:"vehicles" bson:"vehicles"`
}

func (s *SessionSnapshot) Diff(other *SessionSnapshot) (*SessionSnapshot, error) {
	var diff SessionSnapshot
	if err := utils.DeepCopy[SessionSnapshot](s, &diff); err != nil {
		return nil, err
	}

	diff.Global.Subtract(other.Global)
	diff.Rating.Subtract(other.Rating)

	for vehicleID, otherVehicleStats := range other.Vehicles {
		vehicleStats, ok := diff.Vehicles[vehicleID]
		if !ok {
			diff.Vehicles[vehicleID] = otherVehicleStats
		} else {
			vehicleStats.Subtract(otherVehicleStats)
			if vehicleStats.Battles == 0 {
				delete(diff.Vehicles, vehicleID)
			}
		}
	}

	if other.LastBattleTime > diff.LastBattleTime {
		diff.LastBattleTime = other.LastBattleTime
	}

	return &diff, nil
}
