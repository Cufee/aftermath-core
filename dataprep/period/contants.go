package period

import "github.com/cufee/aftermath-core/dataprep"

// var DefaultBlocks = [][]dataprep.Tag{{dataprep.TagWN8}, {dataprep.TagBattles, dataprep.TagWinrate}, {dataprep.TagDamageRatio, dataprep.TagAvgDamage, dataprep.TagAccuracy}}
// var DefaultBlocks = [][]dataprep.Tag{{dataprep.TagBattles, dataprep.TagWinrate, dataprep.TagAvgTier}, {dataprep.TagDamageRatio, dataprep.TagAvgDamage, dataprep.TagAccuracy}}
var DefaultBlocks = [][]dataprep.Tag{{dataprep.TagWinrate, dataprep.TagWN8}, {dataprep.TagBattles, dataprep.TagAvgDamage, dataprep.TagAccuracy}, {dataprep.TagAvgTier, dataprep.TagDamageRatio, dataprep.TagSurvivalPercent}}

type highlight struct {
	preset dataprep.Tag
}

var (
	HighlightAvgDamage = highlight{dataprep.TagAvgDamage}
	HighlightBattles   = highlight{dataprep.TagBattles}
	HighlightWN8       = highlight{dataprep.TagWN8}
)
