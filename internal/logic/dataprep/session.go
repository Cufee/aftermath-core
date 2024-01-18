package dataprep

import (
	"errors"
	"fmt"

	"github.com/cufee/aftermath-core/internal/core/database"
	"github.com/cufee/aftermath-core/internal/core/localization"
	core "github.com/cufee/aftermath-core/internal/core/stats"
	"github.com/cufee/aftermath-core/internal/core/utils"
	"github.com/rs/zerolog/log"
)

type statsBlockPreset string

const (
	BlockPresetWN8         statsBlockPreset = "wn8"
	BlockPresetBattles     statsBlockPreset = "battles"
	BlockPresetWinrate     statsBlockPreset = "winrate"
	BlockPresetAccuracy    statsBlockPreset = "accuracy"
	BlockPresetAvgDamage   statsBlockPreset = "avg_damage"
	BlockPresetDamageRatio statsBlockPreset = "damage_ratio"
)

var DefaultBlockPresets = []statsBlockPreset{BlockPresetBattles, BlockPresetAvgDamage, BlockPresetDamageRatio, BlockPresetWinrate, BlockPresetWN8}

type cardType string

const (
	CardTypeVehicle  cardType = "vehicle"
	CardTypeOverview cardType = "overview"
)

type StatsCard struct {
	Type   cardType     `json:"type"`
	Title  string       `json:"title"`
	Blocks []StatsBlock `json:"blocks"`
}

type SessionCards []StatsCard

type value struct {
	Value  any    `json:"value"`
	String string `json:"string"`
}

type StatsBlock struct {
	Session value  `json:"session"`
	Career  value  `json:"career"`
	Label   string `json:"label"`
}

type ExportInput struct {
	CareerStats           *core.SessionSnapshot
	SessionStats          *core.SessionSnapshot
	SessionVehicles       []*core.ReducedVehicleStats
	GlobalVehicleAverages map[int]*core.ReducedStatsFrame
}

type ExportOptions struct {
	Blocks []statsBlockPreset
	Locale localization.SupportedLanguage
}

func SnapshotToSession(input ExportInput, options ExportOptions) (SessionCards, error) {
	if input.SessionStats == nil || input.CareerStats == nil {
		return nil, errors.New("session or career stats are nil")
	}

	var cards SessionCards
	printer := localization.GetPrinter(options.Locale)

	// Unrated battles
	if input.SessionStats.Rating.Battles == 0 || input.SessionStats.Global.Battles > 0 {
		var unratedBlocks []StatsBlock
		for _, preset := range options.Blocks {
			if preset == BlockPresetWN8 {
				// WN8 is a special case that needs to be calculated from vehicles
				sessionWN8 := calculateSessionWN8(input.SessionStats.Vehicles, input.GlobalVehicleAverages)
				if sessionWN8 != core.InvalidValue {
					unratedBlocks = append(unratedBlocks, StatsBlock{
						Session: statsToValue(sessionWN8),
						Label:   printer("label_wn8"),
					})
					continue
				}
			}
			block, err := preset.StatsBlock(input.SessionStats.Global, input.CareerStats.Global, nil, printer)
			if err != nil {
				return nil, fmt.Errorf("failed to generate a unrated stats from preset: %w", err)
			}
			unratedBlocks = append(unratedBlocks, block)
		}
		cards = append(cards, StatsCard{
			Title:  printer("label_overview_unrated"),
			Blocks: unratedBlocks,
			Type:   CardTypeOverview,
		})
	}

	// Rating battles
	if input.SessionStats.Rating.Battles > 0 {
		var ratingBlocks []StatsBlock
		for _, preset := range options.Blocks {
			p := preset
			if preset == BlockPresetWN8 {
				// Rating battles have no WN8, so we use Accuracy instead of drawing a blank
				p = BlockPresetAccuracy
			}
			ratingBlock, err := p.StatsBlock(input.SessionStats.Rating, input.CareerStats.Rating, nil, printer)
			if err != nil {
				return nil, fmt.Errorf("failed to generate a rating stats from preset: %w", err)
			}
			ratingBlocks = append(ratingBlocks, ratingBlock)
		}
		cards = append(cards, StatsCard{
			Title:  printer("label_overview_rating"),
			Blocks: ratingBlocks,
			Type:   CardTypeOverview,
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
				block, err := preset.StatsBlock(vehicle.ReducedStatsFrame, career, input.GlobalVehicleAverages[vehicle.VehicleID], printer)
				if err != nil {
					return nil, fmt.Errorf("failed to generate vehicle %d stats from preset: %w", vehicle.VehicleID, err)
				}
				vehicleBlocks = append(vehicleBlocks, block)
			}

			glossary := vehiclesGlossary[vehicle.VehicleID]
			cards = append(cards, StatsCard{
				Title:  fmt.Sprintf("%s %s", utils.IntToRoman(glossary.Tier), glossary.Name(options.Locale)),
				Blocks: vehicleBlocks,
				Type:   CardTypeVehicle,
			})
		}
	}

	return cards, nil
}

func calculateSessionWN8(vehicles map[int]*core.ReducedVehicleStats, averages map[int]*core.ReducedStatsFrame) int {
	var wn8VehiclesTotal, wn8VehiclesBattles int
	for id, vehicle := range vehicles {
		if vehicle.Valid(vehicle.WN8(averages[id])) {
			wn8VehiclesTotal += vehicle.WN8(averages[id]) * vehicle.Battles
			wn8VehiclesBattles += vehicle.Battles
		}
	}
	if wn8VehiclesBattles > 0 {
		return wn8VehiclesTotal / wn8VehiclesBattles
	}
	return core.InvalidValue
}

func (p *statsBlockPreset) StatsBlock(session, career, averages *core.ReducedStatsFrame, printer localization.LocalePrinter) (StatsBlock, error) {
	if session == nil {
		return StatsBlock{}, errors.New("session is nil")
	}

	var block StatsBlock
	switch *p {
	case BlockPresetWN8:
		block.Session = statsToValue(session.WN8(averages))
		block.Label = printer("label_wn8")
		if career != nil {
			block.Career = statsToValue(career.WN8(averages))
		}
	case BlockPresetBattles:
		block.Session = statsToValue(session.Battles)
		block.Label = printer("label_battles")
		if career != nil {
			block.Career = statsToValue(career.Battles)
		}
	case BlockPresetWinrate:
		block.Session = statsToValue(session.Winrate())
		block.Label = printer("label_winrate")
		if career != nil {
			block.Career = statsToValue(career.Winrate())
		}
	case BlockPresetAccuracy:
		block.Session = statsToValue(session.Accuracy())
		block.Label = printer("label_accuracy")
		if career != nil {
			block.Career = statsToValue(career.Accuracy())
		}
	case BlockPresetAvgDamage:
		block.Session = statsToValue(int(session.AvgDamage()))
		block.Label = printer("label_avg_damage")
		if career != nil {
			block.Career = statsToValue(int(career.AvgDamage()))
		}
	case BlockPresetDamageRatio:
		block.Session = statsToValue(session.DamageRatio())
		block.Label = printer("label_damage_ratio")
		if career != nil {
			block.Career = statsToValue(career.DamageRatio())
		}
	default:
		return StatsBlock{}, errors.New("invalid preset")
	}

	return block, nil
}
