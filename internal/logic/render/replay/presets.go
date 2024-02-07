package replay

import (
	"github.com/cufee/aftermath-core/dataprep"
	"github.com/cufee/aftermath-core/dataprep/replay"
	"github.com/cufee/aftermath-core/internal/logic/render"
)

var blockWidthPresets = map[dataprep.Tag]float64{
	dataprep.TagWN8:                    75,
	dataprep.TagDamageDealt:            75,
	dataprep.TagDamageBlocked:          75,
	dataprep.TagDamageAssisted:         75,
	dataprep.TagDamageAssistedCombined: 100,
	dataprep.TagFrags:                  30,
}

func statsBlockToBlock(stats replay.StatsBlock) render.Block {
	width, ok := blockWidthPresets[stats.Tag]
	if !ok {
		width = 75
	}
	return render.NewBlocksContent(render.Style{Direction: render.DirectionVertical, Width: width, AlignItems: render.AlignItemsCenter},
		render.NewTextContent(render.Style{
			Font:      &render.FontLarge,
			FontColor: render.TextPrimary,
		}, stats.Value.String),
		render.NewTextContent(render.Style{
			Font:      &render.FontSmall,
			FontColor: render.TextAlt,
		}, stats.Label))
}
