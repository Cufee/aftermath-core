package replay

import (
	"fmt"
	"math"

	"github.com/cufee/aftermath-core/dataprep/replay"
	"github.com/cufee/aftermath-core/internal/logic/external"
	"github.com/cufee/aftermath-core/internal/logic/render"
)

func newTitleBlock(replay *external.Replay, width float64) render.Block {
	var titleBlocks []render.Block
	if replay.Victory {
		titleBlocks = append(titleBlocks, render.NewTextContent(render.Style{
			Font:      &render.FontLarge,
			FontColor: render.TextPrimary,
		}, "Victory"))
	} else {
		titleBlocks = append(titleBlocks, render.NewTextContent(render.Style{
			Font:      &render.FontLarge,
			FontColor: render.TextPrimary,
		}, "Defeat"))
	}

	titleBlocks = append(titleBlocks, render.NewTextContent(render.Style{
		Font:      &render.FontLarge,
		FontColor: render.TextSecondary,
	}, fmt.Sprintf(" - %s", replay.BattleType.Name)))

	style := defaultCardStyle(width, 75)
	style.JustifyContent = render.JustifyContentCenter
	style.Direction = render.DirectionHorizontal
	style.AlignItems = render.AlignItemsCenter

	return render.NewBlocksContent(style, titleBlocks...)
}

func newPlayerCard(style render.Style, card replay.Card, player external.Player, ally bool) render.Block {
	hpBarValue := float64(player.HPLeft) / float64((player.Performance.DamageReceived + player.HPLeft))
	if hpBarValue > 0 {
		hpBarValue = math.Max(hpBarValue, 0.2)
	}

	var hpBar render.Block
	if ally {
		hpBar = newProgressBar(60, int(hpBarValue*100), progressDirectionVertical, hpBarColorAllies)
	} else {
		hpBar = newProgressBar(60, int(hpBarValue*100), progressDirectionVertical, hpBarColorEnemies)
	}

	var nameBlocks []render.Block
	nameBlocks = append(nameBlocks, render.NewTextContent(render.Style{
		Font:      &render.FontLarge,
		FontColor: render.TextSecondary,
	}, player.Nickname))
	if player.ClanTag != "" {
		nameBlocks = append(nameBlocks, render.NewTextContent(render.Style{
			Font:      &render.FontLarge,
			FontColor: render.TextSecondary,
		}, fmt.Sprintf("[%s]", player.ClanTag)))
	}

	leftBlock := render.NewBlocksContent(render.Style{
		Direction:  render.DirectionHorizontal,
		AlignItems: render.AlignItemsCenter,
		// Debug:      true,
		Gap:    10,
		Height: 80,
	}, hpBar, render.NewBlocksContent(render.Style{Direction: render.DirectionVertical},
		render.NewTextContent(render.Style{Font: &render.FontLarge, FontColor: render.TextPrimary}, fmt.Sprint(player.VehicleID)),
		render.NewBlocksContent(render.Style{Direction: render.DirectionHorizontal, Gap: 5, AlignItems: render.AlignItemsCenter, Width: playerCardNameWidth}, nameBlocks...),
	))

	var rightBlocks []render.Block
	for _, block := range card.Blocks {
		rightBlocks = append(rightBlocks, statsBlockToBlock(block))
	}

	style.Direction = render.DirectionHorizontal
	style.AlignItems = render.AlignItemsCenter
	style.JustifyContent = render.JustifyContentSpaceBetween
	style.Debug = true

	return render.NewBlocksContent(style, append([]render.Block{leftBlock}, rightBlocks...)...)
}
