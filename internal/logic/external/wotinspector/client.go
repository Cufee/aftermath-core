package wotinspector

import (
	"crypto/tls"
	"net/http"
	"time"

	"github.com/cufee/aftermath-core/internal/core/utils"
)

var replayUploadUrl = utils.MustGetEnv("WOT_INSPECTOR_REPLAYS_URL")

var insecureClient = &http.Client{Timeout: 10 * time.Second, Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
