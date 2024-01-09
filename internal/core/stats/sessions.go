package stats

type SessionSnapshot struct {
	AccountID      int `json:"accountId" bson:"accountId"`
	LastBattleTime int `json:"lastBattleTime" bson:"lastBattleTime"`

	Global ReducedStatsFrame `json:"global" bson:"global"`
	Rating ReducedStatsFrame `json:"rating" bson:"rating"`

	Vehicles map[int]ReducedVehicleStats `json:"vehicles" bson:"vehicles"`
}

func (s *SessionSnapshot) Add(other SessionSnapshot) {
	s.Global.Add(other.Global)
	s.Rating.Add(other.Rating)

	for vehicleID, otherVehicleStats := range other.Vehicles {
		vehicleStats, ok := s.Vehicles[vehicleID]
		if !ok {
			s.Vehicles[vehicleID] = otherVehicleStats
		} else {
			vehicleStats.Add(otherVehicleStats)
		}
	}

	if other.LastBattleTime > s.LastBattleTime {
		s.LastBattleTime = other.LastBattleTime
	}
}

func (s *SessionSnapshot) Subtract(other SessionSnapshot) {
	s.Global.Subtract(other.Global)
	s.Rating.Subtract(other.Rating)

	for vehicleID, otherVehicleStats := range other.Vehicles {
		vehicleStats, ok := s.Vehicles[vehicleID]
		if !ok {
			s.Vehicles[vehicleID] = otherVehicleStats
		} else {
			vehicleStats.Subtract(otherVehicleStats)
		}
	}

	if other.LastBattleTime > s.LastBattleTime {
		s.LastBattleTime = other.LastBattleTime
	}
}
