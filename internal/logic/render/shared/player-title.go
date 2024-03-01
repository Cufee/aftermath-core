package shared

import (
	"image/color"

	"github.com/cufee/aftermath-core/internal/core/database/models"
	"github.com/cufee/aftermath-core/internal/logic/render"
	"github.com/cufee/aftermath-core/internal/logic/render/badges"
	"github.com/rs/zerolog/log"
)

type TitleCardStyle struct {
	Container render.Style
	Nickname  render.Style
	ClanTag   render.Style
}

func NewPlayerTitleCard(style TitleCardStyle, nickname, clanTag string, subscriptions []models.UserSubscription) render.Block {
	clanTagBlock, hasClanTagBlock := newClanTagBlock(style.ClanTag, clanTag, subscriptions)
	if !hasClanTagBlock {
		return render.NewBlocksContent(style.Container, render.NewTextContent(style.Nickname, nickname))
	}

	content := make([]render.Block, 0, 3)
	style.Container.JustifyContent = render.JustifyContentSpaceBetween

	clanTagImage, err := clanTagBlock.Render()
	if err != nil {
		log.Warn().Err(err).Msg("failed to render clan tag")
		// This error is not fatal, we can just render the name
		return render.NewBlocksContent(style.Container, render.NewTextContent(style.Nickname, nickname))
	}
	content = append(content, render.NewImageContent(render.Style{Width: float64(clanTagImage.Bounds().Dx()), Height: float64(clanTagImage.Bounds().Dy())}, clanTagImage))

	// Nickname
	content = append(content, render.NewTextContent(style.Nickname, nickname))

	// Invisible tag to offset the nickname
	clanTagOffsetBlock := render.NewBlocksContent(render.Style{
		Width:          float64(clanTagImage.Bounds().Dx()),
		JustifyContent: render.JustifyContentEnd,
	}, render.NewTextContent(render.Style{Font: style.ClanTag.Font, FontColor: color.Transparent}, "-"))
	content = append(content, clanTagOffsetBlock)

	return render.NewBlocksContent(style.Container, content...)

}

func newClanTagBlock(style render.Style, clanTag string, subs []models.UserSubscription) (render.Block, bool) {
	if clanTag == "" {
		return render.Block{}, false
	}

	var blocks []render.Block
	blocks = append(blocks, render.NewTextContent(render.Style{Font: &render.FontMedium, FontColor: render.TextSecondary}, clanTag))
	if sub := badges.ClanSubscriptionsBadges(subs); sub != nil {
		iconBlock, err := sub.Block()
		if err == nil {
			blocks = append(blocks, iconBlock)
		}
	}

	return render.NewBlocksContent(style, blocks...), true
}
