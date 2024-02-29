package period

import (
	"github.com/cufee/aftermath-core/dataprep/period"
	"github.com/cufee/aftermath-core/internal/logic/render"
)

type overviewStyle struct {
	container render.Style
}

func (s *overviewStyle) block(flavor period.BlockFlavor) (render.Style, render.Style) {
	switch flavor {
	case period.BlockFlavorSpecial:
		return render.Style{FontColor: render.TextPrimary, Font: &render.FontXL}, render.Style{FontColor: render.TextAlt, Font: &render.FontSmall}
	case period.BlockFlavorSecondary:
		return render.Style{FontColor: render.TextSecondary, Font: &render.FontMedium}, render.Style{FontColor: render.TextAlt, Font: &render.FontSmall}
	default:
		return render.Style{FontColor: render.TextPrimary, Font: &render.FontLarge}, render.Style{FontColor: render.TextAlt, Font: &render.FontSmall}
	}
}

func getOverviewStyle(width float64) overviewStyle {
	return overviewStyle{render.Style{
		Width:          width,
		AlignItems:     render.AlignItemsCenter,
		JustifyContent: render.JustifyContentCenter,
		// Debug:          true,
	}}
}

func DefaultCardStyle(width float64) render.Style {
	style := render.Style{
		JustifyContent:  render.JustifyContentCenter,
		AlignItems:      render.AlignItemsCenter,
		Direction:       render.DirectionVertical,
		PaddingX:        20,
		PaddingY:        40,
		Gap:             20,
		BackgroundColor: render.DefaultCardColor,
		BorderRadius:    20,
		Width:           width,
		// Debug:           true,
	}
	return style
}
