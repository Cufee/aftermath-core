package session

import (
	"image"
	"image/color"
	"slices"

	"github.com/cufee/aftermath-core/dataprep"
	"github.com/cufee/aftermath-core/dataprep/session"
	"github.com/cufee/aftermath-core/internal/core/database/models"
	"github.com/cufee/aftermath-core/internal/core/stats"
	"github.com/cufee/aftermath-core/internal/logic/render"
	"github.com/cufee/aftermath-core/internal/logic/render/assets"
)

func (data *PlayerData) userBadges() ([]render.Block, error) {
	slices.SortFunc(data.Subscriptions, func(i, j models.UserSubscription) int {
		return subscriptionWeight[j.Type] - subscriptionWeight[i.Type]
	})

	var badges []render.Block
	for _, subscription := range data.Subscriptions {
		var header *subscriptionHeader
		switch subscription.Type {
		case models.SubscriptionTypeDeveloper:
			header = subscriptionDeveloper
		case models.SubscriptionTypeServerModerator:
			header = subscriptionServerModerator
		case models.SubscriptionTypeContentTranslator:
			header = subscriptionTranslator
		}

		if header != nil {
			block, err := header.Block()
			if err != nil {
				return nil, err
			}
			badges = append(badges, block)
			break
		}
	}
	for _, subscription := range data.Subscriptions {
		var header *subscriptionHeader
		switch subscription.Type {
		case models.SubscriptionTypeContentTranslator:
			header = subscriptionTranslator
		case models.SubscriptionTypeSupporter:
			header = userSubscriptionSupporter
		}

		if header != nil {
			block, err := header.Block()
			if err != nil {
				return nil, err
			}
			badges = append(badges, block)
			break
		}
	}
	for _, subscription := range data.Subscriptions {
		var header *subscriptionHeader
		switch subscription.Type {
		case models.SubscriptionTypePro:
			header = userSubscriptionPro
		case models.SubscriptionTypePlus:
			header = userSubscriptionPlus
		case models.SubscriptionTypeSupporter:
			header = userSubscriptionSupporter
		case models.SubscriptionTypeServerBooster:
			header = subscriptionServerBooster
		}

		if header != nil {
			block, err := header.Block()
			if err != nil {
				return nil, err
			}
			badges = append(badges, block)
			break
		}
	}

	return badges, nil
}

func (data *PlayerData) clanSubscriptionHeader() *subscriptionHeader {
	var headers []*subscriptionHeader

	for _, subscription := range data.Subscriptions {
		switch subscription.Type {
		case models.SubscriptionTypeProClan:
			headers = append(headers, clanSubscriptionPro)
		case models.SubscriptionTypeVerifiedClan:
			headers = append(headers, clanSubscriptionVerified)
		}
	}

	if len(headers) > 0 {
		return headers[0]
	}

	return nil
}

func styleBlocks(blocks []session.StatsBlock, styles ...render.Style) []styledStatsBlock {
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

func comparisonIconFromBlock(block session.StatsBlock) render.Block {
	if !stats.ValueValid(block.Session.Value) || !stats.ValueValid(block.Career.Value) {
		return blankIconBlock
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
		return blankIconBlock
	}

	return render.NewImageContent(render.Style{Width: 25, Height: 25, BackgroundColor: iconColor}, icon)
}

func blockToWN8Icon(value dataprep.Value, tag dataprep.Tag) render.Block {
	if tag != dataprep.TagWN8 || !stats.ValueValid(value.Value) {
		return blankIconBlock
	}
	return render.NewImageContent(render.Style{Width: 25, Height: 25, BackgroundColor: getWN8Color(int(value.Value))}, wn8Icon)
}
