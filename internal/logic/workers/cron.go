package workers

import (
	"github.com/robfig/cron"
)

func StartCronJobs() {
	c := cron.New()
	// Tasks
	// c.AddFunc("* * * * *", func() { wrk.AutoRunTaskBatch() })
	// c.AddFunc("0 */2 * * *", func() { wrk.AutoRunTaskCleanup() })

	// Glossary
	// c.AddFunc("20 9 * * *", updateGlossaryWrkr)
	// c.AddFunc("10 9 * * *", updateAveragesWrkr)
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
