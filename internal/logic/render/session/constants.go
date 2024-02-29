package session

import (
	"image"
	"image/color"

	"github.com/cufee/aftermath-core/internal/logic/render"
	"github.com/fogleman/gg"
)

func init() {
	{
		ctx := gg.NewContext(iconSize, iconSize)
		ctx.DrawRoundedRectangle(13, 2.5, 6, 17.5, 3)
		ctx.SetColor(color.RGBA{R: 255, G: 255, B: 255, A: 255})
		ctx.Fill()
		wn8Icon = ctx.Image()
	}

	{
		ctx := gg.NewContext(iconSize, 1)
		blankIconBlock = render.NewImageContent(render.Style{Width: float64(iconSize), Height: 1}, ctx.Image())
	}
}

var (
	iconSize       = 25
	wn8Icon        image.Image
	blankIconBlock render.Block

	BaseCardWidth       = 680.0
	BaseStatsBlockWidth = 120.0
	ClanPillWidth       = 80
)

func HighlightCardColor(base color.Color) color.Color {
	casted, ok := base.(color.RGBA)
	if !ok {
		return base
	}
	casted.R += 10
	casted.G += 10
	casted.B += 10
	return casted
}

func DefaultCardStyle(matchToImage image.Image) render.Style {
	style := render.Style{
		JustifyContent:  render.JustifyContentCenter,
		AlignItems:      render.AlignItemsCenter,
		Direction:       render.DirectionVertical,
		PaddingX:        20,
		PaddingY:        20,
		BackgroundColor: render.DefaultCardColor,
		BorderRadius:    20,
		Width:           BaseCardWidth,
		// Debug:           true,
	}
	return style
}

var DefaultStatsBlockStyle = render.Style{
	Direction:  render.DirectionVertical,
	AlignItems: render.AlignItemsCenter,
	Width:      BaseStatsBlockWidth,
	// Debug:      true,
}

func HighlightStatsBlockStyle(bgColor color.Color) render.Style {
	s := DefaultStatsBlockStyle
	s.PaddingY = 10
	s.BorderRadius = 10
	s.BackgroundColor = bgColor
	return s
}
