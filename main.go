package main

import (
	"github.com/cufee/aftermath-core/internal/logic/server"
)

func main() {
	// err := cache.RefreshSessions(cache.SessionTypeDaily, "na", 1013072123)
	// if err != nil {
	// 	panic(err)
	// }

	server.Start()
}
