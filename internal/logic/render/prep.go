package render

import (
	"errors"
	"image"

	"github.com/cufee/aftermath-core/internal/core/localization"
	core "github.com/cufee/aftermath-core/internal/core/stats"
	"github.com/cufee/aftermath-core/internal/logic/render/assets"
	"github.com/cufee/aftermath-core/internal/logic/stats"
	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
)

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

func SnapshotToCardsBlocks(snapshot *stats.Snapshot, averages map[int]core.ReducedStatsFrame, locale localization.SupportedLanguage) ([]Block, error) {
	var cards []Block

	localePrinter := localization.GetPrinter(locale)

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

	{
		// Regular Battles
		blocks, err := FrameToStatsBlocks(snapshot.Diff.Global, snapshot.Selected.Global, nil, localePrinter)
		if err != nil {
			return nil, err
		}
		cards = append(cards, NewCardBlock(NewTextLabel(localePrinter("label_overview_unrated")), blocks))
	}

	if snapshot.Diff.Rating.Battles > 0 {
		// Rating Battles
		blocks, err := FrameToStatsBlocks(snapshot.Diff.Global, snapshot.Selected.Global, nil, localePrinter)
		if err != nil {
			return nil, err
		}
		cards = append(cards, NewCardBlock(NewTextLabel(localePrinter("label_overview_rating")), blocks))
	}

	{
		for _, vehicle := range snapshot.Diff.Vehicles {
			// Vehicle Cards
			tankAverages := averages[vehicle.VehicleID]
			blocks, err := FrameToSlimStatsBlocks(vehicle.ReducedStatsFrame, &tankAverages, localePrinter)
			if err != nil {
				return nil, err
			}
			cards = append(cards, NewCardBlock(NewVehicleLabel("Some Tank", "X"), blocks))
		}
	}

	return cards, nil
}

func RenderStatsImage(snapshot *stats.Snapshot, averages map[int]core.ReducedStatsFrame, locale localization.SupportedLanguage) (image.Image, error) {
	cards, err := SnapshotToCardsBlocks(snapshot, averages, locale)
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
	bgImage, _ := assets.GetBackground("default")
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
