package dataprep

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
	Value  float64 `json:"value"`
	String string  `json:"string"`
}
