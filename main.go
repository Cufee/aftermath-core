package main

import (
	"github.com/cufee/aftermath-core/internal/logic/server"
	"github.com/cufee/aftermath-core/internal/logic/tasks"
)

func main() {
	go tasks.StartCronJobs()
	server.Start()
}
