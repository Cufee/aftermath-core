package external

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/cufee/aftermath-core/internal/core/utils"
)

// Response from https://wotinspector.com/en/
type InspectorVehicle struct {
	NameEn  string `json:"en"`
	NameRu  string `json:"ru"`
	Tier    int    `json:"tier"`
	Type    int    `json:"type"`
	Premium int    `json:"premium"`
}

var inspectorVehiclesURL = utils.MustGetEnv("WOT_INSPECTOR_TANK_DB_URL")

func GetInspectorVehicles() (map[int]InspectorVehicle, error) {
	re := regexp.MustCompile(`(\d{1,9}):`)
	tanks := make(map[int]InspectorVehicle)

	res, err := insecureClient.Get(inspectorVehiclesURL)
	if err != nil || res == nil || res.StatusCode != http.StatusOK {
		return tanks, fmt.Errorf("status code: %+v. error: %s", res, err)
	}
	defer res.Body.Close()
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return tanks, fmt.Errorf("status code: %+v. error: %s", res, err)
	}
	tanksString := strings.ReplaceAll(string(bodyBytes), "TANK_DB = ", "")
	tanksString = re.ReplaceAllString(tanksString, `"$1":`)
	split := strings.SplitAfter(tanksString, "},")
	if len(split) <= 2 {
		return tanks, fmt.Errorf("failed to split string")
	}
	fix := strings.ReplaceAll(split[len(split)-2], "},", "}")
	tanksString = strings.ReplaceAll(tanksString, split[len(split)-2], fix)
	return tanks, json.Unmarshal([]byte(tanksString), &tanks)
}
