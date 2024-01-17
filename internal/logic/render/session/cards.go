package session

import (
	"image"
	"image/color"
	"sync"

	"github.com/cufee/aftermath-core/internal/core/localization"
	"github.com/cufee/aftermath-core/internal/logic/cache"
	"github.com/cufee/aftermath-core/internal/logic/dataprep"
	"github.com/cufee/aftermath-core/internal/logic/render"
	"github.com/cufee/aftermath-core/internal/logic/stats"
	"github.com/cufee/aftermath-core/internal/logic/users"
	"github.com/rs/zerolog/log"
)

type PlayerData struct {
	Blocks        *dataprep.SessionBlocks
	Snapshot      *stats.Snapshot
	Subscriptions []users.UserSubscription
}

type RenderOptions struct {
	Locale          localization.SupportedLanguage
	PromoText       []string
	CardStyle       render.Style
	BackgroundImage image.Image
}

func snapshotToCardsBlocks(player PlayerData, options RenderOptions) ([]render.Block, error) {
	var glossarySync sync.WaitGroup
	glossaryChan := make(chan map[int]cache.VehicleInfo, 1)

	glossarySync.Add(1)
	go func() {
		defer glossarySync.Done()
		defer close(glossaryChan)
		vehicleIDs := make([]int, 0, len(player.Blocks.Vehicles))
		for _, vehicle := range player.Blocks.Vehicles {
			vehicleIDs = append(vehicleIDs, vehicle.ID)
		}
		vehiclesGlossary, err := cache.GetGlossaryVehicles(vehicleIDs...)
		if err != nil {
			// This is definitely not fatal, but will look ugly
			log.Warn().Err(err).Msg("failed to get vehicles glossary")
		}
		glossaryChan <- vehiclesGlossary
	}()

	var cards []render.Block
	localePrinter := localization.GetPrinter(options.Locale)

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
		clanSubBlock := render.NewTextContent(render.Style{Font: &FontMedium, FontColor: color.Transparent}, player.Snapshot.Account.Clan.Tag)
		if sub := player.clanSubscriptionHeader(); sub != nil {
			iconBlock, err := sub.Block()
			if err != nil {
				log.Warn().Err(err).Msg("failed to render clan tag") // This error is not fatal, but we should avoid trying to render the tag
			} else {
				clanSubBlock = iconBlock
			}
		}
		cards = append(cards, newPlayerTitleCard(options.CardStyle, player.Snapshot.Account.Nickname, player.Snapshot.Account.Clan.Tag, clanSubBlock))

	}

	// Styles
	styled := func(blocks []dataprep.StatsBlock) []styledStatsBlock {
		return styleBlocks(blocks, HighlightStatsBlockStyle(options.CardStyle.BackgroundColor), DefaultStatsBlockStyle)
	}
	vehicleHighlightBlockStyle := HighlightStatsBlockStyle(options.CardStyle.BackgroundColor)
	vehicleHighlightBlockStyle.PaddingY = 5

	{
		// Regular Battles
		blocks, err := statsBlocksToCardBlocks(styled(player.Blocks.Regular))
		if err != nil {
			return nil, err
		}
		cards = append(cards, newCardBlock(options.CardStyle, newTextLabel(localePrinter("label_overview_unrated")), blocks))
	}

	// Rating Battles
	if player.Snapshot.Diff.Rating.Battles > 0 {
		blocks, err := statsBlocksToCardBlocks(styled(player.Blocks.Rating))
		if err != nil {
			return nil, err
		}
		cards = append(cards, newCardBlock(options.CardStyle, newTextLabel(localePrinter("label_overview_rating")), blocks))
	}

	{
		glossarySync.Wait()
		glossary := <-glossaryChan
		// Vehicle Cards
		for _, vehicle := range player.Blocks.Vehicles {
			styledBlocks := styleBlocks(vehicle.Blocks, vehicleHighlightBlockStyle, DefaultStatsBlockStyle)
			blocks, err := statsBlocksToCardBlocks(styledBlocks, convertOptions{showSessionStats: true, showCareerStats: true})
			if err != nil {
				return nil, err
			}
			vehicleInfo := glossary[vehicle.ID]
			vehicleInfo.ID = vehicle.ID
			cards = append(cards, newCardBlock(options.CardStyle, newVehicleLabel(vehicleInfo.Name(options.Locale), render.IntToRoman(vehicleInfo.Tier)), blocks))
		}
	}

	return cards, nil
}
