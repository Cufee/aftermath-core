package main

import (
	"os"

	"github.com/cufee/aftermath-core/internal/core/database"
	"github.com/cufee/aftermath-core/internal/core/utils"
	"github.com/cufee/aftermath-core/internal/logic/scheduler"
	"github.com/cufee/aftermath-core/internal/logic/server"
	"github.com/rs/zerolog"
)

func main() {
	level, _ := zerolog.ParseLevel(os.Getenv("LOG_LEVEL"))
	zerolog.SetGlobalLevel(level)

	err := database.Connect(utils.MustGetEnv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}

	err = database.SyncIndexes(database.DefaultClient.Database())
	if err != nil {
		panic(err)
	}

	if os.Getenv("SCHEDULER_ENABLED") != "false" {
		scheduler.StartCronJobs()
	}

	server.Start()
}
