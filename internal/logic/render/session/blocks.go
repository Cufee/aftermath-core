package session

import (
	"github.com/cufee/aftermath-core/dataprep"
	"github.com/cufee/aftermath-core/dataprep/session"
	"github.com/cufee/aftermath-core/internal/logic/render"
)

type convertOptions struct {
	showSessionStats bool
	showCareerStats  bool
	showLabels       bool
	showIcons        bool
}

func statsBlocksToCardBlocks(stats []session.StatsBlock, blockWidth map[int]float64, opts ...convertOptions) ([]render.Block, error) {
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

		containerStyle := defaultStatsBlockStyle(blockWidth[index])
		if index == 0 {
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
