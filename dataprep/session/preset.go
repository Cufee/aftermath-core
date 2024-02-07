package session

import (
	"errors"

	"github.com/cufee/aftermath-core/dataprep"
	"github.com/cufee/aftermath-core/internal/core/localization"
	"github.com/cufee/aftermath-core/internal/core/stats"
)

type Cards []dataprep.StatsCard[StatsBlock, string]

type StatsBlock struct {
	Session dataprep.Value `json:"session"`
	Career  dataprep.Value `json:"career"`
	Label   string         `json:"label"`
	Tag     dataprep.Tag   `json:"tag"`
}

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
		return nil, errors.New("no valid presets")
	}
	return parsed, nil
}

func (p *statsBlockPreset) StatsBlock(session, career, averages *stats.ReducedStatsFrame, printer localization.LocalePrinter) (StatsBlock, error) {
	if session == nil {
		return StatsBlock{}, errors.New("session is nil")
	}

	var block StatsBlock
	switch *p {
	case BlockPresetWN8:
		block.Session = dataprep.StatsToValue(session.WN8(averages))
		block.Label = printer("label_wn8")
		block.Tag = dataprep.TagWN8
		if career != nil {
			block.Career = dataprep.StatsToValue(career.WN8(averages))
		}
	case BlockPresetBattles:
		block.Session = dataprep.StatsToValue(session.Battles)
		block.Label = printer("label_battles")
		block.Tag = dataprep.TagBattles
		if career != nil {
			block.Career = dataprep.StatsToValue(career.Battles)
		}
	case BlockPresetWinrate:
		block.Session = dataprep.StatsToValue(session.Winrate())
		block.Label = printer("label_winrate")
		block.Tag = dataprep.TagWinrate
		if career != nil {
			block.Career = dataprep.StatsToValue(career.Winrate())
		}
	case BlockPresetAccuracy:
		block.Session = dataprep.StatsToValue(session.Accuracy())
		block.Label = printer("label_accuracy")
		block.Tag = dataprep.TagAccuracy
		if career != nil {
			block.Career = dataprep.StatsToValue(career.Accuracy())
		}
	case BlockPresetAvgDamage:
		block.Session = dataprep.StatsToValue(int(session.AvgDamage()))
		block.Label = printer("label_avg_damage")
		block.Tag = dataprep.TagAvgDamage
		if career != nil {
			block.Career = dataprep.StatsToValue(int(career.AvgDamage()))
		}
	case BlockPresetDamageRatio:
		block.Session = dataprep.StatsToValue(session.DamageRatio())
		block.Label = printer("label_damage_ratio")
		block.Tag = dataprep.TagDamageRatio
		if career != nil {
			block.Career = dataprep.StatsToValue(career.DamageRatio())
		}
	default:
		return StatsBlock{}, errors.New("invalid preset")
	}

	return block, nil
}
