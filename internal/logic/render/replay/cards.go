package replay

import (
	"image"

	"github.com/cufee/aftermath-core/internal/logic/external"
	"github.com/cufee/aftermath-core/internal/logic/render"
)

// Tank Name 							// WN8 Winrate Damage Kills
// Player Name [Tag]

func RenderReplayImage(replay *external.Replay) (image.Image, error) {
	var alliesBlocks, enemiesBlocks []render.Block

	// Title Card
	titleBlock := newTitleBlock(replay)

	// Overview Column
	overviewBlocks := []render.Block{newProtagonistBlock(replay)}

	// Highlight Cards
	overviewBlocks = append(overviewBlocks, newHighlightCard(replay))
	overviewBlocks = append(overviewBlocks, newHighlightCard(replay))

	// Summary
	overviewBlocks = append(overviewBlocks, newBattleResultCard(replay))

	// Allies
	for _, player := range replay.Teams.Allies {
		alliesBlocks = append(alliesBlocks, newPlayerCard(&player, true))
	}
	// Enemies
	for _, player := range replay.Teams.Enemies {
		enemiesBlocks = append(enemiesBlocks, newPlayerCard(&player, false))
	}

	var frameBlocks []render.Block
	frameBlocks = append(frameBlocks, render.NewBlocksContent(render.Style{Direction: render.DirectionVertical, Gap: 10}, overviewBlocks...))

	var teamsBlocks []render.Block
	teamsBlocks = append(teamsBlocks, render.NewBlocksContent(render.Style{Direction: render.DirectionVertical, Gap: 10}, alliesBlocks...))
	teamsBlocks = append(teamsBlocks, render.NewBlocksContent(render.Style{Direction: render.DirectionVertical, Gap: 10}, enemiesBlocks...))
	playersBlock := render.NewBlocksContent(render.Style{Direction: render.DirectionHorizontal, Gap: 10}, teamsBlocks...)

	rightSideBlocks := render.NewBlocksContent(render.Style{Direction: render.DirectionVertical, Gap: 10}, titleBlock, playersBlock)
	frameBlocks = append(frameBlocks, render.NewBlocksContent(render.Style{Direction: render.DirectionVertical, Gap: 10}, rightSideBlocks))

	frame := render.NewBlocksContent(frameStyle, frameBlocks...)

	return frame.Render()
}
