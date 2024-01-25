package session

import (
	"image/color"

	"github.com/cufee/aftermath-core/dataprep"
	"github.com/cufee/aftermath-core/internal/logic/render"
	"github.com/rs/zerolog/log"
)

type convertOptions struct {
	showSessionStats bool
	showCareerStats  bool
	showLabels       bool
}

type styledStatsBlock struct {
	dataprep.StatsBlock
	style render.Style
}

func statsBlocksToCardBlocks(stats []styledStatsBlock, opts ...convertOptions) ([]render.Block, error) {
	var options convertOptions = convertOptions{
		showSessionStats: true,
		showCareerStats:  true,
		showLabels:       true,
	}
	if len(opts) > 0 {
		options = opts[0]
	}

	var content []render.Block
	for _, statsBlock := range stats {
		blocks := make([]render.Block, 0, 3)
		if options.showSessionStats {
			blocks = append(blocks, newStatsBlockRow(render.Style{Font: &FontLarge, FontColor: FontLargeColor}, statsBlock.Session, comparisonIconFromBlock(statsBlock.StatsBlock)))
		}
		if options.showCareerStats && statsBlock.Career.String != "" {
			blocks = append(blocks, newStatsBlockRow(render.Style{Font: &FontMedium, FontColor: FontMediumColor}, statsBlock.Career, blockToWN8Icon(statsBlock.Career, statsBlock.Tag)))
		}
		if options.showLabels && statsBlock.Tag != dataprep.TagBattles {
			blocks = append(blocks, render.NewTextContent(render.Style{Font: &FontSmall, FontColor: FontSmallColor}, statsBlock.Label))
		}
		content = append(content, render.NewBlocksContent(statsBlock.style, blocks...))
	}
	return content, nil
}

func newStatsBlockRow(style render.Style, stats dataprep.Value, icon *comparisonIcon) render.Block {
	if icon == nil {
		return render.NewTextContent(style, stats.String)
	}

	return render.NewBlocksContent(
		render.Style{Direction: render.DirectionHorizontal, AlignItems: render.AlignItemsCenter},
		icon.left,
		render.NewTextContent(style, stats.String),
		icon.right,
	)
}

func newPlayerTitleCard(style render.Style, name string, clanTagBlocks []render.Block) render.Block {
	if len(clanTagBlocks) == 0 {
		return render.NewBlocksContent(style, render.NewTextContent(render.Style{Font: &FontLarge, FontColor: FontLargeColor}, name))
	}

	content := make([]render.Block, 0, 3)
	style.JustifyContent = render.JustifyContentSpaceBetween

	// Visible tag
	clanTagBlock := render.NewBlocksContent(render.Style{
		Direction:       render.DirectionHorizontal,
		AlignItems:      render.AlignItemsCenter,
		PaddingX:        10,
		PaddingY:        2.5,
		Gap:             2.5,
		BackgroundColor: HighlightCardColor(style.BackgroundColor),
		BorderRadius:    10,
		// Debug:           true,
	}, clanTagBlocks...)

	clanTagImage, err := clanTagBlock.Render()
	if err != nil {
		log.Warn().Err(err).Msg("failed to render clan tag")
		// This error is not fatal, we can just render the name
		return render.NewBlocksContent(style, render.NewTextContent(render.Style{Font: &FontLarge, FontColor: FontLargeColor}, name))
	}
	content = append(content, render.NewImageContent(render.Style{Width: float64(clanTagImage.Bounds().Dx()), Height: float64(clanTagImage.Bounds().Dy())}, clanTagImage))

	// Nickname
	content = append(content, render.NewTextContent(render.Style{Font: &FontLarge, FontColor: FontLargeColor}, name))

	// Invisible tag to offset the nickname
	clanBlock := render.NewBlocksContent(render.Style{
		Width:          float64(clanTagImage.Bounds().Dx()),
		JustifyContent: render.JustifyContentEnd,
	}, render.NewTextContent(render.Style{Font: &FontLarge, FontColor: color.Transparent}, "-"))

	content = append(content, clanBlock)

	return render.NewBlocksContent(style, render.NewBlocksContent(render.Style{
		JustifyContent: render.JustifyContentSpaceBetween,
		Direction:      render.DirectionHorizontal,
		AlignItems:     render.AlignItemsCenter,
		Width:          BaseCardWidth,
		PaddingX:       20,
		// Debug:          true,
	}, content...))
}

func newTextLabel(label string) render.Block {
	return render.NewTextContent(render.Style{Font: &FontMedium, FontColor: FontMediumColor}, label)
}

func newCardBlock(cardStyle render.Style, label render.Block, stats []render.Block) render.Block {
	var content []render.Block
	content = append(content, label)
	content = append(content, render.NewBlocksContent(render.Style{
		Direction:      render.DirectionHorizontal,
		AlignItems:     render.AlignItemsCenter,
		JustifyContent: render.JustifyContentSpaceBetween,
		Gap:            10,
		// Debug:          true,
	}, stats...))

	return render.NewBlocksContent(cardStyle, content...)
}
