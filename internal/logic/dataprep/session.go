package dataprep

import (
	"errors"
	"fmt"

	"github.com/cufee/aftermath-core/internal/core/localization"
	core "github.com/cufee/aftermath-core/internal/core/stats"
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

type StatsBlock struct {
	Session string
	Career  string
	Label   string
}

type VehicleStatsBlock struct {
	ID     int
	Blocks []StatsBlock
}

type SessionBlocks struct {
	Rating   []StatsBlock
	Regular  []StatsBlock
	Vehicles []VehicleStatsBlock
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

func SnapshotToSession(input ExportInput, options ExportOptions) (SessionBlocks, error) {
	if input.SessionStats == nil || input.CareerStats == nil {
		return SessionBlocks{}, errors.New("session or career stats are nil")
	}

	var sessionBlocks SessionBlocks
	printer := localization.GetPrinter(options.Locale)

	// Vehicles
	for _, vehicle := range input.SessionVehicles {
		var vehicleBlocks []StatsBlock
		for _, preset := range options.Blocks {
			var career *core.ReducedStatsFrame
			if input.CareerStats.Vehicles[vehicle.VehicleID] != nil {
				career = input.CareerStats.Vehicles[vehicle.VehicleID].ReducedStatsFrame
			}
			block, err := preset.StatsBlock(vehicle.ReducedStatsFrame, career, input.GlobalVehicleAverages[vehicle.VehicleID], printer)
			if err != nil {
				return sessionBlocks, fmt.Errorf("failed to generate vehicle %d stats from preset: %w", vehicle.VehicleID, err)
			}
			vehicleBlocks = append(vehicleBlocks, block)
		}
		sessionBlocks.Vehicles = append(sessionBlocks.Vehicles, VehicleStatsBlock{
			ID:     vehicle.VehicleID,
			Blocks: vehicleBlocks,
		})
	}

	// Rating and Unrated battles
	for _, preset := range options.Blocks {
		{
			p := preset
			if preset == BlockPresetWN8 {
				// Rating battles have no WN8, so we use Accuracy instead of drawing a blank
				p = BlockPresetAccuracy
			}
			ratingBlock, err := p.StatsBlock(input.SessionStats.Rating, input.CareerStats.Rating, nil, printer)
			if err != nil {
				return sessionBlocks, fmt.Errorf("failed to generate a rating stats from preset: %w", err)
			}
			sessionBlocks.Rating = append(sessionBlocks.Rating, ratingBlock)
		}
		{
			if preset == BlockPresetWN8 {
				// WN8 is a special case that needs to be calculated from vehicles
				sessionWN8 := calculateSessionWN8(input.SessionStats.Vehicles, input.GlobalVehicleAverages)
				if sessionWN8 != core.InvalidValue {
					sessionBlocks.Regular = append(sessionBlocks.Regular, StatsBlock{
						Session: statsValueToString(sessionWN8),
						Label:   printer("label_wn8"),
					})
					continue
				}
			}
			block, err := preset.StatsBlock(input.SessionStats.Global, input.CareerStats.Global, nil, printer)
			if err != nil {
				return sessionBlocks, fmt.Errorf("failed to generate a unrated stats from preset: %w", err)
			}
			sessionBlocks.Regular = append(sessionBlocks.Regular, block)
		}
	}

	return sessionBlocks, nil
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
		block.Session = statsValueToString(session.WN8(averages))
		block.Label = printer("label_wn8")
		if career != nil {
			block.Career = statsValueToString(career.WN8(averages))
		}
	case BlockPresetBattles:
		block.Session = statsValueToString(session.Battles)
		block.Label = printer("label_battles")
		if career != nil {
			block.Career = statsValueToString(career.Battles)
		}
	case BlockPresetWinrate:
		block.Session = statsValueToString(session.Winrate())
		block.Label = printer("label_winrate")
		if career != nil {
			block.Career = statsValueToString(career.Winrate())
		}
	case BlockPresetAccuracy:
		block.Session = statsValueToString(session.Accuracy())
		block.Label = printer("label_accuracy")
		if career != nil {
			block.Career = statsValueToString(career.Accuracy())
		}
	case BlockPresetAvgDamage:
		block.Session = statsValueToString(int(session.AvgDamage()))
		block.Label = printer("label_avg_damage")
		if career != nil {
			block.Career = statsValueToString(int(career.AvgDamage()))
		}
	case BlockPresetDamageRatio:
		block.Session = statsValueToString(session.DamageRatio())
		block.Label = printer("label_damage_ratio")
		if career != nil {
			block.Career = statsValueToString(career.DamageRatio())
		}
	default:
		return StatsBlock{}, errors.New("invalid preset")
	}

	return block, nil
}
