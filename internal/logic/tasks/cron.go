package tasks

import (
	"time"

	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"
)

func StartCronJobs() {
	log.Info().Msg("starting cron jobs")

	c := gocron.NewScheduler(time.UTC)
	// Tasks
	// c.AddFunc("* * * * *", func() { wrk.AutoRunTaskBatch() })
	// c.AddFunc("0 */2 * * *", func() { wrk.AutoRunTaskCleanup() })

	// Glossary
	c.Cron("20 9 * * *").Do(updateGlossaryWorker)
	c.Cron("10 9 * * *").Do(updateAveragesWorker)
	// c.AddFunc("40 9 * * 0", updateAchievementsWrkr)

	// Sessions
	// c.AddFunc("0 9 * * *", func() { recordSessionsWrkr("NA") })    // NA
	// c.AddFunc("0 1 * * *", func() { recordSessionsWrkr("EU") })    // EU
	// c.AddFunc("0 23 * * *", func() { recordSessionsWrkr("RU") })   // RU
	// c.AddFunc("0 18 * * *", func() { recordSessionsWrkr("ASIA") }) // ASIA

	// Refresh WN8
	// c.AddFunc("45 9 * * *", func() { refreshRatingWrkr("NA") })    // NA
	// c.AddFunc("45 1 * * *", func() { refreshRatingWrkr("EU") })    // EU
	// c.AddFunc("45 23 * * *", func() { refreshRatingWrkr("RU") })   // RU
	// c.AddFunc("45 18 * * *", func() { refreshRatingWrkr("ASIA") }) // ASIA

	// Start the Cron job scheduler
	c.StartAsync()
}

// func updateSessionsWorker(realm string) func() {
// 	return func() {
// 		err := cache.RefreshSessions(cache.SessionTypeDaily, realm)
// 		if err != nil {
// 			log.Errorf("failed to update daily sessions: %s", err.Error())
// 		}
// 	}
// }
