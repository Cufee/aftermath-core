package session

import (
	"errors"
	"image"
	"image/color"

	"github.com/cufee/aftermath-core/internal/logic/render"
	"github.com/cufee/aftermath-core/internal/logic/render/assets"
	"github.com/fogleman/gg"
)

func init() {
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
		BackgroundColor: render.DefaultCardColor,
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
	weight int
	Name   string
	Icon   string
	Style  subscriptionPillStyle
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
		weight: 25,
		Name:   "Supporter",
		Icon:   "images/icons/fire",
		Style: subscriptionPillStyle{
			Container: render.Style{Direction: render.DirectionHorizontal, AlignItems: render.AlignItemsCenter, BackgroundColor: DefaultCardStyle(nil).BackgroundColor, BorderRadius: 15, PaddingX: 10, PaddingY: 5, Height: 32},
			Icon:      render.Style{Width: 16, Height: 16, BackgroundColor: render.TextSubscriptionPlus, PaddingX: 5},
			Text:      render.Style{Font: &render.FontSmall, FontColor: render.TextSecondary, PaddingX: 5},
		},
	}
	userSubscriptionPlus = &subscriptionHeader{
		weight: 50,
		Name:   "Aftermath+",
		Icon:   "images/icons/star",
		Style: subscriptionPillStyle{
			Container: render.Style{Direction: render.DirectionHorizontal, AlignItems: render.AlignItemsCenter, BackgroundColor: DefaultCardStyle(nil).BackgroundColor, BorderRadius: 15, PaddingX: 10, PaddingY: 5, Height: 32},
			Icon:      render.Style{Width: 24, Height: 24, BackgroundColor: render.TextSubscriptionPlus},
			Text:      render.Style{Font: &render.FontSmall, FontColor: render.TextSecondary},
		},
	}
	userSubscriptionPro = &subscriptionHeader{
		weight: 75,
		Name:   "Aftermath Pro",
		Icon:   "images/icons/star",
		Style: subscriptionPillStyle{
			Container: render.Style{Direction: render.DirectionHorizontal, AlignItems: render.AlignItemsCenter, BackgroundColor: DefaultCardStyle(nil).BackgroundColor, BorderRadius: 15, PaddingX: 10, PaddingY: 5, Height: 32},
			Icon:      render.Style{Width: 24, Height: 24, BackgroundColor: render.TextSubscriptionPremium},
			Text:      render.Style{Font: &render.FontSmall, FontColor: render.TextSecondary},
		},
	}
	// Clans
	clanSubscriptionVerified = &subscriptionHeader{
		weight: 10,
		Icon:   "images/icons/verify",
		Style: subscriptionPillStyle{
			Icon:      render.Style{Width: 28, Height: 28, BackgroundColor: render.TextAlt},
			Container: render.Style{Direction: render.DirectionHorizontal},
		},
	}
	clanSubscriptionPro = &subscriptionHeader{
		weight: 75,
		Icon:   "images/icons/star-multiple",
		Style: subscriptionPillStyle{
			Icon:      render.Style{Width: 28, Height: 28, BackgroundColor: render.TextAlt},
			Container: render.Style{Direction: render.DirectionHorizontal},
		},
	}

	// Community
	subscriptionServerModerator = &subscriptionHeader{
		weight: 10000,
		Name:   "Community Moderator",
		Icon:   "images/icons/logo-128",
		Style: subscriptionPillStyle{
			Container: render.Style{Direction: render.DirectionHorizontal, AlignItems: render.AlignItemsCenter, BackgroundColor: DefaultCardStyle(nil).BackgroundColor, BorderRadius: 15, PaddingX: 10, PaddingY: 5, Gap: 5, Height: 32},
			Icon:      render.Style{Width: 20, Height: 20},
			Text:      render.Style{Font: &render.FontSmall, FontColor: render.TextSecondary},
		},
	}
	subscriptionServerBooster = &subscriptionHeader{
		Name:   "Booster",
		weight: 1,
		Icon:   "images/icons/discord-booster",
		Style: subscriptionPillStyle{
			Container: render.Style{Direction: render.DirectionHorizontal, AlignItems: render.AlignItemsCenter, BackgroundColor: DefaultCardStyle(nil).BackgroundColor, BorderRadius: 15, PaddingX: 10, PaddingY: 5, Gap: 5, Height: 32},
			Icon:      render.Style{Width: 20, Height: 20},
			Text:      render.Style{Font: &render.FontSmall, FontColor: render.TextSecondary},
		},
	}
)
