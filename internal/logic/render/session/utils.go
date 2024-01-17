package session

import (
	"github.com/cufee/aftermath-core/internal/logic/dataprep"
	"github.com/cufee/aftermath-core/internal/logic/render"
	"github.com/cufee/aftermath-core/internal/logic/users"
)

func (data *PlayerData) userSubscriptionHeader() *subscriptionHeader {
	for _, subscription := range data.Subscriptions {
		switch subscription.Type {
		case users.SubscriptionTypePro:
			return userSubscriptionPro
		case users.SubscriptionTypePlus:
			return userSubscriptionPlus
		case users.SubscriptionTypeSupporter:
			return userSubscriptionSupporter
		}
	}
	return nil
}

func (data *PlayerData) clanSubscriptionHeader() *subscriptionHeader {
	for _, subscription := range data.Subscriptions {
		switch subscription.Type {
		case users.SubscriptionTypeProClan:
			return clanSubscriptionPro
		case users.SubscriptionTypeSupporterClan:
			return clanSubscriptionSupporter
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
