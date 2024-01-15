package render

import (
	"errors"
	"image"

	"github.com/cufee/aftermath-core/internal/core/localization"
	core "github.com/cufee/aftermath-core/internal/core/stats"
	"github.com/cufee/aftermath-core/internal/logic/cache"
	"github.com/cufee/aftermath-core/internal/logic/render/assets"
	"github.com/cufee/aftermath-core/internal/logic/stats"
	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
)

func FrameToOverviewBlocks(session, career *core.ReducedStatsFrame, sessionWN8, careerWN8 int, localePrinter localization.LocalePrinter) ([]Block, error) {
	if session == nil {
		return nil, errors.New("session is nil")
	}

	var blocks []Block
	{
		// Battles
		values := []interface{}{session.Battles}
		if career != nil {
			values = append(values, career.Battles)
		}
		battlesBlock := NewStatsBlock("", values...)
		// Add some special styling to the block
		battlesBlock.Style.PaddingY = 10
		battlesBlock.Style.BorderRadius = 10
		battlesBlock.Style.BackgroundColor = HighlightCardColor
		blocks = append(blocks, battlesBlock)
	}
	{
		// Avg Damage
		values := []any{int(session.AvgDamage())}
		if career != nil {
			values = append(values, int(career.AvgDamage()))
		}
		blocks = append(blocks, NewStatsBlock(localePrinter("label_avg_damage"), values...))
	}
	{
		// Winrate
		values := []any{session.Winrate()}
		if career != nil {
			values = append(values, career.Winrate())
		}
		blocks = append(blocks, NewStatsBlock(localePrinter("label_winrate"), values...))
	}
	{
		if sessionWN8 != core.InvalidValue {
			// WN8
			values := []any{sessionWN8}
			if careerWN8 != core.InvalidValue {
				values = append(values, careerWN8)
			}
			blocks = append(blocks, NewStatsBlock(localePrinter("label_wn8"), values...))
		} else {
			// Fallback to Accuracy to keep the UI consistent
			values := []any{session.Accuracy()}
			if career != nil {
				values = append(values, career.Accuracy())
			}
			blocks = append(blocks, NewStatsBlock(localePrinter("label_accuracy"), values...))
		}
	}

	return blocks, nil
}

func FrameToStatsBlocks(session, career, averages *core.ReducedStatsFrame, localePrinter localization.LocalePrinter) ([]Block, error) {
	if session == nil {
		return nil, errors.New("session is nil")
	}

	var blocks []Block
	{
		// Battles
		values := []interface{}{session.Battles}
		if career != nil {
			values = append(values, career.Battles)
		}
		battlesBlock := NewStatsBlock("", values...)
		// Add some special styling to the block
		battlesBlock.Style.PaddingY = 10
		battlesBlock.Style.BorderRadius = 10
		battlesBlock.Style.BackgroundColor = HighlightCardColor
		blocks = append(blocks, battlesBlock)
	}
	{
		// Avg Damage
		values := []interface{}{int(session.AvgDamage())}
		if career != nil {
			values = append(values, int(career.AvgDamage()))
		}
		blocks = append(blocks, NewStatsBlock(localePrinter("label_avg_damage"), values...))
	}
	{
		// Winrate
		values := []interface{}{session.Winrate()}
		if career != nil {
			values = append(values, career.Winrate())
		}
		blocks = append(blocks, NewStatsBlock(localePrinter("label_winrate"), values...))
	}
	{
		if session.WN8(averages) != core.InvalidValue {
			// WN8
			values := []interface{}{int(session.WN8(averages))}
			if career != nil {
				values = append(values, int(career.WN8(averages)))
			}
			blocks = append(blocks, NewStatsBlock(localePrinter("label_wn8"), values...))
		} else {
			// Fallback to Accuracy to keep the UI consistent
			values := []interface{}{session.Accuracy()}
			if career != nil {
				values = append(values, career.Accuracy())
			}
			blocks = append(blocks, NewStatsBlock(localePrinter("label_accuracy"), values...))
		}
	}

	return blocks, nil
}

func FrameToSlimStatsBlocks(session, averages *core.ReducedStatsFrame, localePrinter localization.LocalePrinter) ([]Block, error) {
	if session == nil {
		return nil, errors.New("session is nil")
	}

	var blocks []Block
	{
		// Battles
		battlesBlock := NewStatsBlock("", session.Battles)
		// Add some special styling to the block
		battlesBlock.Style.PaddingY = 10
		battlesBlock.Style.BorderRadius = 10
		battlesBlock.Style.BackgroundColor = HighlightCardColor
		blocks = append(blocks, battlesBlock)
	}
	{
		// Avg Damage
		blocks = append(blocks, NewStatsBlock(localePrinter("label_avg_damage"), int(session.AvgDamage())))
	}
	{
		// Winrate
		blocks = append(blocks, NewStatsBlock(localePrinter("label_winrate"), session.Winrate()))
	}
	{
		if session.WN8(averages) != core.InvalidValue {
			// WN8
			blocks = append(blocks, NewStatsBlock(localePrinter("label_wn8"), session.WN8(averages)))
		} else {
			// Fallback to Accuracy to keep the UI consistent
			blocks = append(blocks, NewStatsBlock(localePrinter("label_accuracy"), session.Accuracy()))
		}
	}

	return blocks, nil
}

