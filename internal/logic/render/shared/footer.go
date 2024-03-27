package shared

import "github.com/cufee/aftermath-core/internal/logic/render"

func NewFooterCard(text string) render.Block {
	backgroundColor := render.DefaultCardColor
	backgroundColor.A = 120
	return render.NewBlocksContent(render.Style{
		JustifyContent:  render.JustifyContentCenter,
		AlignItems:      render.AlignItemsCenter,
		Direction:       render.DirectionVertical,
		PaddingX:        12.5,
		PaddingY:        5,
		BackgroundColor: backgroundColor,
		BorderRadius:    15,
		// Debug:           true,
	}, render.NewTextContent(render.Style{Font: &render.FontSmall, FontColor: render.TextSecondary}, text))
}
