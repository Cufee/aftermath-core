package replay

import (
	"fmt"
	"math"

	"github.com/cufee/aftermath-core/dataprep"
	"github.com/cufee/aftermath-core/dataprep/replay"

	"github.com/cufee/aftermath-core/internal/core/localization"
	"github.com/cufee/aftermath-core/internal/logic/render"
	parse "github.com/cufee/aftermath-core/internal/logic/replay"
)

func newTitleBlock(replay *parse.Replay, width float64, printer localization.LocalePrinter) render.Block {
	var titleBlocks []render.Block
	if replay.Victory {
		titleBlocks = append(titleBlocks, render.NewTextContent(render.Style{
			Font:      &render.FontLarge,
			FontColor: render.TextPrimary,
		}, printer("label_victory")))
	} else {
		titleBlocks = append(titleBlocks, render.NewTextContent(render.Style{
			Font:      &render.FontLarge,
			FontColor: render.TextPrimary,
		}, printer("label_defeat")))
	}

	titleBlocks = append(titleBlocks, render.NewTextContent(render.Style{
		Font:      &render.FontLarge,
		FontColor: render.TextSecondary,
	}, fmt.Sprintf(" - %s", printer("label_"+replay.BattleType.String()))))

	style := defaultCardStyle(width, 75)
	style.JustifyContent = render.JustifyContentCenter
	style.Direction = render.DirectionHorizontal
	style.AlignItems = render.AlignItemsCenter

	return render.NewBlocksContent(style, titleBlocks...)
}

func newPlayerCard(style render.Style, sizes map[dataprep.Tag]float64, card replay.Card, player parse.Player, ally, protagonist bool) render.Block {
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

	vehicleColor := render.TextPrimary
	if player.HPLeft <= 0 {
		vehicleColor = render.TextSecondary
	}

	leftBlock := render.NewBlocksContent(render.Style{
		Direction:  render.DirectionHorizontal,
		AlignItems: render.AlignItemsCenter,
		Gap:        10,
		Height:     80,
		// Debug:      true,
	}, hpBar, render.NewBlocksContent(render.Style{Direction: render.DirectionVertical},
		render.NewTextContent(render.Style{Font: &render.FontLarge, FontColor: vehicleColor}, card.Title),
		playerNameBlock(player, protagonist),
	))

	var rightBlocks []render.Block
	for _, block := range card.Blocks {
		rightBlocks = append(rightBlocks, statsBlockToBlock(block, sizes[block.Tag]))
	}
	rightBlock := render.NewBlocksContent(render.Style{
		JustifyContent: render.JustifyContentCenter,
		AlignItems:     render.AlignItemsCenter,
		Gap:            10,
		// Debug: true,
	}, rightBlocks...)

	style.Direction = render.DirectionHorizontal
	style.AlignItems = render.AlignItemsCenter
	style.JustifyContent = render.JustifyContentSpaceBetween
	// style.Debug = true

	return render.NewBlocksContent(style, leftBlock, rightBlock)
}

func playerNameBlock(player parse.Player, protagonist bool) render.Block {
	nameColor := render.TextSecondary
	if protagonist {
		nameColor = protagonistColor
	}

	var nameBlocks []render.Block
	nameBlocks = append(nameBlocks, render.NewTextContent(render.Style{
		Font:      &render.FontLarge,
		FontColor: nameColor,
		// Debug:     true,
	}, player.Nickname))
	if player.ClanTag != "" {
		nameBlocks = append(nameBlocks, render.NewTextContent(render.Style{
			FontColor: render.TextSecondary,
			Font:      &render.FontLarge,
			// Debug:     true,
		}, fmt.Sprintf("[%s]", player.ClanTag)))
	}
	return render.NewBlocksContent(render.Style{Direction: render.DirectionHorizontal, Gap: 5, AlignItems: render.AlignItemsCenter}, nameBlocks...)
}
