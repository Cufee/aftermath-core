package scheduler

import (
	"time"

	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"
)

var (
// SessionsNA =
)

func StartCronJobs() {
	log.Info().Msg("starting cron jobs")

	c := gocron.NewScheduler(time.UTC)
	// Tasks
	c.Cron("* * * * *").Do(runTasksWorker)
	c.Cron("0 * * * *").Do(restartTasksWorker)

	// Glossary - Do it around the same time WG releases game updates
	c.Cron("0 10 * * *").Do(updateGlossaryWorker)
	c.Cron("0 12 * * *").Do(updateGlossaryWorker)
	// c.AddFunc("40 9 * * 0", updateAchievementsWorker)

	// Averages - Update averages shortly after session refreshes
	c.Cron("0 10 * * *").Do(updateAveragesWorker)
	c.Cron("0 2 * * *").Do(updateAveragesWorker)
	c.Cron("0 19 * * *").Do(updateAveragesWorker)

	// Sessions
	c.Cron("0 9 * * *").Do(createSessionTasksWorker("NA"))  // NA
	c.Cron("0 1 * * *").Do(createSessionTasksWorker("EU"))  // EU
	c.Cron("0 18 * * *").Do(createSessionTasksWorker("AS")) // Asia

	// Refresh WN8
	// "45 9 * * *" 	// NA
	// "45 1 * * *" 	// EU
	// "45 18 * * *" 	// Asia

	// Configurations
	c.Cron("0 0 */7 * *").Do(rotateBackgroundPresetsWorker)

	// Start the Cron job scheduler
	c.StartAsync()
}
