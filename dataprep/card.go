package dataprep

import "github.com/cufee/aftermath-core/internal/core/stats"

type cardType string

const (
	CardTypeVehicle  cardType = "vehicle"
	CardTypeOverview cardType = "overview"
)

type StatsCard struct {
	Type   cardType     `json:"type"`
	Title  string       `json:"title"`
	Blocks []StatsBlock `json:"blocks"`
}

type SessionCards []StatsCard

type StatsBlock struct {
	Session Value  `json:"session"`
	Career  Value  `json:"career"`
	Label   string `json:"label"`
	Tag     Tag    `json:"tag"`
}

type Tag string

const (
	TagWN8         Tag = "wn8"
	TagBattles     Tag = "battles"
	TagWinrate     Tag = "winrate"
	TagAccuracy    Tag = "accuracy"
	TagAvgDamage   Tag = "avg_damage"
	TagDamageRatio Tag = "damage_ratio"
)

type Value struct {
	Value  any    `json:"value"`
	String string `json:"string"`
}

func (v *Value) Compare(other Value) int {
	if !stats.ValueValid(v.Value) || !stats.ValueValid(other.Value) {
		return 0
	}
	if v.Float64() > other.Float64() {
		return 1
	}
	if v.Float64() < other.Float64() {
		return -1
	}
	return 0
}

func (v Value) Float64() float64 {
	switch v.Value.(type) {
	case int:
		return float64(v.Value.(int))
	case float64:
		return v.Value.(float64)
	case float32:
		return float64(v.Value.(float32))
	}
	return 0
}
