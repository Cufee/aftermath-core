package replay

import (
	"image/color"

	"github.com/cufee/aftermath-core/internal/logic/render"
)

var (
	frameStyle = render.Style{Direction: render.DirectionVertical, PaddingX: 20, PaddingY: 20, Gap: 20}

	overviewWidth       = 400.0
	playerCardNameWidth = 300.0

	hpBarColorAllies  = color.RGBA{R: 120, G: 255, B: 120, A: 255}
	hpBarColorEnemies = color.RGBA{R: 255, G: 120, B: 120, A: 255}
)

func defaultCardStyle(width, height float64) render.Style {
	return render.Style{
		Direction:       render.DirectionVertical,
		Width:           width + 40,
		Height:          height,
		BackgroundColor: render.DefaultCardColor,
		PaddingX:        10,
		BorderRadius:    15,
	}
}

func playerCardStyle(presets []blockPreset) render.Style {
	var blocksWidth float64
	for _, preset := range presets {
		blocksWidth += preset.width
	}
	return defaultCardStyle(playerCardNameWidth+blocksWidth, 0)
}
