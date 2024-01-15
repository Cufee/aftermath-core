package tasks

import (
	"github.com/cufee/aftermath-core/internal/logic/cache"
	"github.com/rs/zerolog/log"
)

func updateGlossaryWorker() {
	log.Info().Msg("updating glossary cache")
	err := cache.UpdateGlossaryCache()
	if err != nil {
		log.Err(err).Msg("failed to update glossary cache")
	}
}

func updateAveragesWorker() {
	log.Info().Msg("updating averages cache")
	err := cache.UpdateAveragesCache()
	if err != nil {
		log.Err(err).Msg("failed to update averages cache")
	}
}
