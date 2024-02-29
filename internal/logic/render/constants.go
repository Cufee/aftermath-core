package render

import (
	"image/color"

	"github.com/cufee/aftermath-core/internal/logic/render/assets"
	"golang.org/x/image/font"
)

var DiscordBackgroundColor = color.RGBA{49, 51, 56, 255}

var (
	FontExtraLarge font.Face
	FontLarge      font.Face
	FontMedium     font.Face
	FontSmall      font.Face

	TextPrimary   = color.RGBA{255, 255, 255, 255}
	TextSecondary = color.RGBA{204, 204, 204, 255}
	TextAlt       = color.RGBA{150, 150, 150, 255}

	TextSubscriptionPlus    = color.RGBA{72, 167, 250, 255}
	TextSubscriptionPremium = color.RGBA{255, 223, 0, 255}

	DefaultCardColor = color.RGBA{10, 10, 10, 180}
)

func init() {
	fontFaces, ok := assets.GetFontFaces("default", 36, 24, 18, 14)
	if !ok {
		panic("default font not found")
	}
	FontExtraLarge = fontFaces[36]
	FontLarge = fontFaces[24]
	FontMedium = fontFaces[18]
	FontSmall = fontFaces[14]
}
