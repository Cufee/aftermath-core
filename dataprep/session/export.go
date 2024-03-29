package session

import (
	"errors"
	"fmt"

	"github.com/cufee/aftermath-core/dataprep"
	"github.com/cufee/aftermath-core/internal/core/database/models"
	core "github.com/cufee/aftermath-core/internal/core/stats"
	"github.com/cufee/aftermath-core/internal/core/utils"
	"golang.org/x/text/language"
)

type ExportInput struct {
	CareerStats  core.SessionSnapshot
	SessionStats core.SessionSnapshot

	SessionRatingVehicles  []core.ReducedVehicleStats
	SessionUnratedVehicles []core.ReducedVehicleStats

	VehicleGlossary       map[int]models.Vehicle
	GlobalVehicleAverages map[int]core.ReducedStatsFrame
}

type ExportOptions struct {
	IncludeRatingVehicles bool
	Blocks                []dataprep.Tag
	Locale                language.Tag
	LocalePrinter         func(string) string
}

func SnapshotToSession(input ExportInput, options ExportOptions) (Cards, error) {
	if input.SessionStats.AccountID == 0 || input.CareerStats.AccountID == 0 {
		return Cards{}, errors.New("session or career stats have blank accountID")
	}
	if options.LocalePrinter == nil {
		options.LocalePrinter = func(s string) string { return s }
	}
	if input.VehicleGlossary == nil {
		input.VehicleGlossary = make(map[int]models.Vehicle)
	}

	var cards Cards

	var allBattles = input.SessionStats.Rating.Battles + input.SessionStats.Global.Battles

	// Rating battles
	if allBattles == 0 || input.SessionStats.Rating.Battles > 0 {
		var ratingBlocks []StatsBlock
		for _, preset := range options.Blocks {
			if preset == dataprep.TagWN8 {
				// Rating battles have no WN8, so we use Accuracy instead of drawing a blank
				preset = dataprep.TagRankedRating
			}
			ratingBlock, err := presetToBlock(preset, options.LocalePrinter, input.SessionStats.Rating, input.CareerStats.Rating)
			if err != nil {
				return cards, fmt.Errorf("failed to generate a rating stats from preset: %w", err)
			}
			ratingBlocks = append(ratingBlocks, ratingBlock)
		}
		cards.Rating = append(cards.Rating, Card{
			Title:  options.LocalePrinter("label_overview_rating"),
			Blocks: ratingBlocks,
			Type:   dataprep.CardTypeOverview,
		})
	}

	// Rating Vehicles
	if input.SessionStats.Rating.Battles > 0 && options.IncludeRatingVehicles {
		for _, vehicle := range input.SessionRatingVehicles {
			// Wargaming does not provide any stats whatsoever on vehicle stats from Rating Battles
			// we just calculate WN8 based on the entire session and make this the only block
			block, err := presetToBlock(dataprep.TagWN8, options.LocalePrinter, input.SessionStats.Rating, input.CareerStats.Rating, input.GlobalVehicleAverages[vehicle.VehicleID])
			if err != nil {
				return cards, fmt.Errorf("failed to generate vehicle %d stats from preset: %w", vehicle.VehicleID, err)
			}

			glossary := input.VehicleGlossary[vehicle.VehicleID]
			glossary.ID = vehicle.VehicleID
			cards.Rating = append(cards.Rating, Card{
				Title:  fmt.Sprintf("%s %s", utils.IntToRoman(glossary.Tier), glossary.Name(options.Locale)),
				Blocks: []StatsBlock{block},
				Type:   dataprep.CardTypeRatingVehicle,
			})
		}
	}

	// Unrated battles
	if allBattles == 0 || input.SessionStats.Global.Battles > 0 {
		var unratedBlocks []StatsBlock
		for _, preset := range options.Blocks {
			if preset == dataprep.TagWN8 {
				// WN8 is a special case that needs to be calculated from vehicles
				sessionWN8 := calculateWeightedWN8(input.SessionStats.Vehicles, input.GlobalVehicleAverages)
				careerWN8 := calculateWeightedWN8(input.CareerStats.Vehicles, input.GlobalVehicleAverages)
				unratedBlocks = append(unratedBlocks, StatsBlock{
					Session: dataprep.StatsToValue(sessionWN8),
					Career:  dataprep.StatsToValue(careerWN8),
					Label:   options.LocalePrinter("label_" + string(dataprep.TagWN8)),
					Tag:     dataprep.TagWN8,
				})
				continue
			}
			block, err := presetToBlock(preset, options.LocalePrinter, input.SessionStats.Global, input.CareerStats.Global)
			if err != nil {
				return cards, fmt.Errorf("failed to generate a unrated stats from preset: %w", err)
			}
			unratedBlocks = append(unratedBlocks, block)
		}
		cards.Unrated = append(cards.Unrated, Card{
			Title:  options.LocalePrinter("label_overview_unrated"),
			Blocks: unratedBlocks,
			Type:   dataprep.CardTypeOverview,
			Meta:   "unrated",
		})
	}

	// Unrated Vehicles
	if input.SessionStats.Global.Battles > 0 {
		for _, vehicle := range input.SessionUnratedVehicles {
			var vehicleBlocks []StatsBlock
			for _, preset := range options.Blocks {
				var career core.ReducedStatsFrame
				if careerStats, ok := input.CareerStats.Vehicles[vehicle.VehicleID]; ok {
					career = *careerStats.ReducedStatsFrame
				}

				block, err := presetToBlock(preset, options.LocalePrinter, *vehicle.ReducedStatsFrame, career, input.GlobalVehicleAverages[vehicle.VehicleID])
				if err != nil {
					return cards, fmt.Errorf("failed to generate vehicle %d stats from preset: %w", vehicle.VehicleID, err)
				}
				vehicleBlocks = append(vehicleBlocks, block)
			}

			glossary := input.VehicleGlossary[vehicle.VehicleID]
			glossary.ID = vehicle.VehicleID
			cards.Unrated = append(cards.Unrated, Card{
				Title:  fmt.Sprintf("%s %s", utils.IntToRoman(glossary.Tier), glossary.Name(options.Locale)),
				Blocks: vehicleBlocks,
				Type:   dataprep.CardTypeVehicle,
			})
		}
	}

	return cards, nil
}
