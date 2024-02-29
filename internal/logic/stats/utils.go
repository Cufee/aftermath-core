package stats

import (
	"github.com/cufee/aftermath-core/internal/core/stats"
	"github.com/cufee/am-wg-proxy-next/types"
)

func FrameToReducedStatsFrame(frame types.StatsFrame) *stats.ReducedStatsFrame {
	return &stats.ReducedStatsFrame{
		Battles:              frame.Battles,
		BattlesWon:           frame.Wins,
		BattlesSurvived:      frame.SurvivedBattles,
		DamageDealt:          frame.DamageDealt,
		DamageReceived:       frame.DamageReceived,
		ShotsHit:             frame.Hits,
		ShotsFired:           frame.Shots,
		Frags:                frame.Frags,
		MaxFrags:             frame.MaxFrags,
		EnemiesSpotted:       frame.Spotted,
		CapturePoints:        frame.CapturePoints,
		DroppedCapturePoints: frame.DroppedCapturePoints,
	}
}
