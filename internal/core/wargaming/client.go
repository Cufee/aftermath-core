package wargaming

import (
	"time"

	"github.com/cufee/aftermath-core/internal/core/utils"
	client "github.com/cufee/am-wg-proxy-next/client"
)

var Clients struct {
	Live  *client.Client
	Cache *client.Client
}

func init() {
	Clients.Live = client.NewClient(utils.MustGetEnv("LIVE_WG_PROXY_URL"), time.Second*5)
	Clients.Cache = client.NewClient(utils.MustGetEnv("CACHE_WG_PROXY_URL"), time.Second*30)
}
