package wotblitz

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/cufee/aftermath-core/internal/core/utils"
)

var apiBaseUrl = utils.MustGetEnv("WOT_BLITZ_PUBLIC_API_URL_FMT")

var client = http.DefaultClient

func realmToSubdomain(realm string) string {
	switch strings.ToUpper(realm) {
	case "AS":
		return "asia"
	default:
		return strings.ToLower(realm)
	}
}

func apiUrl(realm string, endpoint string) string {
	return fmt.Sprintf(apiBaseUrl, realmToSubdomain(realm)) + endpoint
}
