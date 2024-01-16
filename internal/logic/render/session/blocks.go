package session

import (
	"errors"

	"github.com/cufee/aftermath-core/internal/core/localization"
	core "github.com/cufee/aftermath-core/internal/core/stats"
	"github.com/cufee/aftermath-core/internal/logic/render"
	"github.com/rs/zerolog/log"
)

func FrameToOverviewBlocks(cardStyle render.Style, session, career *core.ReducedStatsFrame, sessionWN8, careerWN8 int, localePrinter localization.LocalePrinter) ([]render.Block, error) {
	if session == nil {
		return nil, errors.New("session is nil")
	}

	var blocks []render.Block
	{
		// Battles
		values := []interface{}{session.Battles}
		if career != nil {
			values = append(values, career.Battles)
		}
		battlesBlock := NewStatsBlock("", values...)
		// Add some special styling to the block
		battlesBlock.Style.PaddingY = 10
		battlesBlock.Style.BorderRadius = 10
		battlesBlock.Style.BackgroundColor = HighlightCardColor(cardStyle.BackgroundColor)
		blocks = append(blocks, battlesBlock)
	}
	{
		// Avg Damage
		values := []any{int(session.AvgDamage())}
		if career != nil {
			values = append(values, int(career.AvgDamage()))
		}
		blocks = append(blocks, NewStatsBlock(localePrinter("label_avg_damage"), values...))
	}
	{
		// Damage Ratio
		values := []any{session.DamageRatio()}
		if career != nil {
			values = append(values, career.DamageRatio())
		}
		blocks = append(blocks, NewStatsBlock(localePrinter("label_damage_ratio"), values...))
	}
	{
		// Winrate
		values := []any{session.Winrate()}
		if career != nil {
			values = append(values, career.Winrate())
		}
		blocks = append(blocks, NewStatsBlock(localePrinter("label_winrate"), values...))
	}
	{
		if sessionWN8 != core.InvalidValue {
			// WN8
			values := []any{sessionWN8}
			if careerWN8 != core.InvalidValue {
				values = append(values, careerWN8)
			}
			blocks = append(blocks, NewStatsBlock(localePrinter("label_wn8"), values...))
		} else {
			// Fallback to Accuracy to keep the UI consistent
			values := []any{session.Accuracy()}
			if career != nil {
				values = append(values, career.Accuracy())
			}
			blocks = append(blocks, NewStatsBlock(localePrinter("label_accuracy"), values...))
		}
	}

	return blocks, nil
}

func FrameToStatsBlocks(cardStyle render.Style, session, career, averages *core.ReducedStatsFrame, localePrinter localization.LocalePrinter) ([]render.Block, error) {
	if session == nil {
		return nil, errors.New("session is nil")
	}

	var blocks []render.Block
	{
		// Battles
		values := []interface{}{session.Battles}
		if career != nil {
			values = append(values, career.Battles)
		}
		battlesBlock := NewStatsBlock("", values...)
		// Add some special styling to the block
		battlesBlock.Style.PaddingY = 10
		battlesBlock.Style.BorderRadius = 10
		battlesBlock.Style.BackgroundColor = HighlightCardColor(cardStyle.BackgroundColor)
		blocks = append(blocks, battlesBlock)
	}
	{
		// Avg Damage
		values := []interface{}{int(session.AvgDamage())}
		if career != nil {
			values = append(values, int(career.AvgDamage()))
		}
		blocks = append(blocks, NewStatsBlock(localePrinter("label_avg_damage"), values...))
	}
	{
		// Damage Ratio
		values := []any{session.DamageRatio()}
		if career != nil {
			values = append(values, career.DamageRatio())
		}
		blocks = append(blocks, NewStatsBlock(localePrinter("label_damage_ratio"), values...))
	}
	{
		// Winrate
		values := []interface{}{session.Winrate()}
		if career != nil {
			values = append(values, career.Winrate())
		}
		blocks = append(blocks, NewStatsBlock(localePrinter("label_winrate"), values...))
	}
	{
		if session.WN8(averages) != core.InvalidValue {
			// WN8
			values := []interface{}{int(session.WN8(averages))}
			if career != nil {
				values = append(values, int(career.WN8(averages)))
			}
			blocks = append(blocks, NewStatsBlock(localePrinter("label_wn8"), values...))
		} else {
			// Fallback to Accuracy to keep the UI consistent
			values := []interface{}{session.Accuracy()}
			if career != nil {
				values = append(values, career.Accuracy())
			}
			blocks = append(blocks, NewStatsBlock(localePrinter("label_accuracy"), values...))
		}
	}

	return blocks, nil
}

