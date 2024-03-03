package replay

import (
	"github.com/cufee/aftermath-core/dataprep"
	"github.com/cufee/aftermath-core/internal/logic/replay"
)

type Cards struct {
	Allies  []Card `json:"allies"`
	Enemies []Card `json:"enemies"`
}

type Card dataprep.StatsCard[StatsBlock, CardMeta]

type CardMeta struct {
	Player replay.Player  `json:"player"`
	Tags   []dataprep.Tag `json:"tags"`
}

type StatsBlock struct {
	Label string `json:"label"`

	Tag   dataprep.Tag   `json:"tag"`
	Value dataprep.Value `json:"value"`
}
