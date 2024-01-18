package session

import (
	"github.com/cufee/aftermath-core/internal/logic/dataprep"
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
	for i, statsBlock := range stats {
		blocks := make([]render.Block, 0, 3)
		if options.showSessionStats {
			blocks = append(blocks, render.NewTextContent(render.Style{Font: &FontLarge, FontColor: FontLargeColor}, statsBlock.Session.String))
		}
		if options.showCareerStats && statsBlock.Career.String != "" {
			blocks = append(blocks, render.NewTextContent(render.Style{Font: &FontMedium, FontColor: FontMediumColor}, statsBlock.Career.String))
		}
		if options.showLabels && i != 0 && statsBlock.Label != "" {
			blocks = append(blocks, render.NewTextContent(render.Style{Font: &FontSmall, FontColor: FontSmallColor}, statsBlock.Label))
		}
		content = append(content, render.NewBlocksContent(statsBlock.style, blocks...))
	}
	return content, nil
}

func newPlayerTitleCard(style render.Style, name, clanTag string, clanSubHeader render.Block) render.Block {
	if clanTag == "" {
		return render.NewBlocksContent(style, render.NewTextContent(render.Style{Font: &FontLarge, FontColor: FontLargeColor}, name))
	}

	content := make([]render.Block, 0, 3)
	style.JustifyContent = render.JustifyContentSpaceBetween

	// Visible tag
	clanTagBlock := render.NewBlocksContent(render.Style{
		Direction:       render.DirectionHorizontal,
		AlignItems:      render.AlignItemsCenter,
		PaddingX:        10,
		PaddingY:        5,
		BackgroundColor: HighlightCardColor(style.BackgroundColor),
		BorderRadius:    10,
		// Debug:           true,
	}, render.NewTextContent(render.Style{Font: &FontMedium, FontColor: FontMediumColor}, clanTag))

	clanTagImage, err := clanTagBlock.Render()
	if err != nil {
		log.Warn().Err(err).Msg("failed to render clan tag")
		// This error is not fatal, we can just render the name
		return render.NewBlocksContent(style, render.NewTextContent(render.Style{Font: &FontLarge, FontColor: FontLargeColor}, name))
	}
	content = append(content, render.NewImageContent(render.Style{Width: float64(clanTagImage.Bounds().Dx()), Height: float64(clanTagImage.Bounds().Dy())}, clanTagImage))

	// Nickname
	content = append(content, render.NewTextContent(render.Style{Font: &FontLarge, FontColor: FontLargeColor}, name))

	clanBlock := render.NewBlocksContent(render.Style{
		Width:          float64(clanTagImage.Bounds().Dx()),
		JustifyContent: render.JustifyContentEnd,
	}, clanSubHeader)

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

func newVehicleLabel(name, tier string) render.Block {
	var blocks []render.Block
	if tier != "" {
		blocks = append(blocks, render.NewTextContent(render.Style{Font: &FontMedium, FontColor: FontMediumColor}, tier))
	}
	blocks = append(blocks, render.NewTextContent(render.Style{Font: &FontMedium, FontColor: FontMediumColor}, name))

	return render.NewBlocksContent(
		render.Style{
			Direction:  render.DirectionHorizontal,
			AlignItems: render.AlignItemsCenter,
			Gap:        5,
			// Debug:      true,
		},
		blocks...,
	)
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
