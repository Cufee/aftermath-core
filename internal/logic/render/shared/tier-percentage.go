package shared

import (
	"fmt"
	"image/color"

	"github.com/cufee/aftermath-core/internal/core/database/models"
	"github.com/cufee/aftermath-core/internal/core/stats"
	"github.com/cufee/aftermath-core/internal/logic/render"
)

func NewTierPercentageCard(style render.Style, vehicles map[int]*stats.ReducedVehicleStats, glossary map[int]models.Vehicle) render.Block {
	var blocks []render.Block
	var elements int = 10

	backgroundSharePrimary := render.DefaultCardColor
	backgroundShareSecondary := color.RGBA{120, 120, 120, 120}

	for i := range elements {
		shade := backgroundSharePrimary
		if i%2 == 0 {
			shade = backgroundShareSecondary
		}

		blocks = append(blocks, render.NewBlocksContent(render.Style{
			BackgroundColor: shade,
			Width:           style.Width / float64(elements),
			JustifyContent:  render.JustifyContentCenter,
		}, render.NewTextContent(render.Style{Font: &render.FontMedium, FontColor: render.TextPrimary}, fmt.Sprint(i))))
	}

	return render.NewBlocksContent(style, blocks...)

}
