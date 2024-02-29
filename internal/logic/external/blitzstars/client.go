package blitzstars

import (
	"crypto/tls"
	"net/http"
	"time"

	"github.com/cufee/aftermath-core/internal/core/utils"
)

var starsStatsApiURL = utils.MustGetEnv("BLITZ_STARS_API_URL")

var insecureClient = &http.Client{Timeout: 10 * time.Second, Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
