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

func presetToBlock(preset dataprep.Tag, session, career, averages *stats.ReducedStatsFrame, printer localization.LocalePrinter) (StatsBlock, error) {
	if session == nil {
		return StatsBlock{}, errors.New("session is nil")
	}

	var block StatsBlock
	block.Label = printer("label_" + string(preset))
	block.Tag = preset
	switch preset {
	case dataprep.TagWN8:
		block.Session = dataprep.StatsToValue(session.WN8(averages))
		if career != nil {
			block.Career = dataprep.StatsToValue(career.WN8(averages))
		}
	case dataprep.TagBattles:
		block.Session = dataprep.StatsToValue(session.Battles)
		if career != nil {
			block.Career = dataprep.StatsToValue(career.Battles)
		}
	case dataprep.TagWinrate:
		block.Session = dataprep.StatsToValue(session.Winrate())
		if career != nil {
			block.Career = dataprep.StatsToValue(career.Winrate())
		}
	case dataprep.TagAccuracy:
		block.Session = dataprep.StatsToValue(session.Accuracy())
		if career != nil {
			block.Career = dataprep.StatsToValue(career.Accuracy())
		}
	case dataprep.TagAvgDamage:
		block.Session = dataprep.StatsToValue(int(session.AvgDamage()))
		if career != nil {
			block.Career = dataprep.StatsToValue(int(career.AvgDamage()))
		}
	case dataprep.TagDamageRatio:
		block.Session = dataprep.StatsToValue(session.DamageRatio())
		if career != nil {
			block.Career = dataprep.StatsToValue(career.DamageRatio())
		}
	default:
		return StatsBlock{}, errors.New("invalid preset")
	}

	return block, nil
}
