package stats

import "github.com/cufee/aftermath-core/internal/core/utils"

func EmptySession(accountID, lastBattle int) SessionSnapshot {
	return SessionSnapshot{
		AccountID:      accountID,
		LastBattleTime: lastBattle,
		Global:         ReducedStatsFrame{},
		Rating:         ReducedStatsFrame{},
		Vehicles:       make(map[int]ReducedVehicleStats),
	}
}

type SessionSnapshot struct {
	AccountID      int `json:"accountId" bson:"accountId"`
	LastBattleTime int `json:"lastBattleTime" bson:"lastBattleTime"`

	Global ReducedStatsFrame `json:"global" bson:"global"`
	Rating ReducedStatsFrame `json:"rating" bson:"rating"`

	Vehicles map[int]ReducedVehicleStats `json:"vehicles" bson:"vehicles"`
}

func (a *SessionSnapshot) Diff(b SessionSnapshot) (SessionSnapshot, error) {
	var diff SessionSnapshot
	if err := utils.DeepCopy[SessionSnapshot](a, &diff); err != nil {
		return SessionSnapshot{}, err
	}

	diff.Global.Subtract(b.Global)
	diff.Rating.Subtract(b.Rating)

	for vehicleID, bVehicle := range b.Vehicles {
		vehicleStats, ok := diff.Vehicles[vehicleID]
		if !ok {
			diff.Vehicles[vehicleID] = bVehicle
		} else {
			vehicleStats.Subtract(bVehicle)
		}
	}

	if b.LastBattleTime > diff.LastBattleTime {
		diff.LastBattleTime = b.LastBattleTime
	}

	return diff, nil
}
