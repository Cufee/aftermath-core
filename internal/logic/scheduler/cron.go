package scheduler

import (
	"time"

	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"
)

func StartCronJobs() {
	log.Info().Msg("starting cron jobs")

	c := gocron.NewScheduler(time.UTC)
	// Tasks
	c.Cron("* * * * *").Do(runTasksWorker)
	// c.AddFunc("0 */2 * * *", func() { wrk.AutoRunTaskCleanup() })

	// Glossary
	c.Cron("20 9 * * *").Do(updateGlossaryWorker)
	c.Cron("10 9 * * *").Do(updateAveragesWorker)
	// c.AddFunc("40 9 * * 0", updateAchievementsWorker)

	// Sessions
	c.Cron("0 9 * * *").Do(createSessionTasksWorker("NA"))  // NA
	c.Cron("0 1 * * *").Do(createSessionTasksWorker("EU"))  // EU
	c.Cron("0 18 * * *").Do(createSessionTasksWorker("AS")) // Asia

	// Refresh WN8
	// "45 9 * * *" 	// NA
	// "45 1 * * *" 	// EU
	// "45 18 * * *" 	// Asia

	// Start the Cron job scheduler
	c.StartAsync()
}
