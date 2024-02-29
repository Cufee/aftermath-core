package session

import (
	"errors"
	"strings"
	"time"

	"github.com/cufee/aftermath-core/dataprep"
	"github.com/cufee/aftermath-core/dataprep/session"
	"github.com/cufee/aftermath-core/internal/core/database/models"
	"github.com/cufee/aftermath-core/internal/logic/render"
	"github.com/cufee/aftermath-core/internal/logic/render/badges"
	"github.com/cufee/aftermath-core/internal/logic/stats/sessions"
	"github.com/cufee/aftermath-core/utils"
	wg "github.com/cufee/am-wg-proxy-next/types"
	"github.com/rs/zerolog/log"
)

type PlayerData struct {
	Clan    *wg.Clan
	Account *wg.Account
	Session *sessions.Snapshot

	Subscriptions []models.UserSubscription
	Cards         session.Cards
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
	if len(player.Cards) == 0 {
		log.Error().Msg("player cards slice is 0 length, this should not happen")
		return nil, errors.New("no cards provided")
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
	if badges, _ := badges.SubscriptionsBadges(player.Subscriptions); len(badges) > 0 {
		cards = append(cards, render.NewBlocksContent(render.Style{Direction: render.DirectionHorizontal, AlignItems: render.AlignItemsCenter, Gap: 10},
			badges...,
		))
	}

	if addPromoText && options.PromoText != nil {
		// Users without a subscription get promo text
		var textBlocks []render.Block
		for _, text := range options.PromoText {
			textBlocks = append(textBlocks, render.NewTextContent(render.Style{Font: &render.FontMedium, FontColor: render.TextPrimary}, text))
		}
		cards = append(cards, render.NewBlocksContent(render.Style{
			Direction:  render.DirectionVertical,
			AlignItems: render.AlignItemsCenter,
		},
			textBlocks...,
		))
	}

	// Title Card
	var clanTagBlocks []render.Block
	if player.Clan != nil && player.Clan.Tag != "" {
		clanTagBlocks = append(clanTagBlocks, render.NewTextContent(render.Style{Font: &render.FontMedium, FontColor: render.TextSecondary}, player.Clan.Tag))
		if sub := badges.ClanSubscriptionsBadges(player.Subscriptions); sub != nil {
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

	var footer []string
	switch strings.ToLower(utils.RealmFromAccountID(player.Account.ID)) {
	case "na":
		footer = append(footer, "North America")
	case "eu":
		footer = append(footer, "Europe")
	case "as":
		footer = append(footer, "Asia")
	}
	if player.Session != nil && player.Session.Selected.LastBattleTime > 0 {
		sessionTo := time.Unix(int64(player.Session.Live.LastBattleTime), 0).Format("January 2")
		sessionFrom := time.Unix(int64(player.Session.Selected.LastBattleTime), 0).Format("January 2")
		if sessionFrom == sessionTo {
			footer = append(footer, sessionTo)
		} else {
			footer = append(footer, sessionFrom+" - "+sessionTo)
		}
	}

	if len(footer) > 0 {
		cards = append(cards, render.NewTextContent(render.Style{Font: &render.FontSmall, FontColor: render.TextAlt}, strings.Join(footer, " • ")))
	}

	return cards, nil
}
