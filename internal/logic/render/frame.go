package render

import (
	"fmt"

	"github.com/cufee/aftermath-core/internal/core/localization"
	core "github.com/cufee/aftermath-core/internal/core/stats"
)

func FrameToLargeStatsBlocks(session, career, averages *core.ReducedStatsFrame, locale *localization.SupportedLanguage, config *BlockRenderConfig) blockSet {
	var blocks []block
	// Battles
	blocks = append(blocks, config.CompleteBlock(blockLabelTagBattles, session.Battles, career.Battles))
	// Avg Damage
	blocks = append(blocks, config.CompleteBlock(blockLabelTagAvgDmg, int(session.AvgDamage()), int(career.AvgDamage())))
	// Winrate
	blocks = append(blocks, config.CompleteBlock(blockLabelTagAccuracy, session.Winrate(), career.Winrate()))
	if session.WN8(averages) != core.InvalidValue {
		// WN8
		blocks = append(blocks, config.CompleteBlock(blockLabelTagWN8, int(session.WN8(averages)), int(career.WN8(averages))))
	} else {
		// Fallback to Accuracy to keep the UI consistent
		blocks = append(blocks, config.CompleteBlock(blockLabelTagAccuracy, session.Accuracy(), career.Accuracy()))
	}

	return blockSet{blocks: blocks, style: config.SetStyle}
}

func FrameToSlimStatsBlocks(session, career, averages *core.ReducedStatsFrame, locale *localization.SupportedLanguage, config *BlockRenderConfig) blockSet {
	var blocks []block
	// Winrate (Battles)
	blocks = append(blocks, config.CompleteBlock(blockLabelTagNone, fmt.Sprintf("%.2f%% (%d)", session.Winrate(), session.Battles), nil))
	// Avg Damage
	blocks = append(blocks, config.CompleteBlock(blockLabelTagNone, int(session.AvgDamage()), nil))
	if session.WN8(averages) != core.InvalidValue {
		// WN8
		blocks = append(blocks, config.CompleteBlock(blockLabelTagNone, int(session.WN8(averages)), nil))
	} else {
		// Fallback to Accuracy to keep the UI consistent
		blocks = append(blocks, config.CompleteBlock(blockLabelTagNone, session.Accuracy(), nil))
	}

	return blockSet{blocks: blocks, style: config.SetStyle}
}
