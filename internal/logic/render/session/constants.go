package session

import (
	"errors"
	"image"
	"image/color"

	"github.com/cufee/aftermath-core/internal/logic/render"
	"github.com/cufee/aftermath-core/internal/logic/render/assets"
	"github.com/fogleman/gg"
	"golang.org/x/image/font"
)

func init() {
	{
		fontFaces, ok := assets.GetFontFaces("default", 24, 18, 14)
		if !ok {
			panic("default font not found")
		}
		FontLarge = fontFaces[24]
		FontMedium = fontFaces[18]
		FontSmall = fontFaces[14]
	}

	{
		ctx := gg.NewContext(iconSize, iconSize)
		ctx.DrawRoundedRectangle(13, 2.5, 6, 17.5, 3)
		ctx.SetColor(color.RGBA{R: 255, G: 255, B: 255, A: 255})
		ctx.Fill()
		wn8Icon = ctx.Image()
	}

	{
		ctx := gg.NewContext(iconSize, 1)
		blankIconBlock = render.NewImageContent(render.Style{Width: float64(iconSize), Height: 1}, ctx.Image())
	}
}

var (
	iconSize       = 25
	wn8Icon        image.Image
	blankIconBlock render.Block

	BaseCardWidth       = 680.0
	BaseStatsBlockWidth = 120.0
	ClanPillWidth       = 80

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

var DefaultStatsBlockStyle = render.Style{
	Direction:  render.DirectionVertical,
	AlignItems: render.AlignItemsCenter,
	Width:      BaseStatsBlockWidth,
	// Debug:      true,
}

func HighlightStatsBlockStyle(bgColor color.Color) render.Style {
	s := DefaultStatsBlockStyle
	s.PaddingY = 10
	s.BorderRadius = 10
	s.BackgroundColor = bgColor
	return s
}

type subscriptionPillStyle struct {
	Text      render.Style
	Icon      render.Style
	Container render.Style
}
type subscriptionHeader struct {
	Name  string
	Icon  string
	Style subscriptionPillStyle
}

func (sub subscriptionHeader) Block() (render.Block, error) {
	if tierImage, ok := assets.GetImage(sub.Icon); ok {
		content := []render.Block{render.NewImageContent(sub.Style.Icon, tierImage)}
		if sub.Name != "" {
			content = append(content, render.NewTextContent(sub.Style.Text, sub.Name))
		}
		return render.NewBlocksContent(sub.Style.Container, content...), nil
	}
	return render.Block{}, errors.New("tier icon not found")
}

var (
	// Personal
	userSubscriptionSupporter = &subscriptionHeader{
		Name: "Supporter",
		Icon: "images/icons/fire",
		Style: subscriptionPillStyle{
			Container: render.Style{Direction: render.DirectionHorizontal, AlignItems: render.AlignItemsCenter, BackgroundColor: DefaultCardStyle(nil).BackgroundColor, BorderRadius: 15, PaddingX: 10, PaddingY: 5},
			Icon:      render.Style{Width: 16, Height: 16, BackgroundColor: FontPlusColor, PaddingX: 5},
			Text:      render.Style{Font: &FontSmall, FontColor: FontMediumColor, PaddingX: 5},
		},
	}
	userSubscriptionPlus = &subscriptionHeader{
		Name: "Aftermath+",
		Icon: "images/icons/star",
		Style: subscriptionPillStyle{
			Container: render.Style{Direction: render.DirectionHorizontal, AlignItems: render.AlignItemsCenter, BackgroundColor: DefaultCardStyle(nil).BackgroundColor, BorderRadius: 15, PaddingX: 10, PaddingY: 5},
			Icon:      render.Style{Width: 24, Height: 24, BackgroundColor: FontPlusColor},
			Text:      render.Style{Font: &FontSmall, FontColor: FontMediumColor},
		},
	}
	userSubscriptionPro = &subscriptionHeader{
		Name: "Aftermath Pro",
		Icon: "images/icons/star",
		Style: subscriptionPillStyle{
			Container: render.Style{Direction: render.DirectionHorizontal, AlignItems: render.AlignItemsCenter, BackgroundColor: DefaultCardStyle(nil).BackgroundColor, BorderRadius: 15, PaddingX: 10, PaddingY: 5},
			Icon:      render.Style{Width: 24, Height: 24, BackgroundColor: FontPremiumColor},
			Text:      render.Style{Font: &FontSmall, FontColor: FontMediumColor},
		},
	}
	// Clans
	clanSubscriptionVerified = &subscriptionHeader{
		Icon: "images/icons/verify",
		Style: subscriptionPillStyle{
			Icon:      render.Style{Width: 28, Height: 28, BackgroundColor: FontSmallColor},
			Container: render.Style{Direction: render.DirectionHorizontal},
		},
	}
	clanSubscriptionPro = &subscriptionHeader{
		Icon: "images/icons/star-multiple",
		Style: subscriptionPillStyle{
			Icon:      render.Style{Width: 28, Height: 28, BackgroundColor: FontSmallColor},
			Container: render.Style{Direction: render.DirectionHorizontal},
		},
	}
)
