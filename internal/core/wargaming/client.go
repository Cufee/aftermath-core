package wargaming

import (
	"time"

	"github.com/cufee/aftermath-core/internal/core/utils"

	"github.com/cufee/am-wg-proxy-next/v2/client"
)

var Clients struct {
	Live  client.Client
	Cache client.Client
}

func init() {
	Clients.Live, _ = client.NewRemoteClient(utils.MustGetEnv("LIVE_WG_PROXY_URL"), time.Second*5)
	Clients.Cache, _ = client.NewRemoteClient(utils.MustGetEnv("CACHE_WG_PROXY_URL"), time.Second*30)
}
