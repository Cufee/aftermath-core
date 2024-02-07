package dataprep

type cardType string

const (
	CardTypeVehicle  cardType = "vehicle"
	CardTypeOverview cardType = "overview"
)

type StatsCard[T, M interface{}] struct {
	Type   cardType `json:"type"`
	Title  string   `json:"title"`
	Blocks []T      `json:"blocks"`
	Meta   M        `json:"meta,omitempty"`
}

type Tag string

const (
	// Global
	TagWN8      Tag = "wn8"
	TagFrags    Tag = "frags"
	TagBattles  Tag = "battles"
	TagWinrate  Tag = "winrate"
	TagAccuracy Tag = "accuracy"

	// Session Specific
	TagAvgDamage   Tag = "avg_damage"
	TagDamageRatio Tag = "damage_ratio"

	// Replay Specific
	TagDamageDealt            Tag = "damage_dealt"
	TagDamageTaken            Tag = "damage_taken"
	TagDamageBlocked          Tag = "spotted"
	TagDamageAssisted         Tag = "assisted"
	TagDamageAssistedCombined Tag = "assisted_combined"
)

type Value struct {
	Value  float64 `json:"value"`
	String string  `json:"string"`
}
