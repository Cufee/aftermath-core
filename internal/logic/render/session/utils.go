package session

import (
	"image"
	"image/color"

	"github.com/cufee/aftermath-core/dataprep"
	"github.com/cufee/aftermath-core/internal/core/database/models"
	"github.com/cufee/aftermath-core/internal/core/stats"
	"github.com/cufee/aftermath-core/internal/logic/render"
	"github.com/cufee/aftermath-core/internal/logic/render/assets"
	"github.com/fogleman/gg"
)

var wn8Icon image.Image

func init() {
	ctx := gg.NewContext(20, 20)
	// ctx.DrawCircle(7.5, 10, 7.5)
	ctx.DrawRoundedRectangle(7, 0, 7, 20, 3.5)
	ctx.SetColor(color.RGBA{R: 255, G: 255, B: 255, A: 255})
	ctx.Fill()
	wn8Icon = ctx.Image()
}

func (data *PlayerData) userSubscriptionHeader() *subscriptionHeader {
	for _, subscription := range data.Subscriptions {
		switch subscription.Type {
		case models.SubscriptionTypePro:
			return userSubscriptionPro
		case models.SubscriptionTypePlus:
			return userSubscriptionPlus
		case models.SubscriptionTypeSupporter:
			return userSubscriptionSupporter
		}
	}
	return nil
}

func (data *PlayerData) clanSubscriptionHeader() *subscriptionHeader {
	for _, subscription := range data.Subscriptions {
		switch subscription.Type {
		case models.SubscriptionTypeProClan:
			return clanSubscriptionPro
		case models.SubscriptionTypeVerifiedClan:
			return clanSubscriptionVerified
		}
	}
	return nil
}

func styleBlocks(blocks []dataprep.StatsBlock, styles ...render.Style) []styledStatsBlock {
	var lastStyle render.Style
	var styledBlocks []styledStatsBlock
	for i, block := range blocks {
		if i < len(styles) {
			lastStyle = styles[i]
		}
		styledBlocks = append(styledBlocks, styledStatsBlock{StatsBlock: block, style: lastStyle})
	}
	return styledBlocks
}

func getWN8Color(r int) color.Color {
	if r > 0 && r < 301 {
		return color.RGBA{255, 0, 0, 180}
	}
	if r > 300 && r < 451 {
		return color.RGBA{251, 83, 83, 180}
	}
	if r > 450 && r < 651 {
		return color.RGBA{255, 160, 49, 180}
	}
	if r > 650 && r < 901 {
		return color.RGBA{255, 244, 65, 180}
	}
	if r > 900 && r < 1201 {
		return color.RGBA{149, 245, 62, 180}
	}
	if r > 1200 && r < 1601 {
		return color.RGBA{103, 190, 51, 180}
	}
	if r > 1600 && r < 2001 {
		return color.RGBA{106, 236, 255, 180}
	}
	if r > 2000 && r < 2451 {
		return color.RGBA{46, 174, 193, 180}
	}
	if r > 2450 && r < 2901 {
		return color.RGBA{208, 108, 255, 180}
	}
	if r > 2900 {
		return color.RGBA{142, 65, 177, 180}
	}
	return color.Transparent
}

type comparisonIcon struct {
	left  render.Block
	right render.Block
}

func comparisonIconFromBlock(block dataprep.StatsBlock) *comparisonIcon {
	if block.Tag == dataprep.TagBattles {
		// Don't show comparison icons for battle count
		return nil
	}
	if !stats.ValueValid(block.Session.Value) || !stats.ValueValid(block.Career.Value) {
		return nil
	}

	if block.Tag == dataprep.TagWN8 {
		// WN8 icons need to show the color
		return blockToWN8Icon(block.Session, block.Tag)
	}

	var icon image.Image
	var iconColor color.Color
	if block.Session.Value > block.Career.Value {
		icon, _ = assets.GetImage("images/icons/chevron-up-single")
		iconColor = color.RGBA{R: 0, G: 255, B: 0, A: 255}
	}
	if block.Session.Value < block.Career.Value {
		icon, _ = assets.GetImage("images/icons/chevron-down-single")
		iconColor = color.RGBA{R: 255, G: 0, B: 0, A: 255}
	}
	if icon == nil {
		return nil
	}
	return &comparisonIcon{
		left:  render.NewImageContent(render.Style{Width: 20, Height: 20, BackgroundColor: iconColor}, icon),
		right: render.NewImageContent(render.Style{Width: 20, Height: 20, BackgroundColor: color.Transparent}, icon),
	}
}

func blockToWN8Icon(value dataprep.Value, tag dataprep.Tag) *comparisonIcon {
	if tag != dataprep.TagWN8 || !stats.ValueValid(value.Value) {
		return nil
	}
	return &comparisonIcon{
		left:  render.NewImageContent(render.Style{Width: 20, Height: 20, BackgroundColor: getWN8Color(int(value.Value))}, wn8Icon),
		right: render.NewImageContent(render.Style{Width: 20, Height: 20, BackgroundColor: color.Transparent}, wn8Icon),
	}
}
