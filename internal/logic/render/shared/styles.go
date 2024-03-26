package shared

import "github.com/cufee/aftermath-core/internal/logic/render"

func FooterCardStyle() render.Style {
	backgroundColor := render.DefaultCardColor
	backgroundColor.A = 120
	return render.Style{
		JustifyContent:  render.JustifyContentCenter,
		AlignItems:      render.AlignItemsCenter,
		Direction:       render.DirectionVertical,
		PaddingX:        10,
		PaddingY:        5,
		BackgroundColor: backgroundColor,
		BorderRadius:    15,
		// Debug:           true,
	}
}
