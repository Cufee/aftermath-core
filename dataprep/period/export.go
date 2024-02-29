package period

import (
	"errors"

	"github.com/cufee/aftermath-core/dataprep"
	"github.com/cufee/aftermath-core/internal/core/localization"
	core "github.com/cufee/aftermath-core/internal/core/stats"
	"github.com/cufee/aftermath-core/internal/logic/stats/period"
)

type ExportInput struct {
	Stats                 *period.PeriodStats
	GlobalVehicleAverages map[int]*core.ReducedStatsFrame
}

type ExportOptions struct {
	Locale     localization.SupportedLanguage
	Blocks     [][]dataprep.Tag
	Highlights []highlight
}

type Cards struct {
	Overview   OverviewCard  `json:"overview"`
	Highlights []VehicleCard `json:"highlights"`
}

type OverviewCard dataprep.StatsCard[[]StatsBlock, string]
type VehicleCard dataprep.StatsCard[[]StatsBlock, string]

type StatsBlock struct {
	Label string         `json:"label"`
	Data  dataprep.Value `json:"data"`
	Tag   dataprep.Tag   `json:"tag"`
}

func SnapshotToSession(input ExportInput, options ExportOptions) (Cards, error) {
	if input.Stats == nil {
		return Cards{}, errors.New("period stats are nil")
	}

	var cards Cards
	printer := localization.GetPrinter(options.Locale)

	// Overview Card
	for _, row := range options.Blocks {
		var rowBlocks []StatsBlock
		for _, preset := range row {
			block, err := presetToBlock(preset, &input.Stats.Stats, nil, printer)
			if err != nil {
				return cards, err
			}
			rowBlocks = append(rowBlocks, block)
		}
		cards.Overview.Blocks = append(cards.Overview.Blocks, rowBlocks)
	}

	// Vehicle Highlights

	// // Vehicles
	// if len(input.SessionVehicles) > 0 {
	// 	var ids []int
	// 	for _, vehicle := range input.SessionVehicles {
	// 		ids = append(ids, vehicle.VehicleID)
	// 	}

	// 	vehiclesGlossary, err := database.GetGlossaryVehicles(ids...)
	// 	if err != nil {
	// 		// This is definitely not fatal, but will look ugly
	// 		log.Warn().Err(err).Msg("failed to get vehicles glossary")
	// 	}

	// 	for _, vehicle := range input.SessionVehicles {
	// 		var vehicleBlocks []StatsBlock
	// 		for _, preset := range options.Blocks {
	// 			var career *core.ReducedStatsFrame
	// 			if input.CareerStats.Vehicles[vehicle.VehicleID] != nil {
	// 				career = input.CareerStats.Vehicles[vehicle.VehicleID].ReducedStatsFrame
	// 			}
	// 			block, err := presetToBlock(preset, vehicle.ReducedStatsFrame, career, input.GlobalVehicleAverages[vehicle.VehicleID], printer)
	// 			if err != nil {
	// 				return nil, fmt.Errorf("failed to generate vehicle %d stats from preset: %w", vehicle.VehicleID, err)
	// 			}
	// 			vehicleBlocks = append(vehicleBlocks, block)
	// 		}

	// 		glossary := vehiclesGlossary[vehicle.VehicleID]
	// 		glossary.ID = vehicle.VehicleID
	// 		cards = append(cards, dataprep.StatsCard[StatsBlock, string]{
	// 			Title:  fmt.Sprintf("%s %s", utils.IntToRoman(glossary.Tier), glossary.Name(options.Locale)),
	// 			Blocks: vehicleBlocks,
	// 			Type:   dataprep.CardTypeVehicle,
	// 		})
	// 	}
	// }

	return cards, nil
}
