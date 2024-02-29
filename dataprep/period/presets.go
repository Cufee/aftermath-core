package period

import (
	"errors"

	"github.com/cufee/aftermath-core/dataprep"
	"github.com/cufee/aftermath-core/internal/core/localization"
	"github.com/cufee/aftermath-core/internal/core/stats"
)

func presetToBlock(preset dataprep.Tag, stats, averages *stats.ReducedStatsFrame, printer localization.LocalePrinter) (StatsBlock, error) {
	if stats == nil {
		return StatsBlock{}, errors.New("session is nil")
	}

	var block StatsBlock
	block.Label = printer("label_" + string(preset))
	block.Tag = preset
	switch preset {
	case dataprep.TagWN8:
		block.Data = dataprep.StatsToValue(stats.WN8(averages))

	case dataprep.TagBattles:
		block.Data = dataprep.StatsToValue(stats.Battles)

	case dataprep.TagWinrate:
		block.Data = dataprep.StatsToValue(stats.Winrate())

	case dataprep.TagAccuracy:
		block.Data = dataprep.StatsToValue(stats.Accuracy())

	case dataprep.TagAvgDamage:
		block.Data = dataprep.StatsToValue(int(stats.AvgDamage()))

	case dataprep.TagDamageRatio:
		block.Data = dataprep.StatsToValue(stats.DamageRatio())

	case dataprep.TagAvgTier:
		block.Data = dataprep.StatsToValue(stats.DamageRatio())

	default:
		return StatsBlock{}, errors.New("invalid preset")
	}

	return block, nil
}