func SnapshotToCardsBlocks(snapshot *stats.Snapshot, vehicles []*core.ReducedVehicleStats, averages map[int]*core.ReducedStatsFrame, locale localization.SupportedLanguage) ([]Block, error) {
	var cards []Block

	localePrinter := localization.GetPrinter(locale)

	vehicleIDs := make([]int, len(vehicles))
	for i, vehicle := range vehicles {
		vehicleIDs[i] = vehicle.VehicleID
	}
	vehiclesGlossary, err := cache.GetGlossaryVehicles(vehicleIDs...)
	if err != nil {
		return nil, err
	}

	{
		// Promo Card
		cards = append(cards, NewBlocksContent(Style{Direction: DirectionVertical, AlignItems: AlignItemsCenter},
			NewTextContent("Aftermath is back!", Style{Font: FontMedium, FontColor: FontTranslucentColor}),
			NewTextContent("amth.one/join", Style{Font: FontMedium, FontColor: FontTranslucentColor}),
		))
	}

	// Title Card
	{
		// TODO: Pass some customization crap, stickers, etc.
		cards = append(cards, NewPlayerTitleCard(snapshot.Account.Nickname, snapshot.Account.Clan.Tag))
	}

	var totalVehicleWN8 int
	var vehicleCards []Block
	{
		for _, vehicle := range vehicles {
			if vehicle.WN8(averages[vehicle.VehicleID]) != core.InvalidValue {
				totalVehicleWN8 += vehicle.WN8(averages[vehicle.VehicleID]) * vehicle.Battles
			}

			// Vehicle Cards
			blocks, err := FrameToSlimStatsBlocks(vehicle.ReducedStatsFrame, averages[vehicle.VehicleID], localePrinter)
			if err != nil {
				return nil, err
			}

			vehicleInfo := vehiclesGlossary[vehicle.VehicleID]
			vehicleInfo.ID = vehicle.VehicleID
			vehicleCards = append(vehicleCards, NewCardBlock(NewVehicleLabel(vehicleInfo.Name(locale), intToRoman(vehicleInfo.Tier)), blocks))
		}
	}

	{
		// Regular Battles
		sessionWN8 := core.InvalidValue
		if snapshot.Diff.Global.Battles > 0 {
			sessionWN8 = totalVehicleWN8 / snapshot.Diff.Global.Battles
		}
		blocks, err := FrameToOverviewBlocks(snapshot.Diff.Global, snapshot.Selected.Global, sessionWN8, core.InvalidValue, localePrinter)
		if err != nil {
			return nil, err
		}
		cards = append(cards, NewCardBlock(NewTextLabel(localePrinter("label_overview_unrated")), blocks))
	}

	if snapshot.Diff.Rating.Battles > 0 {
		// Rating Battles
		blocks, err := FrameToOverviewBlocks(snapshot.Diff.Global, snapshot.Selected.Global, core.InvalidValue, core.InvalidValue, localePrinter)
		if err != nil {
			return nil, err
		}
		cards = append(cards, NewCardBlock(NewTextLabel(localePrinter("label_overview_rating")), blocks))
	}

	return append(cards, vehicleCards...), nil
}

func RenderStatsImage(snapshot *stats.Snapshot, vehicles []*core.ReducedVehicleStats, averages map[int]*core.ReducedStatsFrame, locale localization.SupportedLanguage) (image.Image, error) {
	cards, err := SnapshotToCardsBlocks(snapshot, vehicles, averages, locale)
	if err != nil {
		return nil, err
	}

	// TODO: Some text outside of a card, like session date, can be added here

	allCards := NewBlocksContent(
		Style{
			Direction:  DirectionVertical,
			AlignItems: AlignItemsCenter,
			PaddingX:   20,
			PaddingY:   20,
			Gap:        10,
			// Debug:      true,
		}, cards...)

	// TODO: Custom images from users and default as a fallback image
	bgImage, _ := assets.GetImage("images/backgrounds/default")
	// if !ok {
	// 	// This is always ok because for now
	// }

	cardsImage, err := allCards.Render()
	if err != nil {
		return nil, err
	}

	frameCtx := gg.NewContextForImage(cardsImage)
	// Resize the background image to fit the cards
	bgImage = imaging.Fill(bgImage, frameCtx.Width(), frameCtx.Height(), imaging.Center, imaging.NearestNeighbor)
	bgImage = imaging.Blur(bgImage, 10.0)
	frameCtx.DrawImage(bgImage, 0, 0)
	frameCtx.DrawImage(cardsImage, 0, 0)

	return frameCtx.Image(), nil
}
