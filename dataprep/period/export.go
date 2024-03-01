package period

import (
	"errors"
	"fmt"

	"github.com/cufee/aftermath-core/dataprep"
	"github.com/cufee/aftermath-core/internal/core/database"
	"github.com/cufee/aftermath-core/internal/core/localization"
	"github.com/cufee/aftermath-core/internal/core/utils"
	"github.com/cufee/aftermath-core/internal/logic/stats/period"
	"github.com/rs/zerolog/log"
)

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
type VehicleCard dataprep.StatsCard[StatsBlock, string]

type BlockFlavor string

const (
	BlockFlavorDefault   = "default"
	BlockFlavorSpecial   = "special"
	BlockFlavorSecondary = "secondary"
)

type StatsBlock struct {
	Label  string         `json:"label"`
	Data   dataprep.Value `json:"data"`
	Tag    dataprep.Tag   `json:"tag"`
	Flavor BlockFlavor    `json:"flavor"`
}

func SnapshotToSession(stats *period.PeriodStats, options ExportOptions) (Cards, error) {
	if stats == nil {
		return Cards{}, errors.New("period stats are nil")
	}

	var cards Cards
	printer := localization.GetPrinter(options.Locale)

	var ids []int
	for _, vehicle := range stats.Vehicles {
		ids = append(ids, vehicle.VehicleID)
	}

	vehiclesGlossary, err := database.GetGlossaryVehicles(ids...)
	if err != nil {
		// This is definitely not fatal, but will look ugly
		log.Warn().Err(err).Msg("failed to get vehicles glossary")
	}

	// Overview Card
	for _, row := range options.Blocks {
		var rowBlocks []StatsBlock
		for _, preset := range row {
			if preset == dataprep.TagAvgTier {
				value := calculateAvgTier(stats.Vehicles, vehiclesGlossary)
				rowBlocks = append(rowBlocks, StatsBlock{
					Label:  printer("label_" + string(preset)),
					Data:   dataprep.StatsToValue(value),
					Flavor: BlockFlavorSecondary,
					Tag:    preset,
				})
				continue
			}
			block, err := presetToBlock(preset, &stats.Stats, printer)
			if err != nil {
				return cards, err
			}
			rowBlocks = append(rowBlocks, block)

		}

		cards.Overview.Type = dataprep.CardTypeOverview
		cards.Overview.Blocks = append(cards.Overview.Blocks, rowBlocks)
	}

	if len(stats.Vehicles) < 1 || len(options.Highlights) < 1 {
		return cards, nil
	}

	// Vehicle Highlights
	var minimumBattles int = 5
	periodDays := stats.End.Sub(stats.Start).Hours() / 24
	if periodDays > 90 {
		minimumBattles = 100
	} else if periodDays > 60 {
		minimumBattles = 75
	} else if periodDays > 30 {
		minimumBattles = 50
	} else if periodDays > 14 {
		minimumBattles = 25
	} else if periodDays > 7 {
		minimumBattles = 10
	}

	highlightedVehicles := getHighlightedVehicles(options.Highlights, stats.Vehicles, minimumBattles)
	for _, data := range highlightedVehicles {
		var vehicleBlocks []StatsBlock

		for _, preset := range data.highlight.blocks {
			block, err := presetToBlock(preset, data.vehicle.ReducedStatsFrame, printer)
			if err != nil {
				return cards, fmt.Errorf("failed to generate vehicle %d stats from preset: %w", data.vehicle.VehicleID, err)
			}
			vehicleBlocks = append(vehicleBlocks, block)
		}

		glossary := vehiclesGlossary[data.vehicle.VehicleID]
		glossary.ID = data.vehicle.VehicleID

		cards.Highlights = append(cards.Highlights, VehicleCard{
			Title:  fmt.Sprintf("%s %s", utils.IntToRoman(glossary.Tier), glossary.Name(options.Locale)),
			Type:   dataprep.CardTypeVehicle,
			Blocks: vehicleBlocks,
			Meta:   printer(data.highlight.label),
		})
	}

	return cards, nil
}
