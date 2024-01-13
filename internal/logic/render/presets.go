package render

import (
	"fmt"
	"image/color"

	core "github.com/cufee/aftermath-core/internal/core/stats"
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
	BaseCardWidth       = 500.0
	BaseStatsBlockWidth = 107.5
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

func NewPlayerTitleCard(name, clanTag string) Block {
	var content []Block
	// Visible tag
	content = append(content, NewBlocksContent(Style{
		Direction:       DirectionHorizontal,
		AlignItems:      AlignItemsCenter,
		PaddingX:        10,
		PaddingY:        5,
		BackgroundColor: HighlightCardColor,
		BorderRadius:    10,
		// Debug:           true,
	}, NewTextContent(clanTag, Style{Font: FontMedium, FontColor: FontMediumColor})))

	// Nickname
	content = append(content, NewTextContent(name, Style{Font: FontLarge, FontColor: FontLargeColor}))

	// Invisible tag
	content = append(content, NewBlocksContent(Style{
		Direction:    DirectionHorizontal,
		AlignItems:   AlignItemsCenter,
		PaddingX:     10,
		PaddingY:     5,
		BorderRadius: 10,
		// Debug:        true,
	}, NewTextContent(clanTag, Style{Font: FontMedium, FontColor: color.RGBA{0, 0, 0, 0}})))

	return NewBlocksContent(Style{
		Direction:       DirectionHorizontal,
		JustifyContent:  JustifyContentSpaceBetween,
		AlignItems:      AlignItemsCenter,
		PaddingX:        20,
		PaddingY:        20,
		BackgroundColor: BaseCardColor,
		BorderRadius:    20,
		Width:           BaseCardWidth,
		// Debug:           true,
	}, content...)
}

func NewStatsBlock(label string, values ...any) Block {
	var content []Block
	for i, value := range values {
		style := Style{Font: FontLarge, FontColor: FontLargeColor}
		if i > 0 {
			style = Style{Font: FontMedium, FontColor: FontMediumColor}
		}
		content = append(content, NewTextContent(statsValueToString(value), style))
	}
	if label != "" {
		content = append(content, NewTextContent(label, Style{Font: FontSmall, FontColor: FontSmallColor}))
	}
	return NewBlocksContent(Style{
		Direction:  DirectionVertical,
		AlignItems: AlignItemsCenter,
		Width:      BaseStatsBlockWidth,
		// Debug:      true,
	}, content...)
}

func NewTextLabel(label string) Block {
	return NewTextContent(label, Style{Font: FontMedium, FontColor: FontSmallColor})
}

func NewVehicleLabel(name, tier string) Block {
	return NewBlocksContent(
		Style{
			Direction:  DirectionHorizontal,
			AlignItems: AlignItemsCenter,
			Gap:        5,
			// Debug:      true,
		},
		NewTextContent(tier, Style{Font: FontSmall, FontColor: FontMediumColor}),
		NewTextContent(name, Style{Font: FontMedium, FontColor: FontMediumColor}),
	)
}

func NewCardBlock(label Block, stats []Block) Block {
	var content []Block
	content = append(content, label)
	content = append(content, NewBlocksContent(Style{
		Direction:      DirectionHorizontal,
		JustifyContent: JustifyContentSpaceBetween,
		Gap:            10,
		// Debug:     true,
	}, stats...))

	return NewBlocksContent(
		Style{
			Direction:       DirectionVertical,
			AlignItems:      AlignItemsCenter,
			Gap:             0,
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
	if value == core.InvalidValue {
		return "-"
	}
	switch cast := value.(type) {
	case string:
		return cast
	case float64:
		return fmt.Sprintf("%.2f%%", value)
	case int:
		return fmt.Sprintf("%d", value)
	default:
		return fmt.Sprint(value)
	}
}
