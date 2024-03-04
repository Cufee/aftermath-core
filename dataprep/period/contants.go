package period

import "github.com/cufee/aftermath-core/dataprep"

var DefaultBlocks = [][]dataprep.Tag{{dataprep.TagDamageRatio, dataprep.TagAvgDamage, dataprep.TagAccuracy}, {dataprep.TagWN8, dataprep.TagBattles}, {dataprep.TagAvgTier, dataprep.TagWinrate, dataprep.TagSurvivalPercent}}
var DefaultHighlights = []highlight{HighlightBattles, HighlightWN8, HighlightAvgDamage}
