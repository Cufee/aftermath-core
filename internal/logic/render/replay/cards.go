package replay

import (
	"fmt"
	"image"

	"github.com/cufee/aftermath-core/dataprep"
	"github.com/cufee/aftermath-core/dataprep/replay"
	"golang.org/x/text/language"

	"github.com/cufee/aftermath-core/internal/core/localization"
	"github.com/cufee/aftermath-core/internal/logic/render"
	parse "github.com/cufee/aftermath-core/internal/logic/replay"
)

// Tank Name 							// WN8 Winrate Damage Blocker Assisted Kills
// Player Name [Tag]

type ReplayData struct {
	Cards  replay.Cards
	Replay *parse.Replay
}

type RenderOptions struct {
	Locale language.Tag
}

func RenderReplayImage(data ReplayData, opts RenderOptions) (image.Image, error) {
	var alliesBlocks, enemiesBlocks []render.Block

	printer := localization.GetPrinter(opts.Locale)

	var playerNameWidth float64
	statsSizes := make(map[dataprep.Tag]float64)
	for _, card := range append(data.Cards.Allies, data.Cards.Enemies...) {
		// Measure player name and tag or vehicle name
		name := card.Meta.Player.Nickname
		if card.Meta.Player.ClanTag != "" {
			name += fmt.Sprintf(" [%s]", card.Meta.Player.ClanTag)
		}
		nameSize := render.MeasureString(name, render.FontLarge)
		tankSize := render.MeasureString(card.Title, render.FontLarge)
		size := nameSize
		if tankSize.TotalWidth > nameSize.TotalWidth {
			size = tankSize
		}
		if size.TotalWidth > playerNameWidth {
			playerNameWidth = size.TotalWidth
		}

		// Measure stats value and label
		for _, block := range card.Blocks {
			valueSize := render.MeasureString(block.Value.String, render.FontLarge)
			labelSize := render.MeasureString(block.Label, render.FontSmall)
			w := valueSize.TotalWidth
			if labelSize.TotalWidth > valueSize.TotalWidth {
				w = labelSize.TotalWidth
			}
			if w > statsSizes[block.Tag] {
				statsSizes[block.Tag] = w
			}
		}
	}

	var totalStatsWidth float64
	for _, width := range statsSizes {
		totalStatsWidth += width
	}

	playerStatsCardStyle := defaultCardStyle(playerNameWidth+(float64(len(statsSizes)*10))+totalStatsWidth, 0)
	totalCardsWidth := (playerStatsCardStyle.Width * 2) - 30

	// Allies
	for _, card := range data.Cards.Allies {
		alliesBlocks = append(alliesBlocks, newPlayerCard(playerStatsCardStyle, statsSizes, card, card.Meta.Player, true, card.Meta.Player.ID == data.Replay.Protagonist.ID))
	}
	// Enemies
	for _, card := range data.Cards.Enemies {
		enemiesBlocks = append(enemiesBlocks, newPlayerCard(playerStatsCardStyle, statsSizes, card, card.Meta.Player, false, false))
	}

	// Title Card
	titleBlock := newTitleBlock(data.Replay, totalCardsWidth, printer)

	// Teams
	var teamsBlocks []render.Block
	teamsBlocks = append(teamsBlocks, render.NewBlocksContent(render.Style{Direction: render.DirectionVertical, Gap: 10}, alliesBlocks...))
	teamsBlocks = append(teamsBlocks, render.NewBlocksContent(render.Style{Direction: render.DirectionVertical, Gap: 10}, enemiesBlocks...))
	playersBlock := render.NewBlocksContent(render.Style{Direction: render.DirectionHorizontal, Gap: 10}, teamsBlocks...)
	teamsBlock := render.NewBlocksContent(render.Style{Direction: render.DirectionVertical, Gap: 10}, playersBlock)

	frame := render.NewBlocksContent(frameStyle, titleBlock, teamsBlock)
	return frame.Render()
}
