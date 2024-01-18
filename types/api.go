package types

import (
	"github.com/cufee/aftermath-core/internal/core/server"
	"github.com/cufee/aftermath-core/internal/logic/server/handlers/stats"
	"github.com/cufee/aftermath-core/internal/logic/users"
	wg "github.com/cufee/am-wg-proxy-next/types"
)

type Account wg.Account
type AccountsSearchResponse server.Response[Account]

type EncodedImage string
type RenderSessionResponse server.Response[EncodedImage]

type SessionStats stats.SessionStatsResponse
type StatsSessionResponse server.Response[SessionStats]

type UserConnection users.UserConnection
type UsersConnectionResponse server.Response[UserConnection]
