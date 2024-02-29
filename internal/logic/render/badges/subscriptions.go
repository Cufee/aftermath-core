package badges

import (
	"errors"
	"image/color"
	"slices"

	"github.com/cufee/aftermath-core/internal/core/database/models"
	"github.com/cufee/aftermath-core/internal/logic/render"
	"github.com/cufee/aftermath-core/internal/logic/render/assets"
)

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
	subscriptionWeight = map[models.SubscriptionType]int{
		models.SubscriptionTypeDeveloper: 999,
		// Moderators
		models.SubscriptionTypeServerModerator:  99,
		models.SubscriptionTypeContentModerator: 98,
		// Paid
		models.SubscriptionTypePro:     89,
		models.SubscriptionTypeProClan: 88,
		models.SubscriptionTypePlus:    79,
		//
		models.SubscriptionTypeSupporter:     29,
		models.SubscriptionTypeServerBooster: 28,
		//
		models.SubscriptionTypeVerifiedClan: 19,
	}

	// Personal
	userSubscriptionSupporter = &subscriptionHeader{
		Name: "Supporter",
		Icon: "images/icons/fire",
		Style: subscriptionPillStyle{
			Container: render.Style{Direction: render.DirectionHorizontal, AlignItems: render.AlignItemsCenter, BackgroundColor: render.DefaultCardColor, BorderRadius: 15, PaddingX: 7, PaddingY: 5, Height: 32},
			Icon:      render.Style{Width: 16, Height: 16, BackgroundColor: render.TextSubscriptionPlus},
			Text:      render.Style{Font: &render.FontSmall, FontColor: render.TextSecondary, PaddingX: 5},
		},
	}
	userSubscriptionPlus = &subscriptionHeader{
		Name: "Aftermath+",
		Icon: "images/icons/star",
		Style: subscriptionPillStyle{
			Container: render.Style{Direction: render.DirectionHorizontal, AlignItems: render.AlignItemsCenter, BackgroundColor: render.DefaultCardColor, BorderRadius: 15, PaddingX: 5, PaddingY: 5, Height: 32},
			Icon:      render.Style{Width: 24, Height: 24, BackgroundColor: render.TextSubscriptionPlus},
			Text:      render.Style{Font: &render.FontSmall, FontColor: render.TextSecondary, PaddingX: 5},
		},
	}
	userSubscriptionPro = &subscriptionHeader{
		Name: "Aftermath Pro",
		Icon: "images/icons/star",
		Style: subscriptionPillStyle{
			Container: render.Style{Direction: render.DirectionHorizontal, AlignItems: render.AlignItemsCenter, BackgroundColor: render.DefaultCardColor, BorderRadius: 15, PaddingX: 5, PaddingY: 5, Height: 32},
			Icon:      render.Style{Width: 24, Height: 24, BackgroundColor: render.TextSubscriptionPremium},
			Text:      render.Style{Font: &render.FontSmall, FontColor: render.TextSecondary, PaddingX: 5},
		},
	}
	// Clans
	clanSubscriptionVerified = &subscriptionHeader{
		Icon: "images/icons/verify",
		Style: subscriptionPillStyle{
			Icon:      render.Style{Width: 28, Height: 28, BackgroundColor: render.TextAlt},
			Container: render.Style{Direction: render.DirectionHorizontal},
		},
	}
	clanSubscriptionPro = &subscriptionHeader{
		Icon: "images/icons/star-multiple",
		Style: subscriptionPillStyle{
			Icon:      render.Style{Width: 28, Height: 28, BackgroundColor: render.TextAlt},
			Container: render.Style{Direction: render.DirectionHorizontal},
		},
	}

	// Community
	subscriptionDeveloper = &subscriptionHeader{
		Name: "Developer",
		Icon: "images/icons/github",
		Style: subscriptionPillStyle{
			Container: render.Style{Direction: render.DirectionHorizontal, AlignItems: render.AlignItemsCenter, BackgroundColor: color.RGBA{64, 32, 128, 180}, BorderRadius: 15, PaddingX: 6, PaddingY: 5, Gap: 5, Height: 32},
			Icon:      render.Style{Width: 20, Height: 20, BackgroundColor: render.TextPrimary},
			Text:      render.Style{Font: &render.FontSmall, FontColor: render.TextPrimary, PaddingX: 5},
		},
	}
	subscriptionServerModerator = &subscriptionHeader{
		Name: "Community Moderator",
		Icon: "images/icons/logo-128",
		Style: subscriptionPillStyle{
			Container: render.Style{Direction: render.DirectionHorizontal, AlignItems: render.AlignItemsCenter, BackgroundColor: render.DefaultCardColor, BorderRadius: 15, PaddingX: 7, PaddingY: 5, Gap: 5, Height: 32},
			Icon:      render.Style{Width: 20, Height: 20},
			Text:      render.Style{Font: &render.FontSmall, FontColor: render.TextSecondary, PaddingX: 2},
		},
	}
	subscriptionContentModerator = &subscriptionHeader{
		Name: "Moderator",
		Icon: "images/icons/logo-128",
		Style: subscriptionPillStyle{
			Container: render.Style{Direction: render.DirectionHorizontal, AlignItems: render.AlignItemsCenter, BackgroundColor: render.DefaultCardColor, BorderRadius: 15, PaddingX: 7, PaddingY: 5, Gap: 5, Height: 32},
			Icon:      render.Style{Width: 20, Height: 20},
			Text:      render.Style{Font: &render.FontSmall, FontColor: render.TextSecondary, PaddingX: 2},
		},
	}
	subscriptionServerBooster = &subscriptionHeader{
		Name: "Booster",
		Icon: "images/icons/discord-booster",
		Style: subscriptionPillStyle{
			Container: render.Style{Direction: render.DirectionHorizontal, AlignItems: render.AlignItemsCenter, BackgroundColor: render.DefaultCardColor, BorderRadius: 15, PaddingX: 10, PaddingY: 5, Gap: 5, Height: 32},
			Icon:      render.Style{Width: 20, Height: 20},
			Text:      render.Style{Font: &render.FontSmall, FontColor: render.TextSecondary},
		},
	}
	subscriptionTranslator = &subscriptionHeader{
		Name: "Translator",
		Icon: "images/icons/translator",
		Style: subscriptionPillStyle{
			Container: render.Style{Direction: render.DirectionHorizontal, AlignItems: render.AlignItemsCenter, BackgroundColor: render.DefaultCardColor, BorderRadius: 15, PaddingX: 10, PaddingY: 5, Gap: 5, Height: 32},
			Icon:      render.Style{Width: 20, Height: 20, BackgroundColor: render.TextPrimary},
			Text:      render.Style{Font: &render.FontSmall, FontColor: render.TextSecondary},
		},
	}
)

