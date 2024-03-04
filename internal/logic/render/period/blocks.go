package period

import (
	"image"
	"image/color"

	"github.com/cufee/aftermath-core/dataprep"
	"github.com/cufee/aftermath-core/dataprep/period"
	"github.com/cufee/aftermath-core/internal/logic/render"
	"github.com/cufee/aftermath-core/internal/logic/render/shared"
	"github.com/fogleman/gg"
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
	valueStyle, _ := style.block(stats)
	valueBlock := render.NewTextContent(valueStyle, stats.Data.String)

	var ratingColor color.Color = render.TextAlt
	if stats.Data.Value > 0 {
		ratingColor = shared.GetWN8Color(int(stats.Data.Value))
	}
	iconTop := getIconWN8(ratingColor, defaultRatingIconOptions(1))
	iconBottom := getIconWN8(ratingColor, defaultRatingIconOptions(0))
	iconBlockTop := render.NewImageContent(render.Style{Width: float64(iconTop.Bounds().Dx()), Height: float64(iconTop.Bounds().Dy())}, iconTop)
	iconBlockBottom := render.NewImageContent(render.Style{Width: float64(iconBottom.Bounds().Dx()), Height: float64(iconBottom.Bounds().Dy())}, iconBottom)

	return render.NewBlocksContent(style.blockContainer, iconBlockTop, valueBlock, iconBlockBottom)
}

func getIconWN8(ratingColor color.Color, opts ratingIconOptions) image.Image {
	ctx := gg.NewContext(opts.width(), opts.height())
	for line := range opts.lines {
		height := opts.lineStep + opts.lineStep*line
		if line > opts.lines/2 {
			height = opts.lineStep * (opts.lines - line)
		}

		offset := 0.0
		if opts.direction == 1 {
			offset = float64(opts.height() - height)
		}

		ctx.DrawRoundedRectangle(float64(line*(int(opts.lineWidth+opts.gap))), offset, opts.lineWidth, float64(height), 3)
	}

	ctx.SetColor(ratingColor)
	ctx.Fill()
	return ctx.Image()
}
