package session

import (
	"errors"
	"image"
	"image/color"

	"github.com/cufee/aftermath-core/internal/logic/render"
	"github.com/cufee/aftermath-core/internal/logic/render/assets"
	"golang.org/x/image/font"
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

var (
	BaseCardWidth       = 550.0
	BaseStatsBlockWidth = 120.0

	FontLarge  font.Face
	FontMedium font.Face
	FontSmall  font.Face

	FontLargeColor  = color.RGBA{255, 255, 255, 255} // Session stats values, titles and names
	FontMediumColor = color.RGBA{204, 204, 204, 255} // Career stats values
	FontSmallColor  = color.RGBA{150, 150, 150, 255} // Stats labels

	FontPlusColor    = color.RGBA{72, 167, 250, 255}
	FontPremiumColor = color.RGBA{255, 223, 0, 255}
)

func HighlightCardColor(base color.Color) color.Color {
	casted, ok := base.(color.RGBA)
	if !ok {
		return base
	}
	casted.R += 10
	casted.G += 10
	casted.B += 10
	return casted
}

func DefaultCardStyle(matchToImage image.Image) render.Style {
	style := render.Style{
		JustifyContent:  render.JustifyContentCenter,
		AlignItems:      render.AlignItemsCenter,
		Direction:       render.DirectionVertical,
		PaddingX:        20,
		PaddingY:        20,
		BackgroundColor: color.RGBA{10, 10, 10, 120},
		BorderRadius:    20,
		Width:           BaseCardWidth,
		// Debug:           true,
	}
	return style
}

type subscriptionPillStyle struct {
	Text      render.Style
	Icon      render.Style
	Container render.Style
}
type UserSubscriptionHeader struct {
	Name  string
	Icon  string
	Style subscriptionPillStyle
}

func (sub UserSubscriptionHeader) Block() (render.Block, error) {
	if tierImage, ok := assets.GetImage(sub.Icon); ok {
		blankStyle := sub.Style.Icon
		blankStyle.BackgroundColor = color.Transparent

		return render.NewBlocksContent(sub.Style.Container,
			render.NewImageContent(sub.Style.Icon, tierImage),
			render.NewTextContent(sub.Style.Text, sub.Name),
			// render.NewImageContent(blankStyle, tierImage),
		), nil
	}
	return render.Block{}, errors.New("tier icon not found")
}

var (
	UserSubscriptionSupporter = &UserSubscriptionHeader{
		Name: "Supporter",
		Icon: "images/icons/fire",
		Style: subscriptionPillStyle{
			Container: render.Style{Direction: render.DirectionHorizontal, AlignItems: render.AlignItemsCenter, BackgroundColor: DefaultCardStyle(nil).BackgroundColor, BorderRadius: 15, PaddingX: 10, PaddingY: 5},
			Icon:      render.Style{Width: 16, Height: 16, BackgroundColor: FontPlusColor, PaddingX: 5},
			Text:      render.Style{Font: &FontSmall, FontColor: FontMediumColor, PaddingX: 5},
		},
	}
	UserSubscriptionPlus = &UserSubscriptionHeader{
		Name: "Aftermath+",
		Icon: "images/icons/star",
		Style: subscriptionPillStyle{
			Container: render.Style{Direction: render.DirectionHorizontal, AlignItems: render.AlignItemsCenter, BackgroundColor: DefaultCardStyle(nil).BackgroundColor, BorderRadius: 15, PaddingX: 10, PaddingY: 5},
			Icon:      render.Style{Width: 24, Height: 24, BackgroundColor: FontPlusColor},
			Text:      render.Style{Font: &FontSmall, FontColor: FontMediumColor},
		},
	}
	UserSubscriptionPremium = &UserSubscriptionHeader{
		Name: "Aftermath Pro",
		Icon: "images/icons/star",
		Style: subscriptionPillStyle{
			Container: render.Style{Direction: render.DirectionHorizontal, AlignItems: render.AlignItemsCenter, BackgroundColor: DefaultCardStyle(nil).BackgroundColor, BorderRadius: 15, PaddingX: 10, PaddingY: 5},
			Icon:      render.Style{Width: 24, Height: 24, BackgroundColor: FontPremiumColor},
			Text:      render.Style{Font: &FontSmall, FontColor: FontMediumColor},
		},
	}
	UserSubscriptionPremiumXL = &UserSubscriptionHeader{
		Name: "Aftermath Pro",
		Icon: "images/icons/star-multiple",
		Style: subscriptionPillStyle{
			Container: render.Style{Direction: render.DirectionHorizontal, AlignItems: render.AlignItemsCenter, BackgroundColor: DefaultCardStyle(nil).BackgroundColor, BorderRadius: 15, PaddingX: 10, PaddingY: 5},
			Icon:      render.Style{Width: 24, Height: 24, BackgroundColor: FontPremiumColor},
			Text:      render.Style{Font: &FontSmall, FontColor: FontMediumColor},
		},
	}
)