func FrameToSlimStatsBlocks(cardStyle render.Style, session, averages *core.ReducedStatsFrame, localePrinter localization.LocalePrinter) ([]render.Block, error) {
	if session == nil {
		return nil, errors.New("session is nil")
	}

	var blocks []render.Block
	{
		// Battles
		battlesBlock := NewStatsBlock("", session.Battles)
		// Add some special styling to the block
		battlesBlock.Style.PaddingY = 10
		battlesBlock.Style.BorderRadius = 10
		battlesBlock.Style.BackgroundColor = HighlightCardColor(cardStyle.BackgroundColor)
		blocks = append(blocks, battlesBlock)
	}
	{
		// Avg Damage
		blocks = append(blocks, NewStatsBlock(localePrinter("label_avg_damage"), int(session.AvgDamage())))
	}
	{
		// Damage Ratio
		blocks = append(blocks, NewStatsBlock(localePrinter("label_damage_ratio"), session.DamageRatio()))
	}
	{
		// Winrate
		blocks = append(blocks, NewStatsBlock(localePrinter("label_winrate"), session.Winrate()))
	}
	{
		if session.WN8(averages) != core.InvalidValue {
			// WN8
			blocks = append(blocks, NewStatsBlock(localePrinter("label_wn8"), session.WN8(averages)))
		} else {
			// Fallback to Accuracy to keep the UI consistent
			blocks = append(blocks, NewStatsBlock(localePrinter("label_accuracy"), session.Accuracy()))
		}
	}

	return blocks, nil
}

func NewPlayerTitleCard(style render.Style, name, clanTag string, clanSubHeader render.Block) render.Block {
	if clanTag == "" {
		return render.NewBlocksContent(style, render.NewTextContent(render.Style{Font: &FontLarge, FontColor: FontLargeColor}, name))
	}

	content := make([]render.Block, 0, 3)
	style.JustifyContent = render.JustifyContentSpaceBetween

	// Visible tag
	clanTagBlock := render.NewBlocksContent(render.Style{
		Direction:       render.DirectionHorizontal,
		AlignItems:      render.AlignItemsCenter,
		PaddingX:        10,
		PaddingY:        5,
		BackgroundColor: HighlightCardColor(style.BackgroundColor),
		BorderRadius:    10,
		// Debug:           true,
	}, render.NewTextContent(render.Style{Font: &FontMedium, FontColor: FontMediumColor}, clanTag))

	clanTagImage, err := clanTagBlock.Render()
	if err != nil {
		log.Warn().Err(err).Msg("failed to render clan tag")
		// This error is not fatal, we can just render the name
		return render.NewBlocksContent(style, render.NewTextContent(render.Style{Font: &FontLarge, FontColor: FontLargeColor}, name))
	}
	content = append(content, render.NewImageContent(render.Style{Width: float64(clanTagImage.Bounds().Dx()), Height: float64(clanTagImage.Bounds().Dy())}, clanTagImage))

	// Nickname
	content = append(content, render.NewTextContent(render.Style{Font: &FontLarge, FontColor: FontLargeColor}, name))

	clanBlock := render.NewBlocksContent(render.Style{
		Width:          float64(clanTagImage.Bounds().Dx()),
		JustifyContent: render.JustifyContentEnd,
	}, clanSubHeader)

	content = append(content, clanBlock)

	return render.NewBlocksContent(style, render.NewBlocksContent(render.Style{
		JustifyContent: render.JustifyContentSpaceBetween,
		Direction:      render.DirectionHorizontal,
		AlignItems:     render.AlignItemsCenter,
		Width:          BaseCardWidth,
		PaddingX:       20,
		// Debug:          true,
	}, content...))
}

func NewStatsBlock(label string, values ...any) render.Block {
	var content []render.Block
	for i, value := range values {
		style := render.Style{Font: &FontLarge, FontColor: FontLargeColor}
		if i > 0 {
			style = render.Style{Font: &FontMedium, FontColor: FontMediumColor}
		}
		content = append(content, render.NewTextContent(style, statsValueToString(value)))
	}
	if label != "" {
		content = append(content, render.NewTextContent(render.Style{Font: &FontSmall, FontColor: FontSmallColor}, label))
	}
	return render.NewBlocksContent(render.Style{
		Direction:  render.DirectionVertical,
		AlignItems: render.AlignItemsCenter,
		Width:      BaseStatsBlockWidth,
		// Debug:      true,
	}, content...)
}

func NewTextLabel(label string) render.Block {
	return render.NewTextContent(render.Style{Font: &FontMedium, FontColor: FontMediumColor}, label)
}

func NewVehicleLabel(name, tier string) render.Block {
	var blocks []render.Block
	if tier != "" {
		blocks = append(blocks, render.NewTextContent(render.Style{Font: &FontMedium, FontColor: FontMediumColor}, tier))
	}
	blocks = append(blocks, render.NewTextContent(render.Style{Font: &FontMedium, FontColor: FontMediumColor}, name))

	return render.NewBlocksContent(
		render.Style{
			Direction:  render.DirectionHorizontal,
			AlignItems: render.AlignItemsCenter,
			Gap:        5,
			// Debug:      true,
		},
		blocks...,
	)
}

func NewCardBlock(cardStyle render.Style, label render.Block, stats []render.Block) render.Block {
	var content []render.Block
	content = append(content, label)
	content = append(content, render.NewBlocksContent(render.Style{
		Direction:      render.DirectionHorizontal,
		JustifyContent: render.JustifyContentSpaceBetween,
		Gap:            10,
		// Debug:     true,
	}, stats...))

	return render.NewBlocksContent(cardStyle, content...)
}
