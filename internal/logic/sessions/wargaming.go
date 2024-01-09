package sessions

import (
	"time"

	"github.com/cufee/aftermath-core/internal/core/utils"
	client "github.com/cufee/am-wg-proxy-next/client"
)

var liveClient *client.Client

func init() {
	liveClient = client.NewClient(utils.MustGetEnv("LIVE_WG_PROXY_HOST"), time.Second*30)
}