func SubscriptionsBadges(subscriptions []models.UserSubscription) ([]render.Block, error) {
	slices.SortFunc(subscriptions, func(i, j models.UserSubscription) int {
		return subscriptionWeight[j.Type] - subscriptionWeight[i.Type]
	})

	var badges []render.Block
	for _, subscription := range subscriptions {
		var header *subscriptionHeader
		switch subscription.Type {
		case models.SubscriptionTypeDeveloper:
			header = subscriptionDeveloper
		case models.SubscriptionTypeServerModerator:
			header = subscriptionServerModerator
		case models.SubscriptionTypeContentModerator:
			header = subscriptionContentModerator
		}

		if header != nil {
			block, err := header.Block()
			if err != nil {
				return nil, err
			}
			badges = append(badges, block)
			break
		}
	}
	for _, subscription := range subscriptions {
		var header *subscriptionHeader
		switch subscription.Type {
		case models.SubscriptionTypeContentTranslator:
			header = subscriptionTranslator
		}

		if header != nil {
			block, err := header.Block()
			if err != nil {
				return nil, err
			}
			badges = append(badges, block)
			break
		}
	}
	for _, subscription := range subscriptions {
		var header *subscriptionHeader
		switch subscription.Type {
		case models.SubscriptionTypePro:
			header = userSubscriptionPro
		case models.SubscriptionTypePlus:
			header = userSubscriptionPlus
		case models.SubscriptionTypeServerBooster:
			header = subscriptionServerBooster
		case models.SubscriptionTypeSupporter:
			header = userSubscriptionSupporter
		}

		if header != nil {
			block, err := header.Block()
			if err != nil {
				return nil, err
			}
			badges = append(badges, block)
			break
		}
	}

	return badges, nil
}

func ClanSubscriptionsBadges(subscriptions []models.UserSubscription) *subscriptionHeader {
	var headers []*subscriptionHeader

	for _, subscription := range subscriptions {
		switch subscription.Type {
		case models.SubscriptionTypeProClan:
			headers = append(headers, clanSubscriptionPro)
		case models.SubscriptionTypeVerifiedClan:
			headers = append(headers, clanSubscriptionVerified)
		}
	}

	if len(headers) > 0 {
		return headers[0]
	}

	return nil
}
