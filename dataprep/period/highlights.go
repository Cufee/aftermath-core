package period

import (
	"github.com/cufee/aftermath-core/dataprep"
	"github.com/cufee/aftermath-core/internal/core/stats"
	"github.com/rs/zerolog/log"
)

type highlight struct {
	compareWith dataprep.Tag
	blocks      []dataprep.Tag
	label       string
}

var (
	HighlightAvgDamage = highlight{dataprep.TagAvgDamage, []dataprep.Tag{dataprep.TagBattles, dataprep.TagAvgDamage, dataprep.TagWN8}, "label_highlight_avg_damage"}
	HighlightBattles   = highlight{dataprep.TagBattles, []dataprep.Tag{dataprep.TagBattles, dataprep.TagAvgDamage, dataprep.TagWN8}, "label_highlight_battles"}
	HighlightWN8       = highlight{dataprep.TagWN8, []dataprep.Tag{dataprep.TagBattles, dataprep.TagAvgDamage, dataprep.TagWN8}, "label_highlight_wn8"}
)

type highlightedVehicle struct {
	highlight highlight
	vehicle   *stats.ReducedVehicleStats
	value     float64
}

func getHighlightedVehicles(highlights []highlight, vehicles map[int]*stats.ReducedVehicleStats, minBattles int) []highlightedVehicle {
	leadersMap := make(map[string]highlightedVehicle)
	for _, vehicle := range vehicles {
		if vehicle.Battles < minBattles {
			continue
		}

		for _, highlight := range highlights {
			currentLeader := leadersMap[highlight.label]

			value, err := presetToBlock(highlight.compareWith, vehicle.ReducedStatsFrame, func(s string) string { return s })
			if err != nil {
				log.Warn().Str("highlight", highlight.label).Msg("failed to get preset value for a vehicle highlight")
				continue
			}

			if value.Data.Value > currentLeader.value {
				currentLeader.highlight = highlight
				currentLeader.value = value.Data.Value
				currentLeader.vehicle = vehicle
				leadersMap[highlight.label] = currentLeader
			}
		}
	}

	var highlightedVehicles []highlightedVehicle
	for _, leader := range leadersMap {
		highlightedVehicles = append(highlightedVehicles, leader)
	}
	return highlightedVehicles
}
