package period

import (
	"github.com/cufee/aftermath-core/dataprep"
	"github.com/cufee/aftermath-core/dataprep/period"
	"github.com/cufee/aftermath-core/internal/logic/render"
	"github.com/cufee/aftermath-core/internal/logic/render/shared"
)

func statsBlocksToColumnBlock(style overviewStyle, statsBlocks []period.StatsBlock) (render.Block, error) {
	var content []render.Block

	for _, statsBlock := range statsBlocks {
		if statsBlock.Flavor == period.BlockFlavorSpecial {
			content = append(content, uniqueStatsBlock(style, statsBlock))
		} else {
			content = append(content, defaultStatsBlock(style, statsBlock))
		}
	}
	return render.NewBlocksContent(style.container, content...), nil
}

func uniqueStatsBlock(style overviewStyle, stats period.StatsBlock) render.Block {
	switch stats.Tag {
	case dataprep.TagWN8:
		return uniqueBlockWN8(style, stats)
	default:
		return defaultStatsBlock(style, stats)
	}
}

func defaultStatsBlock(style overviewStyle, stats period.StatsBlock) render.Block {
	valueStyle, labelStyle := style.block(stats)

	blocks := []render.Block{render.NewTextContent(valueStyle, stats.Data.String)}
	blocks = append(blocks, render.NewTextContent(labelStyle, stats.Label))

	return render.NewBlocksContent(style.blockContainer, blocks...)
}

func uniqueBlockWN8(style overviewStyle, stats period.StatsBlock) render.Block {
	var blocks []render.Block

	valueStyle, labelStyle := style.block(stats)
	valueBlock := render.NewTextContent(valueStyle, stats.Data.String)

	ratingColors := shared.GetWN8Colors(int(stats.Data.Value))
	if stats.Data.Value <= 0 {
		ratingColors.Content = render.TextAlt
		ratingColors.Background = render.TextAlt
	}

	iconTop := shared.AftermathLogo(ratingColors.Background, shared.DefaultLogoOptions())
	iconBlockTop := render.NewImageContent(render.Style{Width: float64(iconTop.Bounds().Dx()), Height: float64(iconTop.Bounds().Dy())}, iconTop)

	style.blockContainer.Gap = 10
	blocks = append(blocks, render.NewBlocksContent(style.blockContainer, iconBlockTop, valueBlock))

	if stats.Data.Value >= 0 {
		labelStyle.FontColor = ratingColors.Content
		blocks = append(blocks, render.NewBlocksContent(render.Style{
			PaddingY:        5,
			PaddingX:        10,
			BorderRadius:    15,
			BackgroundColor: ratingColors.Background,
		}, render.NewTextContent(labelStyle, shared.GetWN8TierName(int(stats.Data.Value)))))
	}

	return render.NewBlocksContent(render.Style{Direction: render.DirectionVertical, AlignItems: render.AlignItemsCenter, Gap: 10, PaddingY: 5}, blocks...)
}
