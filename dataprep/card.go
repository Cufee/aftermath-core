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

type Value struct {
	Value  float64 `json:"value"`
	String string  `json:"string"`
}
