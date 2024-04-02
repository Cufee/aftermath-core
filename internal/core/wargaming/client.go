package wargaming

import (
	"time"

	"github.com/cufee/aftermath-core/internal/core/utils"
	"github.com/cufee/am-wg-proxy-next/v2/remote"
)

var Clients struct {
	Live  *remote.Client
	Cache *remote.Client
}

func init() {
	Clients.Live = remote.NewClient(utils.MustGetEnv("LIVE_WG_PROXY_URL"), time.Second*5)
	Clients.Cache = remote.NewClient(utils.MustGetEnv("CACHE_WG_PROXY_URL"), time.Second*30)
}
