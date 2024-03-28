package session

import (
	"github.com/cufee/aftermath-core/dataprep"
	"github.com/cufee/aftermath-core/dataprep/session"
	"github.com/cufee/aftermath-core/internal/logic/render"
)

type convertOptions struct {
	showSessionStats    bool
	showCareerStats     bool
	showLabels          bool
	showIcons           bool
	highlightBlockIndex int
}

func newVehicleCard(style render.Style, card dataprep.StatsCard[session.StatsBlock, string], sizes map[int]float64, opts convertOptions) (render.Block, error) {
	if card.Type == dataprep.CardTypeRatingVehicle {
		return slimVehicleCard(style, card, sizes, opts)
	}

	return defaultVehicleCard(style, card, sizes, opts)
}

func defaultVehicleCard(style render.Style, card dataprep.StatsCard[session.StatsBlock, string], sizes map[int]float64, opts convertOptions) (render.Block, error) {
	blocks, err := statsBlocksToCardBlocks(card.Blocks, sizes, opts)
	if err != nil {
		return render.Block{}, err
	}

	cardContentBlocks := []render.Block{newCardTitle(card.Title)}
	contentWidth := style.Width - style.PaddingX*2

	statsRowBlock := render.NewBlocksContent(statsRowStyle(contentWidth), blocks...)
	cardContentBlocks = append(cardContentBlocks, statsRowBlock)

	return render.NewBlocksContent(style, cardContentBlocks...), nil
}

func slimVehicleCard(style render.Style, card dataprep.StatsCard[session.StatsBlock, string], sizes map[int]float64, opts convertOptions) (render.Block, error) {
	opts.highlightBlockIndex = -1
	opts.showCareerStats = false
	opts.showLabels = false
	opts.showIcons = true

	blocks, err := statsBlocksToCardBlocks(card.Blocks, sizes, opts)
	if err != nil {
		return render.Block{}, err
	}

	titleBlock := render.NewTextContent(ratingVehicleTitleStyle, card.Title)
	statsRowBlock := render.NewBlocksContent(statsRowStyle(0), blocks...)

	containerStyle := style
	containerStyle.Direction = render.DirectionHorizontal
	containerStyle.JustifyContent = render.JustifyContentSpaceBetween

	return render.NewBlocksContent(containerStyle, titleBlock, statsRowBlock), nil
}

func statsBlocksToCardBlocks(stats []session.StatsBlock, blockWidth map[int]float64, opts ...convertOptions) ([]render.Block, error) {
	var options convertOptions = convertOptions{
		showSessionStats:    true,
		showCareerStats:     true,
		showLabels:          true,
		showIcons:           true,
		highlightBlockIndex: 0,
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

		containerStyle := defaultStatsBlockStyle(blockWidth[index])
		if index == options.highlightBlockIndex {
			containerStyle = highlightStatsBlockStyle(blockWidth[index])
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

func newCardTitle(label string) render.Block {
	return render.NewTextContent(defaultBlockStyle.career, label)
}
