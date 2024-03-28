package period

import (
	"fmt"

	"github.com/cufee/aftermath-core/dataprep"
	"github.com/cufee/aftermath-core/internal/core/database/models"
	"github.com/cufee/aftermath-core/internal/core/utils"
	"github.com/cufee/aftermath-core/internal/logic/stats/period"
	"golang.org/x/text/language"
)

type ExportInput struct {
	Stats           period.PeriodStats
	VehicleGlossary map[int]models.Vehicle
}

type ExportOptions struct {
	Locale        language.Tag
	LocalePrinter func(string) string

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

func SnapshotToSession(input ExportInput, options ExportOptions) (Cards, error) {
	if options.LocalePrinter == nil {
		options.LocalePrinter = func(s string) string { return s }
	}
	if input.VehicleGlossary == nil {
		input.VehicleGlossary = make(map[int]models.Vehicle)
	}

	var cards Cards

	// Overview Card
	for _, column := range options.Blocks {
		var columnBlocks []StatsBlock
		for _, preset := range column {
			if preset == dataprep.TagAvgTier {
				value := calculateAvgTier(input.Stats.Vehicles, input.VehicleGlossary)
				columnBlocks = append(columnBlocks, StatsBlock{
					Label:  options.LocalePrinter("label_" + string(preset)),
					Data:   dataprep.StatsToValue(value),
					Flavor: BlockFlavorSecondary,
					Tag:    preset,
				})
				continue
			}
			block, err := presetToBlock(preset, options.LocalePrinter, input.Stats.Stats)
			if err != nil {
				return cards, err
			}
			columnBlocks = append(columnBlocks, block)

		}

		cards.Overview.Type = dataprep.CardTypeOverview
		cards.Overview.Blocks = append(cards.Overview.Blocks, columnBlocks)
	}

	if len(input.Stats.Vehicles) < 1 || len(options.Highlights) < 1 {
		return cards, nil
	}

	// Vehicle Highlights
	var minimumBattles int = 5
	periodDays := input.Stats.End.Sub(input.Stats.Start).Hours() / 24
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

	highlightedVehicles := getHighlightedVehicles(options.Highlights, input.Stats.Vehicles, minimumBattles)
	for _, data := range highlightedVehicles {
		var vehicleBlocks []StatsBlock

		for _, preset := range data.highlight.blocks {
			block, err := presetToBlock(preset, options.LocalePrinter, data.vehicle.ReducedStatsFrame)
			if err != nil {
				return cards, fmt.Errorf("failed to generate vehicle %d stats from preset: %w", data.vehicle.VehicleID, err)
			}
			vehicleBlocks = append(vehicleBlocks, block)
		}

		glossary := input.VehicleGlossary[data.vehicle.VehicleID]
		glossary.ID = data.vehicle.VehicleID

		cards.Highlights = append(cards.Highlights, VehicleCard{
			Title:  fmt.Sprintf("%s %s", utils.IntToRoman(glossary.Tier), glossary.Name(options.Locale)),
			Type:   dataprep.CardTypeVehicle,
			Blocks: vehicleBlocks,
			Meta:   options.LocalePrinter(data.highlight.label),
		})
	}

	return cards, nil
}
