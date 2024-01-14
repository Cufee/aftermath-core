package tasks

import (
	"github.com/cufee/aftermath-core/internal/logic/cache"
	"github.com/gofiber/fiber/v2/log"
	"github.com/robfig/cron"
)

func StartCronJobs() {
	log.Info("starting cron jobs")

	c := cron.New()
	// Tasks
	// c.AddFunc("* * * * *", func() { wrk.AutoRunTaskBatch() })
	// c.AddFunc("0 */2 * * *", func() { wrk.AutoRunTaskCleanup() })

	// Glossary
	c.AddFunc("20 9 * * *", updateGlossaryWorker)
	c.AddFunc("10 9 * * *", updateAveragesWorker)
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
	c.Start()
}

func updateGlossaryWorker() {
	err := cache.UpdateGlossaryCache()
	if err != nil {
		log.Errorf("failed to update glossary cache: %s", err.Error())
	}
}

func updateAveragesWorker() {
	err := cache.UpdateAveragesCache()
	if err != nil {
		log.Errorf("failed to update averages cache: %s", err.Error())
	}
}

// func updateSessionsWorker(realm string) func() {
// 	return func() {
// 		err := cache.RefreshSessions(cache.SessionTypeDaily, realm)
// 		if err != nil {
// 			log.Errorf("failed to update daily sessions: %s", err.Error())
// 		}
// 	}
// }
