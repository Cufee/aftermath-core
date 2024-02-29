package period

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestGetPlayerStats(t *testing.T) {
	stats, err := GetPlayerStats(1013072123, 10)
	if err != nil {
		t.Fatal(err)
	}

	stats.Vehicles = nil

	data, _ := json.MarshalIndent(stats, "", "  ")
	fmt.Println(string(data))
}
