package session

import (
	"errors"

	"github.com/cufee/aftermath-core/dataprep"
	"github.com/cufee/aftermath-core/dataprep/session"
	"github.com/cufee/aftermath-core/internal/core/database/models"
	"github.com/cufee/aftermath-core/internal/logic/render"
	wg "github.com/cufee/am-wg-proxy-next/types"
	"github.com/rs/zerolog/log"
)

type PlayerData struct {
	Clan          *wg.Clan
	Account       *wg.Account
	Cards         session.Cards
	Subscriptions []models.UserSubscription
}

type RenderOptions struct {
	PromoText []string
	CardStyle render.Style
}

func snapshotToCardsBlocks(player PlayerData, options RenderOptions) ([]render.Block, error) {
	if player.Account == nil {
		log.Error().Msg("player account is nil, this should not happen")
		return nil, errors.New("player account is nil")
	}

	var cards []render.Block

	var addPromoText = true
	for _, sub := range player.Subscriptions {
		switch sub.Type {
		case models.SubscriptionTypePro, models.SubscriptionTypePlus, models.SubscriptionTypeDeveloper:
			addPromoText = false
		}
		if !addPromoText {
			break
		}
	}
	// User Subscription Badge and promo text
	if badges, _ := player.userBadges(); len(badges) > 0 {
		cards = append(cards, render.NewBlocksContent(render.Style{Direction: render.DirectionHorizontal, AlignItems: render.AlignItemsCenter, Gap: 10},
			badges...,
		))
	}

	if addPromoText {
		// Users without a subscription get promo text
		if options.PromoText != nil {
			var textBlocks []render.Block
			for _, text := range options.PromoText {
				textBlocks = append(textBlocks, render.NewTextContent(render.Style{Font: &render.FontMedium, FontColor: render.TextSecondary}, text))
			}
			cards = append(cards, render.NewBlocksContent(render.Style{Direction: render.DirectionVertical, AlignItems: render.AlignItemsCenter},
				textBlocks...,
			))
		}
	}

	// Title Card
	var clanTagBlocks []render.Block
	if player.Clan != nil && player.Clan.Tag != "" {
		clanTagBlocks = append(clanTagBlocks, render.NewTextContent(render.Style{Font: &render.FontMedium, FontColor: render.TextSecondary}, player.Clan.Tag))
		if sub := player.clanSubscriptionHeader(); sub != nil {
			iconBlock, err := sub.Block()
			if err == nil {
				clanTagBlocks = append(clanTagBlocks, iconBlock)
			}
		}
	}
	cards = append(cards, newPlayerTitleCard(options.CardStyle, player.Account.Nickname, clanTagBlocks))

	styled := func(blocks []session.StatsBlock) []styledStatsBlock {
		return styleBlocks(blocks, HighlightStatsBlockStyle(options.CardStyle.BackgroundColor), DefaultStatsBlockStyle)
	}

	for _, card := range player.Cards {
		var hasCareer bool
		var hasSession bool
		for _, block := range card.Blocks {
			if block.Tag == dataprep.TagBattles {
				hasSession = block.Session.Value > 0
				hasCareer = block.Career.Value > 0
				break
			}
		}

		opts := convertOptions{true, hasCareer, true, hasCareer && hasSession}
		if card.Type == dataprep.CardTypeVehicle {
			opts = convertOptions{true, hasCareer, false, hasCareer && hasSession}
		}

		blocks, err := statsBlocksToCardBlocks(styled(card.Blocks), opts)
		if err != nil {
			return nil, err
		}
		cards = append(cards, newCardBlock(options.CardStyle, newTextLabel(card.Title), blocks))
	}
	return cards, nil
}
