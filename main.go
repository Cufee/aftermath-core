package main

import (
	"github.com/cufee/aftermath-core/internal/logic/cache"
	"github.com/cufee/aftermath-core/internal/logic/server"
	"github.com/cufee/aftermath-core/internal/logic/tasks"
	"github.com/gofiber/fiber/v2/log"
)

func main() {
	err := cache.UpdateGlossaryCache()
	if err != nil {
		log.Errorf("failed to update glossary cache: %s", err.Error())
	}

	err = cache.UpdateAveragesCache()
	if err != nil {
		log.Errorf("failed to update averages cache: %s", err.Error())
	}

	go tasks.StartCronJobs()
	server.Start()
}
