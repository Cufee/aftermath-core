package types

import (
	"github.com/cufee/aftermath-core/internal/core/server"
	"github.com/cufee/aftermath-core/internal/logic/server/handlers/stats"
	"github.com/cufee/aftermath-core/internal/logic/users"
	wg "github.com/cufee/am-wg-proxy-next/types"
)

type AccountsSearchResponse server.Response[wg.Account]

type RenderSessionResponse server.Response[string]

type StatsSessionResponse server.Response[stats.SessionStatsResponse]

type UsersConnectionResponse server.Response[users.UserConnection]
