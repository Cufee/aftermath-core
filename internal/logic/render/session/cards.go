package session

import (
	"errors"
	"slices"
	"sync"

	"github.com/cufee/aftermath-core/dataprep"
	"github.com/cufee/aftermath-core/dataprep/session"
	"github.com/cufee/aftermath-core/internal/core/database/models"
	"github.com/cufee/aftermath-core/internal/core/utils"
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

type renderedCard struct {
	index int
	card  render.Block
}

func snapshotToCardsBlocks(player PlayerData, options RenderOptions) ([]render.Block, error) {
	if player.Account == nil {
		log.Error().Msg("player account is nil, this should not happen")
		return nil, errors.New("player account is nil")
	}

	var cardsLock sync.Mutex
	var cardsSync sync.WaitGroup
	cards := make([]utils.DataWithError[renderedCard], 0, 3+len(player.Cards))

	cardsSync.Add(1)
	go func() {
		defer cardsSync.Done()
		// User Subscription Badge and promo text
		switch sub := player.userSubscriptionHeader(); sub {
		case userSubscriptionSupporter, subscriptionServerBooster:
			// Supporters and Boosters get a badge and promo text
			subscriptionBlock, err := sub.Block()
			if err != nil {
				cardsLock.Lock()
				cards = append(cards, utils.DataWithError[renderedCard]{Err: err})
				cardsLock.Unlock()
				return
			}
			cardsLock.Lock()
			cards = append(cards, utils.DataWithError[renderedCard]{Data: renderedCard{index: 0, card: subscriptionBlock}})
			cardsLock.Unlock()
			fallthrough
		case nil:
			// Users without a subscription get promo text
			if options.PromoText != nil {
				var textBlocks []render.Block
				for _, text := range options.PromoText {
					textBlocks = append(textBlocks, render.NewTextContent(render.Style{Font: &render.FontMedium, FontColor: render.TextSecondary}, text))
				}
				promoCard := render.NewBlocksContent(render.Style{Direction: render.DirectionVertical, AlignItems: render.AlignItemsCenter}, textBlocks...)
				cardsLock.Lock()
				cards = append(cards, utils.DataWithError[renderedCard]{Data: renderedCard{index: 1, card: promoCard}})
				cardsLock.Unlock()
			}
		default:
			// All other subscriptions get a badge
			subscriptionBlock, err := sub.Block()
			if err != nil {
				cardsLock.Lock()
				cards = append(cards, utils.DataWithError[renderedCard]{Err: err})
				cardsLock.Unlock()
				return
			}
			cardsLock.Lock()
			cards = append(cards, utils.DataWithError[renderedCard]{Data: renderedCard{index: 0, card: subscriptionBlock}})
			cardsLock.Unlock()
		}
	}()

	cardsSync.Add(1)
	go func() {
		defer cardsSync.Done()
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
		cardsLock.Lock()
		cards = append(cards, utils.DataWithError[renderedCard]{Data: renderedCard{index: 5, card: newPlayerTitleCard(options.CardStyle, player.Account.Nickname, clanTagBlocks)}})
		cardsLock.Unlock()
	}()

	styled := func(blocks []session.StatsBlock) []styledStatsBlock {
		return styleBlocks(blocks, HighlightStatsBlockStyle(options.CardStyle.BackgroundColor), DefaultStatsBlockStyle)
	}

	for i, card := range player.Cards {
		cardsSync.Add(1)
		go func(i int, card dataprep.StatsCard[session.StatsBlock, string]) {
			defer cardsSync.Done()
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
				cardsLock.Lock()
				cards = append(cards, utils.DataWithError[renderedCard]{Err: err})
				cardsLock.Unlock()
				return
			}
			cardsLock.Lock()
			cards = append(cards, utils.DataWithError[renderedCard]{Data: renderedCard{index: 10 + i, card: newCardBlock(options.CardStyle, newTextLabel(card.Title), blocks)}})
			cardsLock.Unlock()
		}(i, card)
	}

	cardsSync.Wait()
	slices.SortFunc(cards, func(a, b utils.DataWithError[renderedCard]) int {
		return a.Data.index - b.Data.index
	})

	var cardsSlice []render.Block
	for _, card := range cards {
		if card.Err != nil {
			return nil, card.Err
		}
		cardsSlice = append(cardsSlice, card.Data.card)
	}

	return cardsSlice, nil
}
