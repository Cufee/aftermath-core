package period

import "github.com/cufee/aftermath-core/dataprep"

// var DefaultBlocks = [][]dataprep.Tag{{dataprep.TagWN8}, {dataprep.TagBattles, dataprep.TagWinrate}, {dataprep.TagDamageRatio, dataprep.TagAvgDamage, dataprep.TagAccuracy}}
// var DefaultBlocks = [][]dataprep.Tag{{dataprep.TagBattles, dataprep.TagWinrate, dataprep.TagAvgTier}, {dataprep.TagDamageRatio, dataprep.TagAvgDamage, dataprep.TagAccuracy}}
var DefaultBlocks = [][]dataprep.Tag{{dataprep.TagAvgDamage, dataprep.TagDamageRatio, dataprep.TagAvgTier}, {dataprep.TagBattles, dataprep.TagWN8}, {dataprep.TagWinrate, dataprep.TagAccuracy, dataprep.TagSurvivalPercent}}
var DefaultHighlights = []highlight{HighlightWN8, HighlightAvgDamage, HighlightBattles}
