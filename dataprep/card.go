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
	switch v.Value.(type) {
	case int:
		if v.Value.(int) > other.Value.(int) {
			return 1
		} else if v.Value.(int) < other.Value.(int) {
			return -1
		}
	case float64:
		if v.Value.(float64) > other.Value.(float64) {
			return 1
		} else if v.Value.(float64) < other.Value.(float64) {
			return -1
		}
	case float32:
		if v.Value.(float32) > other.Value.(float32) {
			return 1
		}
		if v.Value.(float32) < other.Value.(float32) {
			return -1
		}
	}
	return 0
}
