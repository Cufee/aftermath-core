package main

import (
	"os"

	"github.com/cufee/aftermath-core/internal/logic/server"
	"github.com/cufee/aftermath-core/internal/logic/tasks"
	"github.com/rs/zerolog"
)

func main() {
	level, _ := zerolog.ParseLevel(os.Getenv("LOG_LEVEL"))
	zerolog.SetGlobalLevel(level)

	tasks.StartCronJobs()
	server.Start()
}
