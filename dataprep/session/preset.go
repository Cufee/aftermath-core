package session

import (
	"errors"

	"github.com/cufee/aftermath-core/dataprep"
	"github.com/cufee/aftermath-core/internal/core/stats"
)

type Card dataprep.StatsCard[StatsBlock, string]

type Cards struct {
	Unrated []Card `json:"unrated"`
	Rating  []Card `json:"rating"`
}

type StatsBlock struct {
	Session dataprep.Value `json:"session"`
	Career  dataprep.Value `json:"career"`
	Label   string         `json:"label"`
	Tag     dataprep.Tag   `json:"tag"`
}

func presetToBlock(preset dataprep.Tag, printer func(string) string, session, career stats.ReducedStatsFrame, averages ...stats.ReducedStatsFrame) (StatsBlock, error) {
	var block StatsBlock
	block.Label = printer("label_" + string(preset))
	block.Tag = preset
	switch preset {
	case dataprep.TagWN8:
		block.Session = dataprep.StatsToValue(session.WN8(averages...))
		if career.Battles > 0 {
			block.Career = dataprep.StatsToValue(career.WN8(averages...))
		}
	case dataprep.TagBattles:
		block.Session = dataprep.StatsToValue(session.Battles)
		if career.Battles > 0 {
			block.Career = dataprep.StatsToValue(career.Battles)
		}
	case dataprep.TagWinrate:
		block.Session = dataprep.StatsToValue(session.Winrate())
		if career.Battles > 0 {
			block.Career = dataprep.StatsToValue(career.Winrate())
		}
	case dataprep.TagAccuracy:
		block.Session = dataprep.StatsToValue(session.Accuracy())
		if career.Battles > 0 {
			block.Career = dataprep.StatsToValue(career.Accuracy())
		}
	case dataprep.TagAvgDamage:
		block.Session = dataprep.StatsToValue(int(session.AvgDamage()))
		if career.Battles > 0 {
			block.Career = dataprep.StatsToValue(int(career.AvgDamage()))
		}
	case dataprep.TagDamageRatio:
		block.Session = dataprep.StatsToValue(session.DamageRatio())
		if career.Battles > 0 {
			block.Career = dataprep.StatsToValue(career.DamageRatio())
		}
	case dataprep.TagRankedRating:
		block.Session = dataprep.StatsToValue(session.Rating())
		if career.Battles > 0 {
			block.Career = dataprep.StatsToValue(career.Rating())
		}
	default:
		return StatsBlock{}, errors.New("invalid preset")
	}

	return block, nil
}
