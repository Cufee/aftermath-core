package session

import (
	"image"
	"image/color"

	"github.com/cufee/aftermath-core/dataprep"
	"github.com/cufee/aftermath-core/dataprep/session"
	"github.com/cufee/aftermath-core/internal/core/stats"
	"github.com/cufee/aftermath-core/internal/logic/render"
	"github.com/cufee/aftermath-core/internal/logic/render/assets"
	"github.com/cufee/aftermath-core/internal/logic/render/shared"
)

func styleBlocks(blocks []session.StatsBlock, styles ...render.Style) []styledStatsBlock {
	var lastStyle render.Style
	var styledBlocks []styledStatsBlock
	for i, block := range blocks {
		if i < len(styles) {
			lastStyle = styles[i]
		}
		styledBlocks = append(styledBlocks, styledStatsBlock{StatsBlock: block, style: lastStyle})
	}
	return styledBlocks
}

func comparisonIconFromBlock(block session.StatsBlock) render.Block {
	if !stats.ValueValid(block.Session.Value) || !stats.ValueValid(block.Career.Value) {
		return blankIconBlock
	}

	if block.Tag == dataprep.TagWN8 {
		// WN8 icons need to show the color
		return blockToWN8Icon(block.Session, block.Tag)
	}

	var icon image.Image
	var iconColor color.Color
	if block.Session.Value > block.Career.Value {
		icon, _ = assets.GetImage("images/icons/chevron-up-single")
		iconColor = color.RGBA{R: 0, G: 255, B: 0, A: 255}
	}
	if block.Session.Value < block.Career.Value {
		icon, _ = assets.GetImage("images/icons/chevron-down-single")
		iconColor = color.RGBA{R: 255, G: 0, B: 0, A: 255}
	}
	if icon == nil {
		return blankIconBlock
	}

	return render.NewImageContent(render.Style{Width: 25, Height: 25, BackgroundColor: iconColor}, icon)
}

func blockToWN8Icon(value dataprep.Value, tag dataprep.Tag) render.Block {
	if tag != dataprep.TagWN8 || !stats.ValueValid(value.Value) {
		return blankIconBlock
	}
	return render.NewImageContent(render.Style{Width: 25, Height: 25, BackgroundColor: shared.GetWN8Color(int(value.Value))}, wn8Icon)
}
