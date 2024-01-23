package mock

import (
	"embed"
	"encoding/json"

	"github.com/cufee/aftermath-core/dataprep"
)

var PreviewStatsCards dataprep.SessionCards

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
