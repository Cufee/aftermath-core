package period

import (
	"errors"
	"strings"

	"github.com/cufee/aftermath-core/internal/logic/render"

	"github.com/cufee/aftermath-core/utils"
	"github.com/rs/zerolog/log"
)

func snapshotToCardsBlocks(player PlayerData, options RenderOptions) ([]render.Block, error) {
	if len(player.Cards.Overview.Blocks) == 0 && len(player.Cards.Highlights) == 0 {
		log.Error().Msg("player cards slice is 0 length, this should not happen")
		return nil, errors.New("no cards provided")
	}

	var cards []render.Block

	// var addPromoText = true
	// for _, sub := range player.Subscriptions {
	// 	switch sub.Type {
	// 	case models.SubscriptionTypePro, models.SubscriptionTypePlus, models.SubscriptionTypeDeveloper:
	// 		addPromoText = false
	// 	}
	// 	if !addPromoText {
	// 		break
	// 	}
	// }
	// // User Subscription Badge and promo text
	// if badges, _ := player.userBadges(); len(badges) > 0 {
	// 	cards = append(cards, render.NewBlocksContent(render.Style{Direction: render.DirectionHorizontal, AlignItems: render.AlignItemsCenter, Gap: 10},
	// 		badges...,
	// 	))
	// }

	// if addPromoText && options.PromoText != nil {
	// 	// Users without a subscription get promo text
	// 	var textBlocks []render.Block
	// 	for _, text := range options.PromoText {
	// 		textBlocks = append(textBlocks, render.NewTextContent(render.Style{Font: &render.FontMedium, FontColor: render.TextPrimary}, text))
	// 	}
	// 	cards = append(cards, render.NewBlocksContent(render.Style{
	// 		Direction:  render.DirectionVertical,
	// 		AlignItems: render.AlignItemsCenter,
	// 	},
	// 		textBlocks...,
	// 	))
	// }

	// Title Card
	// var clanTagBlocks []render.Block
	// if player.Clan != nil && player.Clan.Tag != "" {
	// 	clanTagBlocks = append(clanTagBlocks, render.NewTextContent(render.Style{Font: &render.FontMedium, FontColor: render.TextSecondary}, player.Clan.Tag))
	// 	if sub := player.clanSubscriptionHeader(); sub != nil {
	// 		iconBlock, err := sub.Block()
	// 		if err == nil {
	// 			clanTagBlocks = append(clanTagBlocks, iconBlock)
	// 		}
	// 	}
	// }
	// cards = append(cards, newPlayerTitleCard(options.CardStyle, player.Account.Nickname, clanTagBlocks))

	// styled := func(blocks []session.StatsBlock) []styledStatsBlock {
	// 	return styleBlocks(blocks, HighlightStatsBlockStyle(options.CardStyle.BackgroundColor), DefaultStatsBlockStyle)
	// }

	var overviewCardBlocks []render.Block
	for _, row := range player.Cards.Overview.Blocks {
		rowBlock, err := statsBlocksToRowBlock(row)
		if err != nil {
			return nil, err
		}
		overviewCardBlocks = append(overviewCardBlocks, rowBlock)
	}
	cards = append(cards, render.NewBlocksContent(render.Style{Direction: render.DirectionVertical}, overviewCardBlocks...))

	var footer []string
	switch strings.ToLower(utils.RealmFromAccountID(player.Stats.Account.ID)) {
	case "na":
		footer = append(footer, "North America")
	case "eu":
		footer = append(footer, "Europe")
	case "as":
		footer = append(footer, "Asia")
	}

	sessionTo := player.Stats.End.Format("January 2")
	sessionFrom := player.Stats.Start.Format("January 2")
	if sessionFrom == sessionTo {
		footer = append(footer, sessionTo)
	} else {
		footer = append(footer, sessionFrom+" - "+sessionTo)
	}

	cards = append(cards, render.NewTextContent(render.Style{Font: &render.FontSmall, FontColor: render.TextAlt}, strings.Join(footer, " • ")))

	return cards, nil
}
