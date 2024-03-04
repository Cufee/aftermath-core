package period

import "github.com/cufee/aftermath-core/dataprep"

var DefaultBlocks = [][]dataprep.Tag{{dataprep.TagAvgDamage, dataprep.TagDamageRatio, dataprep.TagAvgTier}, {dataprep.TagBattles, dataprep.TagWN8}, {dataprep.TagWinrate, dataprep.TagAccuracy, dataprep.TagSurvivalPercent}}
var DefaultHighlights = []highlight{HighlightBattles, HighlightWN8, HighlightAvgDamage}
