package mock

import (
	"embed"
	"encoding/json"

	session "github.com/cufee/aftermath-core/dataprep/session"
)

var PreviewStatsCards session.Cards

//go:embed mock.json
var mockData embed.FS

func init() {
	file, err := mockData.ReadFile("mock.json")
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(file, &PreviewStatsCards)
	if err != nil {
		panic(err)
	}
}
