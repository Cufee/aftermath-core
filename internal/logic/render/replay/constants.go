package replay

import (
	"image/color"

	"github.com/cufee/aftermath-core/internal/logic/render"
)

var (
	frameStyle = render.Style{Direction: render.DirectionHorizontal, PaddingX: 20, PaddingY: 20, Gap: 20}

	overviewWidth   = 400.0
	playerCardWidth = 500.0

	hpBarColorAllies  = color.RGBA{R: 120, G: 255, B: 120, A: 255}
	hpBarColorEnemies = color.RGBA{R: 255, G: 120, B: 120, A: 255}
)

func defaultCardStyle(width, height float64) render.Style {
	return render.Style{
		Direction:       render.DirectionVertical,
		Width:           width,
		Height:          height,
		BackgroundColor: render.DefaultCardColor,
		PaddingX:        20,
		PaddingY:        10,
		BorderRadius:    20,
	}
}

func highlightCardStyle(width, height float64) render.Style {
	return render.Style{
		Direction:       render.DirectionVertical,
		Width:           width,
		Height:          height,
		BackgroundColor: color.RGBA{render.DefaultCardColor.R, render.DefaultCardColor.G, render.DefaultCardColor.B, render.DefaultCardColor.A + 30},
		PaddingX:        20,
		PaddingY:        10,
		BorderRadius:    20,
	}
}

func deadPlayerCardStyle(width, height float64) render.Style {
	return render.Style{
		Direction:       render.DirectionVertical,
		Width:           width,
		Height:          height,
		BackgroundColor: color.RGBA{render.DefaultCardColor.R - 10, render.DefaultCardColor.G - 10, render.DefaultCardColor.B - 10, render.DefaultCardColor.A},
		PaddingX:        20,
		PaddingY:        10,
		BorderRadius:    20,
	}
}
