package replay

import (
	"github.com/cufee/aftermath-core/dataprep/replay"
	"github.com/cufee/aftermath-core/internal/logic/render"
)

func statsBlockToBlock(stats replay.StatsBlock, width float64) render.Block {
	return render.NewBlocksContent(render.Style{Direction: render.DirectionVertical, AlignItems: render.AlignItemsCenter, Width: width},
		render.NewTextContent(render.Style{
			Font:      &render.FontLarge,
			FontColor: render.TextPrimary,
		}, stats.Value.String),
		render.NewTextContent(render.Style{
			Font:      &render.FontSmall,
			FontColor: render.TextAlt,
		}, stats.Label))
}
