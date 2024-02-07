package replay

import (
	"image/color"

	"github.com/cufee/aftermath-core/internal/logic/render"
	"github.com/fogleman/gg"
)

type progressDirection int

const (
	progressDirectionHorizontal progressDirection = iota
	progressDirectionVertical
)

const progressBarWidth = 8

func newProgressBar(size int, progress int, direction progressDirection, fillColor color.Color) render.Block {
	var width, height int
	if direction == progressDirectionHorizontal {
		width = (size)
		height = progressBarWidth
	} else {
		width = progressBarWidth
		height = (size)
	}

	ctx := gg.NewContext((width), (height))
	ctx.SetColor(color.RGBA{70, 70, 70, 255})
	ctx.DrawRoundedRectangle(0, 0, float64(width), float64(height), 5)
	ctx.Fill()

	if progress > 0 {
		ctx.SetColor(fillColor)
		if direction == progressDirectionHorizontal {
			ctx.DrawRoundedRectangle(0, 0, float64(progress)/100*float64(width), float64(height), 5)
		} else {
			ctx.DrawRoundedRectangle(0, float64(height)-float64(progress)/100*float64(height), float64(width), float64(progress)/100*float64(height), 5)
		}
		ctx.Fill()
	}

	if direction == progressDirectionHorizontal {
		return render.NewImageContent(render.Style{Width: float64(size), Height: progressBarWidth}, ctx.Image())
	}
	return render.NewImageContent(render.Style{Width: progressBarWidth, Height: float64(size)}, ctx.Image())
}
