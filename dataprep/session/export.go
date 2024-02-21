package session

import (
	"errors"
	"fmt"

	"github.com/cufee/aftermath-core/dataprep"
	"github.com/cufee/aftermath-core/internal/core/database"
	"github.com/cufee/aftermath-core/internal/core/localization"
	core "github.com/cufee/aftermath-core/internal/core/stats"
	"github.com/cufee/aftermath-core/internal/core/utils"
	"github.com/rs/zerolog/log"
)

type ExportInput struct {
	CareerStats           *core.SessionSnapshot
	SessionStats          *core.SessionSnapshot
	SessionVehicles       []*core.ReducedVehicleStats
	GlobalVehicleAverages map[int]*core.ReducedStatsFrame
}

type ExportOptions struct {
	Blocks []dataprep.Tag
	Locale localization.SupportedLanguage
}

func SnapshotToSession(input ExportInput, options ExportOptions) (Cards, error) {
	if input.SessionStats == nil || input.CareerStats == nil {
		return nil, errors.New("session or career stats are nil")
	}

	var cards Cards
	printer := localization.GetPrinter(options.Locale)

	// Rating battles
	if input.SessionStats.Rating.Battles > 0 {
		var ratingBlocks []StatsBlock
		for _, preset := range options.Blocks {
			if preset == dataprep.TagWN8 {
				// Rating battles have no WN8, so we use Accuracy instead of drawing a blank
				preset = dataprep.TagAccuracy
			}
			ratingBlock, err := presetToBlock(preset, input.SessionStats.Rating, input.CareerStats.Rating, nil, printer)
			if err != nil {
				return nil, fmt.Errorf("failed to generate a rating stats from preset: %w", err)
			}
			ratingBlocks = append(ratingBlocks, ratingBlock)
		}
		cards = append(cards, dataprep.StatsCard[StatsBlock, string]{
			Title:  printer("label_overview_rating"),
			Blocks: ratingBlocks,
			Type:   dataprep.CardTypeOverview,
		})
	}

	// Unrated battles
	if input.SessionStats.Rating.Battles == 0 || input.SessionStats.Global.Battles > 0 {
		var unratedBlocks []StatsBlock
		for _, preset := range options.Blocks {
			if preset == dataprep.TagWN8 {
				// WN8 is a special case that needs to be calculated from vehicles
				sessionWN8 := calculateSessionWN8(input.SessionStats.Vehicles, input.GlobalVehicleAverages)
				if sessionWN8 != core.InvalidValueInt {
					unratedBlocks = append(unratedBlocks, StatsBlock{
						Session: dataprep.StatsToValue(sessionWN8),
						Label:   printer("label_" + string(dataprep.TagWN8)),
						Tag:     dataprep.TagWN8,
					})
					continue
				}
			}
			block, err := presetToBlock(preset, input.SessionStats.Global, input.CareerStats.Global, nil, printer)
			if err != nil {
				return nil, fmt.Errorf("failed to generate a unrated stats from preset: %w", err)
			}
			unratedBlocks = append(unratedBlocks, block)
		}
		cards = append(cards, dataprep.StatsCard[StatsBlock, string]{
			Title:  printer("label_overview_unrated"),
			Blocks: unratedBlocks,
			Type:   dataprep.CardTypeOverview,
		})
	}

	// Vehicles
	if len(input.SessionVehicles) > 0 {
		var ids []int
		for _, vehicle := range input.SessionVehicles {
			ids = append(ids, vehicle.VehicleID)
		}

		vehiclesGlossary, err := database.GetGlossaryVehicles(ids...)
		if err != nil {
			// This is definitely not fatal, but will look ugly
			log.Warn().Err(err).Msg("failed to get vehicles glossary")
		}

		for _, vehicle := range input.SessionVehicles {
			var vehicleBlocks []StatsBlock
			for _, preset := range options.Blocks {
				var career *core.ReducedStatsFrame
				if input.CareerStats.Vehicles[vehicle.VehicleID] != nil {
					career = input.CareerStats.Vehicles[vehicle.VehicleID].ReducedStatsFrame
				}
				block, err := presetToBlock(preset, vehicle.ReducedStatsFrame, career, input.GlobalVehicleAverages[vehicle.VehicleID], printer)
				if err != nil {
					return nil, fmt.Errorf("failed to generate vehicle %d stats from preset: %w", vehicle.VehicleID, err)
				}
				vehicleBlocks = append(vehicleBlocks, block)
			}

			glossary := vehiclesGlossary[vehicle.VehicleID]
			glossary.ID = vehicle.VehicleID
			cards = append(cards, dataprep.StatsCard[StatsBlock, string]{
				Title:  fmt.Sprintf("%s %s", utils.IntToRoman(glossary.Tier), glossary.Name(options.Locale)),
				Blocks: vehicleBlocks,
				Type:   dataprep.CardTypeVehicle,
			})
		}
	}

	return cards, nil
}
