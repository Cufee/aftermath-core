package replay

import (
	"image"

	"github.com/cufee/aftermath-core/dataprep"
	"github.com/cufee/aftermath-core/dataprep/replay"
	"github.com/cufee/aftermath-core/internal/core/database/models"
	"github.com/cufee/aftermath-core/internal/core/stats"
	"github.com/cufee/aftermath-core/internal/logic/external"
	"github.com/cufee/aftermath-core/internal/logic/render"
)

// Tank Name 							// WN8 Winrate Damage Blocker Assisted Kills
// Player Name [Tag]

type ReplayData struct {
	Cards  replay.Cards
	Replay *external.Replay
	// Protagonist *wg.Account
	Glossary map[int]models.Vehicle
	Averages map[int]*stats.ReducedStatsFrame
}

func RenderReplayImage(data ReplayData) (image.Image, error) {
	var alliesBlocks, enemiesBlocks []render.Block

	var tags []dataprep.Tag
	for _, card := range data.Cards.Allies {
		if card.Meta.Tags != nil {
			tags = card.Meta.Tags
			break
		}
	}

	playerStatsCardStyle := playerCardStyle(tags)
	totalCardsWidth := (playerStatsCardStyle.Width * 2) - 30

	// Allies
	for _, card := range data.Cards.Allies {
		alliesBlocks = append(alliesBlocks, newPlayerCard(playerStatsCardStyle, card, card.Meta.Player, true))
	}
	// Enemies
	for _, card := range data.Cards.Enemies {
		enemiesBlocks = append(enemiesBlocks, newPlayerCard(playerStatsCardStyle, card, card.Meta.Player, false))
	}

	// Title Card
	titleBlock := newTitleBlock(data.Replay, totalCardsWidth)

	// Teams
	var teamsBlocks []render.Block
	teamsBlocks = append(teamsBlocks, render.NewBlocksContent(render.Style{Direction: render.DirectionVertical, Gap: 10}, alliesBlocks...))
	teamsBlocks = append(teamsBlocks, render.NewBlocksContent(render.Style{Direction: render.DirectionVertical, Gap: 10}, enemiesBlocks...))
	playersBlock := render.NewBlocksContent(render.Style{Direction: render.DirectionHorizontal, Gap: 10}, teamsBlocks...)
	teamsBlock := render.NewBlocksContent(render.Style{Direction: render.DirectionVertical, Gap: 10}, playersBlock)

	frame := render.NewBlocksContent(frameStyle, titleBlock, teamsBlock)
	return frame.Render()
}
