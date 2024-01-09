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

	session, err := stats.GetCurrentPlayerSession("na", 1013072123) // 1013379500 1013072123
	if err != nil {
		panic(err)
	}

	data, _ := json.MarshalIndent(session.Live.Global, "", "	")
	println(string(data))

	data, _ = json.MarshalIndent(session.Selected, "", "	")
	println(string(data))
}
