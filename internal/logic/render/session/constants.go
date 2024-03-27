package session

import (
	"image"
	"image/color"

	"github.com/cufee/aftermath-core/internal/logic/render"
	"github.com/fogleman/gg"
)

type blockStyle struct {
	session render.Style
	career  render.Style
	label   render.Style
}

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
)

var (
	promoTextStyle = render.Style{Font: &render.FontMedium, FontColor: render.TextPrimary}

	defaultBlockStyle = blockStyle{
		render.Style{Font: &render.FontLarge, FontColor: render.TextPrimary},
		render.Style{Font: &render.FontMedium, FontColor: render.TextSecondary},
		render.Style{Font: &render.FontSmall, FontColor: render.TextAlt},
	}
)

func highlightCardColor() color.Color {
	backgroundColor := render.DefaultCardColor
	backgroundColor.R += 10
	backgroundColor.G += 10
	backgroundColor.B += 10
	return backgroundColor
}

func defaultCardStyle(width float64) render.Style {
	style := render.Style{
		JustifyContent:  render.JustifyContentCenter,
		AlignItems:      render.AlignItemsCenter,
		Direction:       render.DirectionVertical,
		PaddingX:        20,
		PaddingY:        20,
		BackgroundColor: render.DefaultCardColor,
		BorderRadius:    20,
		Width:           width,
		// Debug:           true,
	}
	return style
}

func tierPercentageCardStyle(width float64) render.Style {
	style := defaultCardStyle(width)
	style.Direction = render.DirectionHorizontal
	style.BackgroundColor = nil
	style.BorderRadius = 0
	style.PaddingX = 0
	style.PaddingY = 5
	return style
}

func defaultStatsBlockStyle(width float64) render.Style {
	style := render.Style{
		Direction:  render.DirectionVertical,
		AlignItems: render.AlignItemsCenter,
		Width:      width,
		// Debug:      true,
	}
	return style
}

func statsRowStyle(width float64) render.Style {
	return render.Style{
		JustifyContent: render.JustifyContentSpaceBetween,
		Direction:      render.DirectionHorizontal,
		AlignItems:     render.AlignItemsCenter,
		Width:          width,
		Gap:            10,
	}
}

func highlightStatsBlockStyle(width float64) render.Style {
	s := defaultStatsBlockStyle(width)
	s.PaddingY = 10
	s.BorderRadius = 10
	s.BackgroundColor = highlightCardColor()
	return s
}
