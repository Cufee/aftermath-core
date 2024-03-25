package replay

import (
	"image/color"

	"github.com/cufee/aftermath-core/internal/logic/render"
)

var (
	frameStyle = render.Style{Direction: render.DirectionVertical, PaddingX: 30, PaddingY: 30, Gap: 10}

	hpBarColorAllies  = color.RGBA{R: 120, G: 255, B: 120, A: 255}
	hpBarColorEnemies = color.RGBA{R: 255, G: 120, B: 120, A: 255}

	protagonistColor = color.RGBA{255, 223, 0, 255}
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
