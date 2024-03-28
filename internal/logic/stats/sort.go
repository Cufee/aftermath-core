package stats

import (
	"sort"

	core "github.com/cufee/aftermath-core/internal/core/stats"
)

type SortOptions struct {
	By    vehicleSortOptions
	Limit int
}
type vehicleSortOptions string

const (
	SortByBattlesDesc   = vehicleSortOptions("-battles")
	SortByBattlesAsc    = vehicleSortOptions("battles")
	SortByWinrateDesc   = vehicleSortOptions("-winrate")
	SortByWinrateAsc    = vehicleSortOptions("winrate")
	SortByWN8Desc       = vehicleSortOptions("-wn8")
	SortByWN8Asc        = vehicleSortOptions("wn8")
	SortByAvgDamageDesc = vehicleSortOptions("-avgDamage")
	SortByAvgDamageAsc  = vehicleSortOptions("avgDamage")
	SortByLastBattle    = vehicleSortOptions("lastBattleTime")
)

func SortVehicles(vehicles map[int]core.ReducedVehicleStats, averages map[int]core.ReducedStatsFrame, options ...SortOptions) []core.ReducedVehicleStats {
	opts := SortOptions{By: SortByLastBattle, Limit: 10}
	if len(options) > 0 {
		opts = options[0]
	}

	var sorted []core.ReducedVehicleStats
	for _, vehicle := range vehicles {
		sorted = append(sorted, vehicle)
	}

	switch opts.By {
	case SortByBattlesDesc:
		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i].Battles > sorted[j].Battles
		})
	case SortByBattlesAsc:
		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i].Battles < sorted[j].Battles
		})
	case SortByWinrateDesc:
		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i].Winrate() > sorted[j].Winrate()
		})
	case SortByWinrateAsc:
		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i].Winrate() < sorted[j].Winrate()
		})
	case SortByWN8Desc:
		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i].WN8(averages[sorted[i].VehicleID]) > sorted[j].WN8(averages[sorted[j].VehicleID])
		})
	case SortByWN8Asc:
		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i].WN8(averages[sorted[i].VehicleID]) < sorted[j].WN8(averages[sorted[j].VehicleID])
		})
	case SortByAvgDamageDesc:
		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i].AvgDamage() > sorted[j].AvgDamage()
		})
	case SortByAvgDamageAsc:
		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i].AvgDamage() < sorted[j].AvgDamage()
		})
	case SortByLastBattle:
		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i].LastBattleTime > sorted[j].LastBattleTime
		})
	}

	if opts.Limit > 0 && opts.Limit < len(sorted) {
		sorted = sorted[:opts.Limit]
	}
	return sorted
}

func ParseSortOptions(sort string) vehicleSortOptions {
	switch sort {
	case string(SortByBattlesDesc):
		return SortByBattlesDesc
	case string(SortByBattlesAsc):
		return SortByBattlesAsc
	case string(SortByWinrateDesc):
		return SortByWinrateDesc
	case string(SortByWinrateAsc):
		return SortByWinrateAsc
	case string(SortByWN8Desc):
		return SortByWN8Desc
	case string(SortByWN8Asc):
		return SortByWN8Asc
	case string(SortByAvgDamageDesc):
		return SortByAvgDamageDesc
	case string(SortByAvgDamageAsc):
		return SortByAvgDamageAsc
	case string(SortByLastBattle):
		return SortByLastBattle
	default:
		return SortByLastBattle
	}
}
