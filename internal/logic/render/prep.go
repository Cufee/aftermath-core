package render

import (
	"errors"
	"image"

	"github.com/cufee/aftermath-core/internal/core/localization"
	core "github.com/cufee/aftermath-core/internal/core/stats"
	"github.com/cufee/aftermath-core/internal/logic/stats"
)

func FrameToStatsBlocks(session, career, averages *core.ReducedStatsFrame, locale localization.SupportedLanguage) ([]Block, error) {
	if session == nil {
		return nil, errors.New("session is nil")
	}

	// localePrinter := localization.GetPrinter(locale)

	var blocks []Block
	{
		// Battles
		values := []interface{}{session.Battles}
		if career != nil {
			values = append(values, career.Battles)
		}
		blocks = append(blocks, NewStatsBlock("battles", values...))
	}
	{
		// Avg Damage
		values := []interface{}{int(session.AvgDamage())}
		if career != nil {
			values = append(values, int(career.AvgDamage()))
		}
		blocks = append(blocks, NewStatsBlock("avg_damage", values...))
	}
	{
		// Winrate
		values := []interface{}{session.Winrate()}
		if career != nil {
			values = append(values, career.Winrate())
		}
		blocks = append(blocks, NewStatsBlock("winrate", values...))
	}
	{
		if session.WN8(averages) != core.InvalidValue {
			// WN8
			values := []interface{}{int(session.WN8(averages))}
			if career != nil {
				values = append(values, int(career.WN8(averages)))
			}
			blocks = append(blocks, NewStatsBlock("wn8", values...))
		} else {
			// Fallback to Accuracy to keep the UI consistent
			values := []interface{}{session.Accuracy()}
			if career != nil {
				values = append(values, career.Accuracy())
			}
			blocks = append(blocks, NewStatsBlock("accuracy", values...))
		}
	}

	return blocks, nil
}

func SnapshotToCardsBlocks(snapshot *stats.Snapshot, averages *core.ReducedStatsFrame, locale localization.SupportedLanguage) ([]Block, error) {
	var cards []Block

	{
		// Title Card
		cards = append(cards, NewPlayerTitleCard("NameGoesHere", "TAG")) // TODO: Pass some customization crap, stickers, etc.
	}

	{
		// Regular Battles
		blocks, err := FrameToStatsBlocks(snapshot.Diff.Global, snapshot.Selected.Global, averages, locale)
		if err != nil {
			return nil, err
		}
		cards = append(cards, NewCardBlock(NewTextLabel("overview_unrated"), blocks))
	}

	{
		// Rating Battles
		blocks, err := FrameToStatsBlocks(snapshot.Diff.Global, snapshot.Selected.Global, averages, locale)
		if err != nil {
			return nil, err
		}
		cards = append(cards, NewCardBlock(NewTextLabel("overview_rating"), blocks))
	}

	{
		for _, vehicle := range snapshot.Diff.Vehicles {
			// Vehicle Cards
			blocks, err := FrameToStatsBlocks(vehicle.ReducedStatsFrame, snapshot.Selected.Vehicles[vehicle.VehicleID].ReducedStatsFrame, averages, locale)
			if err != nil {
				return nil, err
			}
			cards = append(cards, NewCardBlock(NewVehicleLabel("Some Tank", "X"), blocks))
		}
	}

	return cards, nil
}

func RenderStatsImage(snapshot *stats.Snapshot, averages *core.ReducedStatsFrame, locale localization.SupportedLanguage) (image.Image, error) {
	cards, err := SnapshotToCardsBlocks(snapshot, averages, locale)
	if err != nil {
		return nil, err
	}

	allCards := NewBlocksContent(
		Style{
			Font:      FontLarge,
			FontColor: FontMediumColor,

			Direction:  DirectionVertical,
			AlignItems: AlignItemsCenter,
			PaddingX:   20,
			PaddingY:   20,
			Gap:        20,
			Debug:      false,
		}, cards...)

	return allCards.Render()
}
