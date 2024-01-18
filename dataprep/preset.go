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

func statsToValue(v any) value {
	switch cast := v.(type) {
	case string:
		return value{String: cast, Value: v}
	case float32:
		if int(cast) == stats.InvalidValue {
			return value{String: "-", Value: v}
		}
		return value{String: fmt.Sprintf("%.2f", cast), Value: v}
	case float64:
		if int(cast) == stats.InvalidValue {
			return value{String: "-", Value: v}
		}
		return value{String: fmt.Sprintf("%.2f%%", cast), Value: v}
	case int:
		if cast == stats.InvalidValue {
			return value{String: "-", Value: v}
		}
		return value{String: fmt.Sprint(cast), Value: v}
	default:
		return value{String: "-", Value: v}
	}
}
