package replay

import (
	"image"
	"slices"

	"github.com/cufee/aftermath-core/internal/logic/external"
	"github.com/cufee/aftermath-core/internal/logic/render"
	wg "github.com/cufee/aftermath-core/types"
)

// Tank Name 							// WN8 Winrate Damage Blocker Assisted Kills
// Player Name [Tag]

type ReplayData struct {
	Replay      *external.Replay
	Protagonist *wg.Account
}

func RenderReplayImage(data ReplayData) (image.Image, error) {
	var alliesBlocks, enemiesBlocks []render.Block

	// Overview Column
	overviewBlocks := []render.Block{newProtagonistBlock(data.Replay)}
	// Highlight Cards
	highlightCards := render.NewBlocksContent(render.Style{Direction: render.DirectionVertical, Gap: 10}, newHighlightCard(data.Replay), newHighlightCard(data.Replay))
	// Summary
	overviewBlocks = append(overviewBlocks, highlightCards, newBattleResultCard(data.Replay))

	presets := []blockPreset{blockPresetDamageDealt, blockPresetDamageAssistedAndBlocked, blockPresetKills}
	if !data.Replay.GameMode.Special {
		presets = append([]blockPreset{blockPresetWN8}, presets...)
	}
	playerStatsCardStyle := playerCardStyle(presets)

	// Title Card
	titleBlock := newTitleBlock(data.Replay, (playerStatsCardStyle.Width*2)-30)

	sortTeams(data.Replay.Teams)
	// Allies
	for _, player := range data.Replay.Teams.Allies {
		alliesBlocks = append(alliesBlocks, newPlayerCard(&player, true, presets))
	}
	// Enemies
	for _, player := range data.Replay.Teams.Enemies {
		enemiesBlocks = append(enemiesBlocks, newPlayerCard(&player, false, presets))
	}

	var frameBlocks []render.Block
	frameBlocks = append(frameBlocks, render.NewBlocksContent(render.Style{Direction: render.DirectionHorizontal, Gap: 20}, overviewBlocks...))

	var teamsBlocks []render.Block
	teamsBlocks = append(teamsBlocks, render.NewBlocksContent(render.Style{Direction: render.DirectionVertical, Gap: 10}, alliesBlocks...))
	teamsBlocks = append(teamsBlocks, render.NewBlocksContent(render.Style{Direction: render.DirectionVertical, Gap: 10}, enemiesBlocks...))
	playersBlock := render.NewBlocksContent(render.Style{Direction: render.DirectionHorizontal, Gap: 10}, teamsBlocks...)

	rightSideBlocks := render.NewBlocksContent(render.Style{Direction: render.DirectionVertical, Gap: 10}, titleBlock, playersBlock)
	frameBlocks = append(frameBlocks, render.NewBlocksContent(render.Style{Direction: render.DirectionVertical, Gap: 10}, rightSideBlocks))

	frame := render.NewBlocksContent(frameStyle, frameBlocks...)

	return frame.Render()
}

func sortTeams(teams external.Teams) {
	sortPlayers(teams.Allies)
	sortPlayers(teams.Enemies)
}

func sortPlayers(players []external.Player) {
	slices.SortFunc(players, func(j, i external.Player) int {
		return (i.Performance.DamageDealt + i.Performance.DamageAssisted + i.Performance.DamageBlocked) - (j.Performance.DamageDealt - j.Performance.DamageAssisted - j.Performance.DamageBlocked)
	})
}
