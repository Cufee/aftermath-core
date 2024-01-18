package session

import (
	"errors"
	"image"

	"github.com/cufee/aftermath-core/internal/core/database/models"
	"github.com/cufee/aftermath-core/internal/logic/dataprep"
	"github.com/cufee/aftermath-core/internal/logic/render"
	wg "github.com/cufee/am-wg-proxy-next/types"
	"github.com/rs/zerolog/log"
)

type PlayerData struct {
	Clan          *wg.Clan
	Account       *wg.Account
	Cards         dataprep.SessionCards
	Subscriptions []models.UserSubscription
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

	// User Subscription Badge and promo text
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
	var clanTagBlocks []render.Block
	if player.Clan != nil {
		clanTagBlocks = append(clanTagBlocks, render.NewTextContent(render.Style{Font: &FontMedium, FontColor: FontMediumColor}, player.Clan.Tag))
		if sub := player.clanSubscriptionHeader(); sub != nil {
			iconBlock, err := sub.Block()
			if err == nil {
				clanTagBlocks = append(clanTagBlocks, iconBlock)
			}
		}
	}
	cards = append(cards, newPlayerTitleCard(options.CardStyle, player.Account.Nickname, clanTagBlocks))

	styled := func(blocks []dataprep.StatsBlock) []styledStatsBlock {
		return styleBlocks(blocks, HighlightStatsBlockStyle(options.CardStyle.BackgroundColor), DefaultStatsBlockStyle)
	}

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
