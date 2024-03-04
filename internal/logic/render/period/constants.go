package period

import (
	"github.com/cufee/aftermath-core/dataprep/period"
	"github.com/cufee/aftermath-core/internal/logic/render"
	"github.com/cufee/aftermath-core/internal/logic/render/shared"
)

type overviewStyle struct {
	container render.Style
}

type highlightStyle struct {
	container  render.Style
	cardTitle  render.Style
	tankName   render.Style
	blockLabel render.Style
	blockValue render.Style
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

func titleCardStyle(containerStyle render.Style) shared.TitleCardStyle {
	container := containerStyle
	container.AlignItems = render.AlignItemsCenter
	container.Direction = render.DirectionHorizontal
	container.PaddingX = 20

	return shared.TitleCardStyle{
		Container: container,
		Nickname:  render.Style{Font: &render.FontLarge, FontColor: render.TextPrimary},
		ClanTag:   render.Style{Font: &render.FontMedium, FontColor: render.TextSecondary},
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

func defaultCardStyle(width float64) render.Style {
	style := render.Style{
		JustifyContent:  render.JustifyContentCenter,
		AlignItems:      render.AlignItemsCenter,
		Direction:       render.DirectionVertical,
		PaddingY:        10,
		Gap:             20,
		BackgroundColor: render.DefaultCardColor,
		BorderRadius:    20,
		Width:           width,
		// Debug:           true,
	}
	return style
}

func overviewCardStyle(width float64) render.Style {
	style := defaultCardStyle(width)
	style.PaddingY = 25
	style.PaddingX = 20
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
