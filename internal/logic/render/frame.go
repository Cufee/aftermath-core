package render

import (
	"github.com/cufee/aftermath-core/internal/core/localization"
	core "github.com/cufee/aftermath-core/internal/core/stats"
)

type blockLabelTag string

func (tag blockLabelTag) String() string {
	return string(tag)
}

const (
	blockLabelTagBattles  blockLabelTag = "battle"
	blockLabelTagAvgDmg   blockLabelTag = "avgDmg"
	blockLabelTagWinrate  blockLabelTag = "winrate"
	blockLabelTagAccuracy blockLabelTag = "accuracy"
	blockLabelTagWN8      blockLabelTag = "wn8"

	blockLabelTagNone blockLabelTag = "none"
)

func FrameToBlocks(session, career, averages *core.ReducedStatsFrame, locale localization.SupportedLanguage, options *RenderOptions) []block {
	if session == nil {
		return nil
	}

	localePrinter := localization.GetPrinter(locale)

	var blocks []block
	{
		// Battles
		values := []interface{}{session.Battles}
		if career != nil {
			values = append(values, career.Battles)
		}
		blocks = append(blocks, NewBlock(localePrinter(blockLabelTagBattles.String()), options, values...))
	}
	{
		// Avg Damage
		values := []interface{}{int(session.AvgDamage())}
		if career != nil {
			values = append(values, int(career.AvgDamage()))
		}
		blocks = append(blocks, NewBlock(localePrinter(blockLabelTagAvgDmg.String()), options, values...))
	}
	{
		// Winrate
		values := []interface{}{session.Winrate()}
		if career != nil {
			values = append(values, career.Winrate())
		}
		blocks = append(blocks, NewBlock(localePrinter(blockLabelTagWinrate.String()), options, values...))
	}
	{
		if session.WN8(averages) != core.InvalidValue {
			// WN8
			values := []interface{}{int(session.WN8(averages))}
			if career != nil {
				values = append(values, int(career.WN8(averages)))
			}
			blocks = append(blocks, NewBlock(localePrinter(blockLabelTagWN8.String()), options, values...))
		} else {
			// Fallback to Accuracy to keep the UI consistent
			values := []interface{}{session.Accuracy()}
			if career != nil {
				values = append(values, career.Accuracy())
			}
			blocks = append(blocks, NewBlock(localePrinter(blockLabelTagAccuracy.String()), options, values...))
		}
	}

	return blocks
}

// func FrameToSlimStatsBlocks(session, career, averages *core.ReducedStatsFrame, config *BlockRenderConfig) []block {
// 	var blocks []block
// 	// Winrate (Battles)
// 	// TODO: Make this into two blocks to allow for different styling
// 	blocks = append(blocks, config.CompleteBlock(blockLabelTagNone, fmt.Sprintf("%.2f%% (%d)", session.Winrate(), session.Battles), nil))
// 	// Avg Damage
// 	blocks = append(blocks, config.CompleteBlock(blockLabelTagNone, int(session.AvgDamage()), nil))
// 	if session.WN8(averages) != core.InvalidValue {
// 		// WN8
// 		blocks = append(blocks, config.CompleteBlock(blockLabelTagNone, int(session.WN8(averages)), nil))
// 	} else {
// 		// Fallback to Accuracy to keep the UI consistent
// 		blocks = append(blocks, config.CompleteBlock(blockLabelTagNone, session.Accuracy(), nil))
// 	}

// 	return blocks
// }
