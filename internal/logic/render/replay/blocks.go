package replay

import (
	"fmt"
	"math"

	"github.com/cufee/aftermath-core/internal/logic/external"
	"github.com/cufee/aftermath-core/internal/logic/render"
)

func newTitleBlock(replay *external.Replay) render.Block {
	return render.NewBlocksContent(defaultCardStyle(playerCardWidth*2+10, 75), render.NewTextContent(render.Style{
		Font:      &render.FontLarge,
		FontColor: render.TextPrimary,
	}, "Title Card"))

}

func newProtagonistBlock(replay *external.Replay) render.Block {
	return render.NewBlocksContent(highlightCardStyle(overviewWidth, 200), render.NewTextContent(render.Style{
		Font:      &render.FontLarge,
		FontColor: render.TextPrimary,
	}, "Protagonist Card"))
}

func newHighlightCard(replay *external.Replay) render.Block {
	return render.NewBlocksContent(defaultCardStyle(overviewWidth, 100), render.NewTextContent(render.Style{
		Font:      &render.FontLarge,
		FontColor: render.TextPrimary,
	}, "Highlight Card"))
}

func newPlayerCard(player *external.Player, ally bool) render.Block {
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
		render.NewBlocksContent(render.Style{Direction: render.DirectionHorizontal, Gap: 5, AlignItems: render.AlignItemsCenter}, nameBlocks...),
	))

	style := defaultCardStyle(playerCardWidth, 100)
	if player.HPLeft == 0 {
		style = deadPlayerCardStyle(style.Width, style.Height)
	}

	return render.NewBlocksContent(style, leftBlock)
}

func newBattleResultCard(replay *external.Replay) render.Block {
	return render.NewBlocksContent(render.Style{Width: overviewWidth}, render.NewTextContent(render.Style{
		Font:      &render.FontLarge,
		FontColor: render.TextPrimary,
	}, "Battle Result Card, medals etc"))
}
