package session

import (
	"github.com/cufee/aftermath-core/internal/core/database/models"
	"github.com/cufee/aftermath-core/internal/logic/dataprep"
	"github.com/cufee/aftermath-core/internal/logic/render"
)

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
