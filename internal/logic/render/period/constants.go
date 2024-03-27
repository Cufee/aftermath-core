package period

import (
	"github.com/cufee/aftermath-core/dataprep/period"
	"github.com/cufee/aftermath-core/internal/logic/render"
)

type overviewStyle struct {
	container      render.Style
	blockContainer render.Style
}

type highlightStyle struct {
	container  render.Style
	cardTitle  render.Style
	tankName   render.Style
	blockLabel render.Style
	blockValue render.Style
}

func (s *overviewStyle) block(block period.StatsBlock) (render.Style, render.Style) {
	switch block.Flavor {
	case period.BlockFlavorSpecial:
		return render.Style{FontColor: render.TextPrimary, Font: &render.FontXL}, render.Style{FontColor: render.TextAlt, Font: &render.FontSmall}
	case period.BlockFlavorSecondary:
		return render.Style{FontColor: render.TextSecondary, Font: &render.FontMedium}, render.Style{FontColor: render.TextAlt, Font: &render.FontSmall}
	default:
		return render.Style{FontColor: render.TextPrimary, Font: &render.FontLarge}, render.Style{FontColor: render.TextAlt, Font: &render.FontSmall}
	}
}

func getOverviewStyle(columnWidth float64) overviewStyle {
	return overviewStyle{render.Style{
		Direction:      render.DirectionVertical,
		AlignItems:     render.AlignItemsCenter,
		JustifyContent: render.JustifyContentCenter,
		PaddingX:       15,
		PaddingY:       0,
		Gap:            10,
		Width:          columnWidth,
	}, render.Style{
		Direction:  render.DirectionVertical,
		AlignItems: render.AlignItemsCenter,
		// Debug:      true,
	}}
}

func defaultCardStyle(width float64) render.Style {
	style := render.Style{
		JustifyContent:  render.JustifyContentCenter,
		AlignItems:      render.AlignItemsCenter,
		Direction:       render.DirectionVertical,
		BackgroundColor: render.DefaultCardColor,
		BorderRadius:    20,
		PaddingY:        10,
		PaddingX:        20,
		Gap:             20,
		Width:           width,
		// Debug:           true,
	}
	return style
}

func titleCardStyle(width float64) render.Style {
	style := defaultCardStyle(width)
	style.PaddingX = style.PaddingY
	return style
}

func overviewCardStyle(width float64) render.Style {
	style := defaultCardStyle(width)
	style.AlignItems = render.AlignItemsEnd
	style.Direction = render.DirectionHorizontal
	style.JustifyContent = render.JustifyContentSpaceAround
	style.PaddingY = 25
	style.PaddingX = 0
	style.Gap = 0
	// style.Debug = true
	return style
}

func highlightCardStyle(containerStyle render.Style) highlightStyle {
	container := containerStyle
	container.Gap = 5
	container.PaddingX = 20
	container.PaddingY = 15
	container.Direction = render.DirectionHorizontal
	container.JustifyContent = render.JustifyContentSpaceBetween

	return highlightStyle{
		container:  container,
		cardTitle:  render.Style{Font: &render.FontSmall, FontColor: render.TextSecondary},
		tankName:   render.Style{Font: &render.FontMedium, FontColor: render.TextPrimary},
		blockValue: render.Style{Font: &render.FontMedium, FontColor: render.TextPrimary},
		blockLabel: render.Style{Font: &render.FontSmall, FontColor: render.TextAlt},
	}
}
