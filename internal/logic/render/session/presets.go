package session

import (
	"fmt"
	"image/color"

	core "github.com/cufee/aftermath-core/internal/core/stats"
	"github.com/cufee/aftermath-core/internal/logic/render"
	"github.com/cufee/aftermath-core/internal/logic/render/assets"
	"golang.org/x/image/font"
)

var (
	// debugColorRed   = color.RGBA{255, 0, 0, 255}
	// debugColorPink = color.RGBA{255, 192, 203, 255}
	// debugColorGreen = color.RGBA{20, 160, 20, 255}

	// FontSizeHeader   = 36.0
	// FontSize         = 24.0
	// TextMargin       = FontSize / 2
	// FrameWidth       = 900
	// FrameMargin      = 50
	BaseCardWidth       = 550.0
	BaseStatsBlockWidth = 120.0
	BaseCardColor       = color.RGBA{30, 30, 30, 120}
	HighlightCardColor  = color.RGBA{50, 50, 50, 120}
	// DecorLinesColor  = color.RGBA{80, 80, 80, 255}

	FontLarge  font.Face
	FontMedium font.Face
	FontSmall  font.Face

	FontTranslucentColor = color.RGBA{255, 255, 255, 50}
	FontLargeColor       = color.RGBA{255, 255, 255, 255} // Session stats values, titles and names
	FontMediumColor      = color.RGBA{204, 204, 204, 255} // Career stats values
	FontSmallColor       = color.RGBA{150, 150, 150, 255} // Stats labels

	// FontPremiumColor = color.RGBA{255, 223, 0, 255}   // Premium Vehicle
	// FontVerifiedColor = color.RGBA{72, 167, 250, 255} // Verified Account
)

func init() {
	fontFaces, ok := assets.GetFontFaces("default", 24, 18, 14)
	if !ok {
		panic("default font not found")
	}
	FontLarge = fontFaces[24]
	FontMedium = fontFaces[18]
	FontSmall = fontFaces[14]
}

func NewPlayerTitleCard(name, clanTag string) render.Block {
	var content []render.Block

	justify := render.JustifyContentCenter
	if clanTag != "" {
		justify = render.JustifyContentSpaceBetween
		// Visible tag
		content = append(content, render.NewBlocksContent(render.Style{
			Direction:       render.DirectionHorizontal,
			AlignItems:      render.AlignItemsCenter,
			PaddingX:        10,
			PaddingY:        5,
			BackgroundColor: HighlightCardColor,
			BorderRadius:    10,
			// Debug:           true,
		}, render.NewTextContent(clanTag, render.Style{Font: FontMedium, FontColor: FontMediumColor})))
	}

	// Nickname
	content = append(content, render.NewTextContent(name, render.Style{Font: FontLarge, FontColor: FontLargeColor}))

	if clanTag != "" {
		// Invisible tag
		content = append(content, render.NewBlocksContent(render.Style{
			Direction:    render.DirectionHorizontal,
			AlignItems:   render.AlignItemsCenter,
			PaddingX:     10,
			PaddingY:     5,
			BorderRadius: 10,
			// Debug:        true,
		}, render.NewTextContent(clanTag, render.Style{Font: FontMedium, FontColor: color.RGBA{0, 0, 0, 0}})))
	}

	return render.NewBlocksContent(render.Style{
		JustifyContent:  justify,
		AlignItems:      render.AlignItemsCenter,
		Direction:       render.DirectionHorizontal,
		PaddingX:        20,
		PaddingY:        20,
		BackgroundColor: BaseCardColor,
		BorderRadius:    20,
		Width:           BaseCardWidth,
		// Debug:           true,
	}, content...)
}

func NewStatsBlock(label string, values ...any) render.Block {
	var content []render.Block
	for i, value := range values {
		style := render.Style{Font: FontLarge, FontColor: FontLargeColor}
		if i > 0 {
			style = render.Style{Font: FontMedium, FontColor: FontMediumColor}
		}
		content = append(content, render.NewTextContent(statsValueToString(value), style))
	}
	if label != "" {
		content = append(content, render.NewTextContent(label, render.Style{Font: FontSmall, FontColor: FontSmallColor}))
	}
	return render.NewBlocksContent(render.Style{
		Direction:  render.DirectionVertical,
		AlignItems: render.AlignItemsCenter,
		Width:      BaseStatsBlockWidth,
		// Debug:      true,
	}, content...)
}

func NewTextLabel(label string) render.Block {
	return render.NewTextContent(label, render.Style{Font: FontMedium, FontColor: FontSmallColor})
}

func NewVehicleLabel(name, tier string) render.Block {
	var blocks []render.Block
	if tier != "" {
		blocks = append(blocks, render.NewTextContent(tier, render.Style{Font: FontSmall, FontColor: FontMediumColor}))
	}
	blocks = append(blocks, render.NewTextContent(name, render.Style{Font: FontMedium, FontColor: FontMediumColor}))

	return render.NewBlocksContent(
		render.Style{
			Direction:  render.DirectionHorizontal,
			AlignItems: render.AlignItemsCenter,
			Gap:        5,
			// Debug:      true,
		},
		blocks...,
	)
}

func NewCardBlock(label render.Block, stats []render.Block) render.Block {
	var content []render.Block
	content = append(content, label)
	content = append(content, render.NewBlocksContent(render.Style{
		Direction:      render.DirectionHorizontal,
		JustifyContent: render.JustifyContentSpaceBetween,
		Gap:            10,
		// Debug:     true,
	}, stats...))

	return render.NewBlocksContent(
		render.Style{
			Direction:       render.DirectionVertical,
			AlignItems:      render.AlignItemsCenter,
			Gap:             5,
			PaddingX:        20,
			PaddingY:        20,
			BackgroundColor: BaseCardColor,
			BorderRadius:    20,
			Width:           BaseCardWidth,
			// Debug:           true,
		},
		content...)
}

func statsValueToString(value any) string {
	switch cast := value.(type) {
	case string:
		return cast
	case float64:
		if int(cast) == core.InvalidValue {
			return "-"
		}
		return fmt.Sprintf("%.2f%%", value)
	case int:
		if value == core.InvalidValue {
			return "-"
		}
		return fmt.Sprintf("%d", value)
	default:
		return fmt.Sprint(value)
	}
}
