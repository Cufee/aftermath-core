package session

import (
	"errors"
	"strings"
	"time"

	"github.com/cufee/aftermath-core/dataprep"
	"github.com/cufee/aftermath-core/dataprep/session"
	"github.com/cufee/aftermath-core/internal/core/database/models"
	helpers "github.com/cufee/aftermath-core/internal/core/utils"
	"github.com/cufee/aftermath-core/internal/logic/render"
	"github.com/cufee/aftermath-core/internal/logic/render/badges"
	"github.com/cufee/aftermath-core/internal/logic/render/shared"
	"github.com/cufee/aftermath-core/internal/logic/stats/sessions"
	"github.com/cufee/aftermath-core/utils"
	wg "github.com/cufee/am-wg-proxy-next/types"
	"github.com/rs/zerolog/log"
)

type PlayerData struct {
	Clan    wg.Clan
	Account wg.Account
	Session sessions.Snapshot

	Subscriptions []models.UserSubscription
	Cards         session.Cards
}

type RenderOptions struct {
	PromoText []string
}

func snapshotToCardsBlocks(player PlayerData, options RenderOptions) ([]render.Block, error) {
	if player.Account.ID <= 0 {
		log.Error().Msg("player account is not set, this should not happen")
		return nil, errors.New("player account is nil")
	}
	if len(player.Cards) == 0 {
		log.Error().Msg("player cards slice is 0 length, this should not happen")
		return nil, errors.New("no cards provided")
	}

	// Calculate minimal card width to fit all the content
	var cardWidth float64
	cardBlockSizes := make(map[int]float64)
	{
		{
			titleStyle := shared.DefaultPlayerTitleStyle(titleCardStyle(cardWidth))
			clanSize := render.MeasureString(player.Clan.Tag, *titleStyle.ClanTag.Font)
			nameSize := render.MeasureString(player.Account.Nickname, *titleStyle.Nickname.Font)
			cardWidth = helpers.Max(cardWidth, titleStyle.TotalPaddingAndGaps()+nameSize.TotalWidth+clanSize.TotalWidth*2)
		}
		{
			for _, text := range options.PromoText {
				size := render.MeasureString(text, *promoTextStyle.Font)
				cardWidth = helpers.Max(size.TotalWidth, cardWidth)
			}
		}

		{
			for _, card := range player.Cards {
				var blockSizeMax float64

				for index, block := range card.Blocks {
					var blockWidth float64
					{
						size := render.MeasureString(block.Session.String, *defaultBlockStyle.session.Font)
						blockWidth = helpers.Max(size.TotalWidth+defaultBlockStyle.session.PaddingX*2+defaultBlockStyle.session.Gap, blockWidth)
					}
					{
						size := render.MeasureString(block.Career.String, *defaultBlockStyle.career.Font)
						blockWidth = helpers.Max(size.TotalWidth+defaultBlockStyle.career.PaddingX*2+defaultBlockStyle.career.Gap, blockWidth)
					}
					if card.Type != dataprep.CardTypeVehicle {
						size := render.MeasureString(block.Label, *defaultBlockStyle.label.Font)
						blockWidth = helpers.Max(size.TotalWidth+defaultBlockStyle.label.PaddingX*2+defaultBlockStyle.label.Gap, blockWidth)
					}

					totalBlockWidth := blockWidth + float64(iconSize)
					if index == 0 {
						totalBlockWidth += highlightStatsBlockStyle(0).PaddingX * 2
					}

					cardBlockSizes[index] = helpers.Max(cardBlockSizes[index], totalBlockWidth)
				}

				if card.Type == dataprep.CardTypeRatingVehicle {
					vehicleNameSize := render.MeasureString(card.Title, *ratingVehicleTitleStyle.Font)
					cardWidth = helpers.Max(cardWidth, (defaultCardStyle(0).PaddingX*4)+(defaultCardStyle(0).Gap*float64(len(card.Blocks)-1))+vehicleNameSize.TotalWidth+blockSizeMax)
				}

			}

			// Find the minimum required width to fix card content for the largest card
			var totalContentSize float64
			for _, size := range cardBlockSizes {
				totalContentSize += size
			}

			// why padding is *4? did not care to debug, but smells like a bug with how card width vs content width is calculated
			cardWidth = helpers.Max(cardWidth, (defaultCardStyle(0).PaddingX*4)+(defaultCardStyle(0).Gap*float64(len(cardBlockSizes)-1))+totalContentSize)
		}
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
	if addPromoText && options.PromoText != nil {
		// Users without a subscription get promo text
		var textBlocks []render.Block
		for _, text := range options.PromoText {
			textBlocks = append(textBlocks, render.NewTextContent(promoTextStyle, text))
		}
		cards = append(cards, render.NewBlocksContent(render.Style{
			Direction:  render.DirectionVertical,
			AlignItems: render.AlignItemsCenter,
		},
			textBlocks...,
		))
	}
	if badges, _ := badges.SubscriptionsBadges(player.Subscriptions); len(badges) > 0 {
		cards = append(cards, render.NewBlocksContent(render.Style{Direction: render.DirectionHorizontal, AlignItems: render.AlignItemsCenter, Gap: 10},
			badges...,
		))
	}

	cards = append(cards, shared.NewPlayerTitleCard(shared.DefaultPlayerTitleStyle(titleCardStyle(cardWidth)), player.Account.Nickname, player.Clan.Tag, player.Subscriptions))

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

		opts := convertOptions{true, hasCareer, true, hasCareer && hasSession, 0}
		if card.Type == dataprep.CardTypeVehicle {
			opts = convertOptions{true, hasCareer, false, hasCareer && hasSession, 0}
		}

		card, err := newVehicleCard(defaultCardStyle(cardWidth), card, cardBlockSizes, opts)
		if err != nil {
			return nil, err
		}
		cards = append(cards, card)
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
	if player.Session.Selected.LastBattleTime > 0 {
		sessionTo := time.Unix(int64(player.Session.Live.LastBattleTime), 0).Format("January 2")
		sessionFrom := time.Unix(int64(player.Session.Selected.LastBattleTime), 0).Format("January 2")
		if sessionFrom == sessionTo {
			footer = append(footer, sessionTo)
		} else {
			footer = append(footer, sessionFrom+" - "+sessionTo)
		}
	}

	if len(footer) > 0 {
		cards = append(cards, shared.NewFooterCard(strings.Join(footer, " • ")))
	}

	return cards, nil
}
