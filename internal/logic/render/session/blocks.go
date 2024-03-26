package session

import (
	"image/color"

	"github.com/cufee/aftermath-core/dataprep"
	"github.com/cufee/aftermath-core/dataprep/session"
	"github.com/cufee/aftermath-core/internal/logic/render"
	"github.com/rs/zerolog/log"
)

type convertOptions struct {
	showSessionStats bool
	showCareerStats  bool
	showLabels       bool
	showIcons        bool
}

func statsBlocksToCardBlocks(stats []session.StatsBlock, blockWidth map[dataprep.Tag]float64, opts ...convertOptions) ([]render.Block, error) {
	var options convertOptions = convertOptions{
		showSessionStats: true,
		showCareerStats:  true,
		showLabels:       true,
		showIcons:        true,
	}
	if len(opts) > 0 {
		options = opts[0]
	}

	var content []render.Block
	for index, statsBlock := range stats {
		blocks := make([]render.Block, 0, 3)
		if options.showSessionStats {
			if options.showIcons && statsBlock.Tag != dataprep.TagBattles {
				blocks = append(blocks, newStatsBlockRow(defaultBlockStyle.session, statsBlock.Session.String, comparisonIconFromBlock(statsBlock)))
			} else {
				blocks = append(blocks, render.NewTextContent(defaultBlockStyle.session, statsBlock.Session.String))
			}
		}
		if options.showCareerStats && statsBlock.Career.String != "" {
			if options.showIcons && statsBlock.Tag != dataprep.TagBattles {
				blocks = append(blocks, newStatsBlockRow(defaultBlockStyle.career, statsBlock.Career.String, blockToWN8Icon(statsBlock.Career, statsBlock.Tag)))
			} else {
				blocks = append(blocks, render.NewTextContent(defaultBlockStyle.career, statsBlock.Career.String))
			}
		}
		if options.showLabels && statsBlock.Tag != dataprep.TagBattles {
			if options.showIcons {
				blocks = append(blocks, newStatsBlockRow(defaultBlockStyle.label, statsBlock.Label, blankIconBlock))
			} else {
				blocks = append(blocks, render.NewTextContent(defaultBlockStyle.label, statsBlock.Label))
			}
		}

		containerStyle := defaultStatsBlockStyle(blockWidth[statsBlock.Tag])
		if index == 0 {
			containerStyle = highlightStatsBlockStyle(blockWidth[statsBlock.Tag])
		}
		content = append(content, render.NewBlocksContent(containerStyle, blocks...))
	}
	return content, nil
}

func newStatsBlockRow(style render.Style, value string, icon render.Block) render.Block {
	return render.NewBlocksContent(
		render.Style{Direction: render.DirectionHorizontal, AlignItems: render.AlignItemsCenter},
		icon,
		render.NewTextContent(style, value),
	)
}

func newPlayerTitleCard(style render.Style, name string, clanTagBlocks []render.Block) render.Block {
	if len(clanTagBlocks) == 0 {
		return render.NewBlocksContent(style, render.NewTextContent(playerNameStyle, name))
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
		BackgroundColor: highlightCardColor(),
		BorderRadius:    10,
		// Debug:           true,
	}, clanTagBlocks...)

	clanTagImage, err := clanTagBlock.Render()
	if err != nil {
		log.Warn().Err(err).Msg("failed to render clan tag")
		// This error is not fatal, we can just render the name
		return render.NewBlocksContent(style, render.NewTextContent(playerNameStyle, name))
	}
	content = append(content, render.NewImageContent(render.Style{Width: float64(clanTagImage.Bounds().Dx()), Height: float64(clanTagImage.Bounds().Dy())}, clanTagImage))

	// Nickname
	content = append(content, render.NewTextContent(playerNameStyle, name))

	// Invisible tag to offset the nickname
	invisibleStyle := clanTagStyle
	invisibleStyle.FontColor = color.Transparent
	clanBlock := render.NewBlocksContent(render.Style{
		Width:          float64(clanTagImage.Bounds().Dx()),
		JustifyContent: render.JustifyContentEnd,
	}, render.NewTextContent(invisibleStyle, "-"))

	content = append(content, clanBlock)

	containerStyle := style
	containerStyle.JustifyContent = render.JustifyContentSpaceBetween
	return render.NewBlocksContent(containerStyle, content...)
}

func newCardTitle(label string) render.Block {
	return render.NewTextContent(defaultBlockStyle.career, label)
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
