package stats

import (
	"strings"
	"time"

	core "github.com/cufee/aftermath-core/internal/core/stats"
	"github.com/gorhill/cronexpr"
)

type PeriodStats struct {
	Account AccountWithClan

	Start time.Time `json:"start"`
	End   time.Time `json:"end"`

	Stats    *core.ReducedStatsFrame
	Vehicles map[int]core.ReducedVehicleStats
}

var sessionsCronNA = cronexpr.MustParse("0 9 * * *")
var sessionsCronEU = cronexpr.MustParse("0 1 * * *")
var sessionsCronAsia = cronexpr.MustParse("0 18 * * *")

func GetPlayerStats(accountId int, days int) (*PeriodStats, error) {

	return nil, nil
}

func daysToRealmTime(realm string, days int) time.Time {
	switch strings.ToLower(realm) {
	case "na":
		return sessionsCronNA.Next(time.Now()).Add(time.Hour * 24 * -1)
	case "eu":
		return sessionsCronEU.Next(time.Now()).Add(time.Hour * 24 * -1)
	case "as":
		return sessionsCronAsia.Next(time.Now()).Add(time.Hour * 24 * -1)
	default:
		return time.Now()
	}
}
