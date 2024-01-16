package session

import (
	"image"
	"image/color"

	"github.com/cufee/aftermath-core/internal/core/localization"
	core "github.com/cufee/aftermath-core/internal/core/stats"
	"github.com/cufee/aftermath-core/internal/logic/cache"
	"github.com/cufee/aftermath-core/internal/logic/render"
	"github.com/cufee/aftermath-core/internal/logic/stats"
	"github.com/rs/zerolog/log"
)

type PlayerData struct {
	Snapshot *stats.Snapshot
	Vehicles []*core.ReducedVehicleStats
	Averages map[int]*core.ReducedStatsFrame
}

type RenderOptions struct {
	Locale                 localization.SupportedLanguage
	PromoText              []string
	CardStyle              render.Style
	BackgroundImage        image.Image
	UserSubscriptionHeader *SubscriptionHeader
	ClanSubscriptionHeader *SubscriptionHeader
}

func SnapshotToCardsBlocks(player PlayerData, options RenderOptions) ([]render.Block, error) {
	var cards []render.Block

	localePrinter := localization.GetPrinter(options.Locale)

	vehicleIDs := make([]int, len(player.Vehicles))
	for i, vehicle := range player.Vehicles {
		vehicleIDs[i] = vehicle.VehicleID
	}
	vehiclesGlossary, err := cache.GetGlossaryVehicles(vehicleIDs...)
	if err != nil {
		return nil, err
	}

	// Check if a user is premium
	// User Status Badge
	switch options.UserSubscriptionHeader {
	case UserSubscriptionSupporter:
		subscriptionBlock, err := options.UserSubscriptionHeader.Block()
		if err != nil {
			return nil, err
		}
		cards = append(cards, subscriptionBlock)
		fallthrough
	case nil:
		if options.PromoText != nil {
			// Promo Card
			var textBlocks []render.Block
			for _, text := range options.PromoText {
				textBlocks = append(textBlocks, render.NewTextContent(render.Style{Font: &FontMedium, FontColor: FontMediumColor}, text))
			}
			cards = append(cards, render.NewBlocksContent(render.Style{Direction: render.DirectionVertical, AlignItems: render.AlignItemsCenter},
				textBlocks...,
			))
		}
	default:
		subscriptionBlock, err := options.UserSubscriptionHeader.Block()
		if err != nil {
			return nil, err
		}
		cards = append(cards, subscriptionBlock)
	}

	// Title Card
	{
		clanSubBlock := render.NewTextContent(render.Style{Font: &FontMedium, FontColor: color.Transparent}, player.Snapshot.Account.Clan.Tag)
		if options.ClanSubscriptionHeader != nil {
			iconBlock, err := options.ClanSubscriptionHeader.Block()
			if err != nil {
				log.Warn().Err(err).Msg("failed to render clan tag") // This error is not fatal, but we should avoid trying to render the tag
			} else {
				clanSubBlock = iconBlock
			}
		}
		cards = append(cards, NewPlayerTitleCard(options.CardStyle, player.Snapshot.Account.Nickname, player.Snapshot.Account.Clan.Tag, clanSubBlock))

	}

	var totalVehicleWN8 int
	var vehicleCards []render.Block
	{
		for _, vehicle := range player.Vehicles {
			if vehicle.WN8(player.Averages[vehicle.VehicleID]) != core.InvalidValue {
				totalVehicleWN8 += vehicle.WN8(player.Averages[vehicle.VehicleID]) * vehicle.Battles
			}

			// Vehicle Cards
			blocks, err := FrameToSlimStatsBlocks(options.CardStyle, vehicle.ReducedStatsFrame, player.Averages[vehicle.VehicleID], localePrinter)
			if err != nil {
				return nil, err
			}

			vehicleInfo := vehiclesGlossary[vehicle.VehicleID]
			vehicleInfo.ID = vehicle.VehicleID
			vehicleCards = append(vehicleCards, NewCardBlock(options.CardStyle, NewVehicleLabel(vehicleInfo.Name(options.Locale), render.IntToRoman(vehicleInfo.Tier)), blocks))
		}
	}

	{
		// Regular Battles
		sessionWN8 := core.InvalidValue
		if player.Snapshot.Diff.Global.Battles > 0 {
			sessionWN8 = totalVehicleWN8 / player.Snapshot.Diff.Global.Battles
		}
		blocks, err := FrameToOverviewBlocks(options.CardStyle, player.Snapshot.Diff.Global, player.Snapshot.Selected.Global, sessionWN8, core.InvalidValue, localePrinter)
		if err != nil {
			return nil, err
		}
		cards = append(cards, NewCardBlock(options.CardStyle, NewTextLabel(localePrinter("label_overview_unrated")), blocks))
	}

	if player.Snapshot.Diff.Rating.Battles > 0 {
		// Rating Battles
		blocks, err := FrameToOverviewBlocks(options.CardStyle, player.Snapshot.Diff.Global, player.Snapshot.Selected.Global, core.InvalidValue, core.InvalidValue, localePrinter)
		if err != nil {
			return nil, err
		}
		cards = append(cards, NewCardBlock(options.CardStyle, NewTextLabel(localePrinter("label_overview_rating")), blocks))
	}

	return append(cards, vehicleCards...), nil
}
