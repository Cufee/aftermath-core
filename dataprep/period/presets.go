package period

import (
	"errors"

	"github.com/cufee/aftermath-core/dataprep"
	"github.com/cufee/aftermath-core/internal/core/localization"
	"github.com/cufee/aftermath-core/internal/core/stats"
)

func presetToBlock(preset dataprep.Tag, stats *stats.ReducedStatsFrame, printer localization.LocalePrinter) (StatsBlock, error) {
	if stats == nil {
		return StatsBlock{}, errors.New("session is nil")
	}

	var block StatsBlock
	block.Label = printer("label_" + string(preset))
	block.Tag = preset
	switch preset {
	case dataprep.TagWN8:
		block.Data = dataprep.StatsToValue(stats.WN8(nil))
		block.Flavor = BlockFlavorSpecial

	case dataprep.TagBattles:
		block.Data = dataprep.StatsToValue(stats.Battles)
		block.Flavor = BlockFlavorDefault

	case dataprep.TagWinrate:
		block.Data = dataprep.StatsToValue(stats.Winrate())
		block.Flavor = BlockFlavorSpecial

	case dataprep.TagSurvivalRatio:
		block.Data = dataprep.StatsToValue(stats.SurvivalRatio())
		block.Flavor = BlockFlavorSecondary

	case dataprep.TagSurvivalPercent:
		block.Data = dataprep.StatsToValue(stats.SurvivalPercent())
		block.Flavor = BlockFlavorSecondary

	case dataprep.TagAccuracy:
		block.Data = dataprep.StatsToValue(stats.Accuracy())
		block.Flavor = BlockFlavorDefault

	case dataprep.TagAvgDamage:
		block.Data = dataprep.StatsToValue(int(stats.AvgDamage()))
		block.Flavor = BlockFlavorSpecial

	case dataprep.TagDamageRatio:
		block.Data = dataprep.StatsToValue(stats.DamageRatio())
		block.Flavor = BlockFlavorSecondary

	default:
		return StatsBlock{}, errors.New("invalid preset")
	}

	return block, nil
}
