package dataprep

import (
	"errors"
	"fmt"

	"github.com/cufee/aftermath-core/internal/core/localization"
	"github.com/cufee/aftermath-core/internal/core/stats"
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

func ParsePresets(presets ...string) ([]statsBlockPreset, error) {
	var parsed []statsBlockPreset
	for _, preset := range presets {
		if preset == "" {
			continue
		}
		switch preset {
		case string(BlockPresetWN8):
			parsed = append(parsed, BlockPresetWN8)
		case string(BlockPresetBattles):
			parsed = append(parsed, BlockPresetBattles)
		case string(BlockPresetWinrate):
			parsed = append(parsed, BlockPresetWinrate)
		case string(BlockPresetAccuracy):
			parsed = append(parsed, BlockPresetAccuracy)
		case string(BlockPresetAvgDamage):
			parsed = append(parsed, BlockPresetAvgDamage)
		case string(BlockPresetDamageRatio):
			parsed = append(parsed, BlockPresetDamageRatio)
		default:
			return nil, errors.New("invalid preset" + preset)
		}
	}

	if len(parsed) == 0 {
		return DefaultBlockPresets, nil
	}
	return parsed, nil
}

var DefaultBlockPresets = []statsBlockPreset{BlockPresetBattles, BlockPresetAvgDamage, BlockPresetDamageRatio, BlockPresetWinrate, BlockPresetWN8}

func (p *statsBlockPreset) StatsBlock(session, career, averages *stats.ReducedStatsFrame, printer localization.LocalePrinter) (StatsBlock, error) {
	if session == nil {
		return StatsBlock{}, errors.New("session is nil")
	}

	var block StatsBlock
	switch *p {
	case BlockPresetWN8:
		block.Session = statsToValue(session.WN8(averages))
		block.Label = printer("label_wn8")
		block.Tag = TagWN8
		if career != nil {
			block.Career = statsToValue(career.WN8(averages))
		}
	case BlockPresetBattles:
		block.Session = statsToValue(session.Battles)
		block.Label = printer("label_battles")
		block.Tag = TagBattles
		if career != nil {
			block.Career = statsToValue(career.Battles)
		}
	case BlockPresetWinrate:
		block.Session = statsToValue(session.Winrate())
		block.Label = printer("label_winrate")
		block.Tag = TagWinrate
		if career != nil {
			block.Career = statsToValue(career.Winrate())
		}
	case BlockPresetAccuracy:
		block.Session = statsToValue(session.Accuracy())
		block.Label = printer("label_accuracy")
		block.Tag = TagAccuracy
		if career != nil {
			block.Career = statsToValue(career.Accuracy())
		}
	case BlockPresetAvgDamage:
		block.Session = statsToValue(int(session.AvgDamage()))
		block.Label = printer("label_avg_damage")
		block.Tag = TagAvgDamage
		if career != nil {
			block.Career = statsToValue(int(career.AvgDamage()))
		}
	case BlockPresetDamageRatio:
		block.Session = statsToValue(session.DamageRatio())
		block.Label = printer("label_damage_ratio")
		block.Tag = TagDamageRatio
		if career != nil {
			block.Career = statsToValue(career.DamageRatio())
		}
	default:
		return StatsBlock{}, errors.New("invalid preset")
	}

	return block, nil
}

func statsToValue(v any) Value {
	switch cast := v.(type) {
	case float32:
		if int(cast) == stats.InvalidValue {
			return Value{String: "-", Value: float64(cast)}
		}
		return Value{String: fmt.Sprintf("%.2f", cast), Value: float64(cast)}
	case float64:
		if int(cast) == stats.InvalidValue {
			return Value{String: "-", Value: cast}
		}
		return Value{String: fmt.Sprintf("%.2f%%", cast), Value: cast}
	case int:
		if cast == stats.InvalidValue {
			return Value{String: "-", Value: float64(cast)}
		}
		return Value{String: fmt.Sprint(cast), Value: float64(cast)}
	default:
		return Value{String: "-", Value: 0}
	}
}
