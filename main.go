package main

import (
	"encoding/json"

	"github.com/cufee/aftermath-core/internal/logic/stats"
)

func main() {
	// err := cache.RefreshSessions(cache.SessionTypeDaily, "na", 1013072123)
	// if err != nil {
	// 	panic(err)
	// }
	session, err := stats.GetCurrentPlayerSession("na", 1013072123)
	if err != nil {
		panic(err)
	}

	data, _ := json.MarshalIndent(session, "", "	")
	println(string(data))
}
