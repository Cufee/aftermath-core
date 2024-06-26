package period

import (
	"errors"

	"github.com/cufee/aftermath-core/dataprep"
	"github.com/cufee/aftermath-core/internal/core/stats"
)

func presetToBlock(preset dataprep.Tag, printer func(string) string, stats stats.ReducedStatsFrame) (StatsBlock, error) {
	var block StatsBlock
	block.Label = printer("label_" + string(preset))
	block.Tag = preset
	switch preset {
	case dataprep.TagWN8:
		block.Data = dataprep.StatsToValue(stats.WN8())
		block.Flavor = BlockFlavorSpecial

	case dataprep.TagBattles:
		block.Data = dataprep.StatsToValue(stats.Battles)
		block.Flavor = BlockFlavorSpecial

	case dataprep.TagWinrate:
		block.Data = dataprep.StatsToValue(stats.Winrate())
		block.Flavor = BlockFlavorDefault

	case dataprep.TagSurvivalRatio:
		block.Data = dataprep.StatsToValue(stats.SurvivalRatio())
		block.Flavor = BlockFlavorSecondary

	case dataprep.TagSurvivalPercent:
		block.Data = dataprep.StatsToValue(stats.SurvivalPercent())
		block.Flavor = BlockFlavorSecondary

	case dataprep.TagAccuracy:
		block.Data = dataprep.StatsToValue(stats.Accuracy())
		block.Flavor = BlockFlavorSecondary

	case dataprep.TagAvgDamage:
		block.Data = dataprep.StatsToValue(int(stats.AvgDamage()))
		block.Flavor = BlockFlavorDefault

	case dataprep.TagDamageRatio:
		block.Data = dataprep.StatsToValue(stats.DamageRatio())
		block.Flavor = BlockFlavorSecondary

	default:
		return StatsBlock{}, errors.New("invalid preset")
	}

	return block, nil
}
