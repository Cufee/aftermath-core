package render

import (
	"image/color"

	"github.com/cufee/aftermath-core/internal/core/localization"
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
	// BaseCardWidth    = FrameWidth - (2 * FrameMargin)
	// BaseCardHeigh    = 150
	// BaseCardColor    = color.RGBA{30, 30, 30, 204}
	// DecorLinesColor  = color.RGBA{80, 80, 80, 255}

	BigTextColor   = color.RGBA{255, 255, 255, 255} // Session stats values, titles and names
	SmallTextColor = color.RGBA{204, 204, 204, 255} // Career stats values
	AltTextColor   = color.RGBA{100, 100, 100, 255} // Stats labels
	// PremiumColor  = color.RGBA{255, 223, 0, 255} // Premium Vehicle
	// VerifiedColor = color.RGBA{72, 167, 250, 255} // Verified Account
)

type RenderConfig struct {
	Font      font.Face
	FontColor color.RGBA
}

type BlockRenderConfig struct {
	Session RenderConfig `json:"session"`
	Career  RenderConfig `json:"career"`
	Label   RenderConfig `json:"label"`

	RowOptions RenderOptions                  `json:"rowOptions"`
	Locale     localization.SupportedLanguage `json:"locale"`
}

var defaultBlockRenderConfig BlockRenderConfig

func init() {
	fontFaces, ok := assets.GetFontFaces("default", 24, 18, 14)
	if !ok {
		panic("default font not found")
	}

	defaultBlockRenderConfig = BlockRenderConfig{
		Session: RenderConfig{
			Font:      fontFaces[24],
			FontColor: BigTextColor,
		},
		Career: RenderConfig{
			Font:      fontFaces[18],
			FontColor: SmallTextColor,
		},
		Label: RenderConfig{
			Font:      fontFaces[14],
			FontColor: AltTextColor,
		},
		RowOptions: RenderOptions{
			Style: Style{
				Direction:  DirectionVertical,
				AlignItems: AlignItemsCenter,
				Width:      100,
			},
			Debug: false,
		},
		Locale: localization.LanguageEN,
	}
}
