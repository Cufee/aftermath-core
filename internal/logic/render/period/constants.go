package period

import (
	"github.com/cufee/aftermath-core/dataprep/period"
	"github.com/cufee/aftermath-core/internal/logic/render"
	"github.com/cufee/aftermath-core/internal/logic/render/shared"
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

func getOverviewStyle(columnWidth float64) overviewStyle {
	return overviewStyle{render.Style{
		Direction:      render.DirectionVertical,
		AlignItems:     render.AlignItemsCenter,
		JustifyContent: render.JustifyContentCenter,
		PaddingX:       10,
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
		Gap:             20,
		Width:           width,
		// Debug:           true,
	}
	return style
}

func overviewCardStyle(width float64) render.Style {
	style := defaultCardStyle(width)
	style.Direction = render.DirectionHorizontal
	style.JustifyContent = render.JustifyContentCenter
	style.PaddingY = 30
	style.PaddingX = 25
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

type ratingIconOptions struct {
	lines     int
	direction int

	lineStep  int
	lineWidth float64
	jump      float64
	gap       float64
}

func (opts ratingIconOptions) height() int {
	return ((opts.lines/2+1)*opts.lineStep + (opts.lines/2)*int(opts.jump))
}
func (opts ratingIconOptions) width() int {
	return opts.lines * (int(opts.lineWidth + opts.gap))
}

func defaultRatingIconOptions(direction int) ratingIconOptions {
	return ratingIconOptions{
		gap:       4,
		jump:      6,
		lines:     7,
		lineStep:  10,
		lineWidth: 6,
		direction: direction,
	}
}
