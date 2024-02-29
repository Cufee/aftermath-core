package period

import "github.com/cufee/aftermath-core/dataprep"

var DefaultBlocks = [][]dataprep.Tag{{dataprep.TagBattles, dataprep.TagAvgDamage, dataprep.TagWN8}, {dataprep.TagWinrate, dataprep.TagAvgTier, dataprep.TagAccuracy}}

type highlight struct {
	preset dataprep.Tag
}

var (
	HighlightAvgDamage = highlight{dataprep.TagAvgDamage}
	HighlightBattles   = highlight{dataprep.TagBattles}
	HighlightWN8       = highlight{dataprep.TagWN8}
)
