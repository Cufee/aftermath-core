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
			block, err := preset.StatsBlock(vehicle.ReducedStatsFrame, input.CareerStats.Vehicles[vehicle.VehicleID].ReducedStatsFrame, input.GlobalVehicleAverages[vehicle.VehicleID], printer)
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
			wn8VehiclesTotal += vehicle.WN8(averages[id])
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
	if career == nil {
		career = core.EmptySession(0, 0).Global
	}

	switch *p {
	case BlockPresetWN8:
		return StatsBlock{
			Session: statsValueToString(session.WN8(averages)),
			Career:  statsValueToString(career.WN8(averages)),
			Label:   printer("label_wn8"),
		}, nil
	case BlockPresetBattles:
		return StatsBlock{
			Session: statsValueToString(session.Battles),
			Career:  statsValueToString(career.Battles),
			Label:   printer("label_battles"),
		}, nil
	case BlockPresetWinrate:
		return StatsBlock{
			Session: statsValueToString(session.Winrate()),
			Career:  statsValueToString(career.Winrate()),
			Label:   printer("label_winrate"),
		}, nil
	case BlockPresetAccuracy:
		return StatsBlock{
			Session: statsValueToString(session.Accuracy()),
			Career:  statsValueToString(career.Accuracy()),
			Label:   printer("label_accuracy"),
		}, nil
	case BlockPresetAvgDamage:
		return StatsBlock{
			Session: statsValueToString(int(session.AvgDamage())),
			Career:  statsValueToString(int(career.AvgDamage())),
			Label:   printer("label_avg_damage"),
		}, nil
	case BlockPresetDamageRatio:
		return StatsBlock{
			Session: statsValueToString(session.DamageRatio()),
			Career:  statsValueToString(career.DamageRatio()),
			Label:   printer("label_damage_ratio"),
		}, nil
	}
	return StatsBlock{}, errors.New("invalid preset")
}
