package scheduler

import (
	"github.com/cufee/aftermath-core/internal/logic/cache"
	"github.com/cufee/aftermath-core/internal/logic/scheduler/tasks"
	"github.com/rs/zerolog/log"
)

func updateGlossaryWorker() {
	// We just run the logic directly as it's not a heavy task and it doesn't matter if it fails due to the app failing
	log.Info().Msg("updating glossary cache")
	err := cache.UpdateGlossaryCache()
	if err != nil {
		log.Err(err).Msg("failed to update glossary cache")
	} else {
		log.Info().Msg("glossary cache updated")
	}
}

func updateAveragesWorker() {
	// We just run the logic directly as it's not a heavy task and it doesn't matter if it fails due to the app failing
	log.Info().Msg("updating averages cache")
	err := cache.UpdateAveragesCache()
	if err != nil {
		log.Err(err).Msg("failed to update averages cache")
	} else {
		log.Info().Msg("averages cache updated")
	}
}

func recordSessionsWorker(realm string) func() {
	return func() {
		err := tasks.CreateSessionUpdateTasks(realm)
		if err != nil {
			log.Err(err).Msg("failed to create session update tasks")
		}
	}
}