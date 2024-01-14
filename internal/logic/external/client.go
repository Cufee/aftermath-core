package external

import (
	"crypto/tls"
	"net/http"
	"time"
)

var insecureClient = &http.Client{Timeout: 10 * time.Second, Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
