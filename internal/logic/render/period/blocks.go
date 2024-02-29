package period

import (
	"github.com/cufee/aftermath-core/dataprep/period"
	"github.com/cufee/aftermath-core/internal/logic/render"
)

func statsBlocksToRowBlock(style overviewStyle, statsBlocks []period.StatsBlock) (render.Block, error) {
	var content []render.Block

	for _, statsBlock := range statsBlocks {
		valueStyle, labelStyle := style.block(statsBlock.Flavor)

		blocks := []render.Block{render.NewTextContent(valueStyle, statsBlock.Data.String)}
		blocks = append(blocks, render.NewTextContent(labelStyle, statsBlock.Label))

		content = append(content, render.NewBlocksContent(render.Style{
			Width:      style.container.Width / float64(len(statsBlocks)),
			Direction:  render.DirectionVertical,
			AlignItems: render.AlignItemsCenter,
		},
			blocks...,
		))
	}
	return render.NewBlocksContent(style.container, content...), nil
}
