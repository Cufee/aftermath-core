package main

import (
	"os"

	"github.com/cufee/aftermath-core/internal/logic/scheduler"
	"github.com/cufee/aftermath-core/internal/logic/server"
	"github.com/rs/zerolog"
)

func main() {
	level, _ := zerolog.ParseLevel(os.Getenv("LOG_LEVEL"))
	zerolog.SetGlobalLevel(level)

	scheduler.StartCronJobs()
	server.Start()
}
