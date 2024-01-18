package session

import (
	"errors"
	"image"
	"image/color"

	"github.com/cufee/aftermath-core/internal/logic/dataprep"
	"github.com/cufee/aftermath-core/internal/logic/render"
	"github.com/cufee/aftermath-core/internal/logic/users"
	wg "github.com/cufee/am-wg-proxy-next/types"
	"github.com/rs/zerolog/log"
)

type PlayerData struct {
	Clan          *wg.Clan
	Account       *wg.Account
	Cards         dataprep.SessionCards
	Subscriptions []users.UserSubscription
}

type RenderOptions struct {
	PromoText       []string
	CardStyle       render.Style
	BackgroundImage image.Image
}

func snapshotToCardsBlocks(player PlayerData, options RenderOptions) ([]render.Block, error) {
	if player.Account == nil {
		log.Error().Msg("player account is nil, this should not happen")
		return nil, errors.New("player account is nil")
	}

	var cards []render.Block

	// User Status Badge
	switch sub := player.userSubscriptionHeader(); sub {
	case userSubscriptionSupporter:
		// Supporters get a badge and a promo text
		subscriptionBlock, err := sub.Block()
		if err != nil {
			return nil, err
		}
		cards = append(cards, subscriptionBlock)
		fallthrough
	case nil:
		// Users without a subscription and supporters get a promo text
		if options.PromoText != nil {
			var textBlocks []render.Block
			for _, text := range options.PromoText {
				textBlocks = append(textBlocks, render.NewTextContent(render.Style{Font: &FontMedium, FontColor: FontMediumColor}, text))
			}
			cards = append(cards, render.NewBlocksContent(render.Style{Direction: render.DirectionVertical, AlignItems: render.AlignItemsCenter},
				textBlocks...,
			))
		}
	default:
		// All other subscriptions get a badge
		subscriptionBlock, err := sub.Block()
		if err != nil {
			return nil, err
		}
		cards = append(cards, subscriptionBlock)
	}

	// Title Card
	{
		var clanTag string
		if player.Clan != nil {
			clanTag = player.Clan.Tag
		}

		clanSubBlock := render.NewTextContent(render.Style{Font: &FontMedium, FontColor: color.Transparent}, clanTag)
		if sub := player.clanSubscriptionHeader(); sub != nil {
			iconBlock, err := sub.Block()
			if err != nil {
				log.Warn().Err(err).Msg("failed to render clan tag") // This error is not fatal, but we should avoid trying to render the tag
			} else {
				clanSubBlock = iconBlock
			}
		}
		cards = append(cards, newPlayerTitleCard(options.CardStyle, player.Account.Nickname, clanTag, clanSubBlock))
	}

	// Styles
	styled := func(blocks []dataprep.StatsBlock) []styledStatsBlock {
		return styleBlocks(blocks, HighlightStatsBlockStyle(options.CardStyle.BackgroundColor), DefaultStatsBlockStyle)
	}
	vehicleHighlightBlockStyle := HighlightStatsBlockStyle(options.CardStyle.BackgroundColor)
	vehicleHighlightBlockStyle.PaddingY = 5

	for _, card := range player.Cards {
		opts := convertOptions{true, true, true}
		if card.Type == dataprep.CardTypeVehicle {
			opts = convertOptions{true, true, false}
		}

		blocks, err := statsBlocksToCardBlocks(styled(card.Blocks), opts)
		if err != nil {
			return nil, err
		}
		cards = append(cards, newCardBlock(options.CardStyle, newTextLabel(card.Title), blocks))
	}
	return cards, nil
}
