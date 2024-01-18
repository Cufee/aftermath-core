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

type value struct {
	Value  any    `json:"value"`
	String string `json:"string"`
}

type StatsBlock struct {
	Session value  `json:"session"`
	Career  value  `json:"career"`
	Label   string `json:"label"`
}
