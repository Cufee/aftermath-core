package period

import (
	"github.com/cufee/aftermath-core/dataprep/period"
	"github.com/cufee/aftermath-core/internal/logic/render"
)

func statsBlocksToRowBlock(blocks []period.StatsBlock) (render.Block, error) {
	var content []render.Block
	for _, statsBlock := range blocks {
		blocks := make([]render.Block, 0, 2)
		blocks = append(blocks, render.NewTextContent(render.Style{Font: &render.FontLarge, FontColor: render.TextPrimary}, statsBlock.Data.String))
		blocks = append(blocks, render.NewTextContent(render.Style{Font: &render.FontLarge, FontColor: render.TextSecondary}, statsBlock.Label))

		content = append(content, render.NewBlocksContent(render.Style{AlignItems: render.AlignItemsCenter, Direction: render.DirectionVertical}, blocks...))
	}
	return render.NewBlocksContent(render.Style{}, content...), nil
}
