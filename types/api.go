package types

import (
	"github.com/cufee/aftermath-core/dataprep"
	"github.com/cufee/aftermath-core/internal/core/database/models"
	"github.com/cufee/aftermath-core/internal/core/server"
	wg "github.com/cufee/am-wg-proxy-next/types"
)

type Account wg.Account
type AccountsSearchResponse server.Response[Account]

type EncodedImage string
type RenderSessionResponse server.Response[EncodedImage]

type StatsSessionResponse server.Response[dataprep.SessionStats]

type UserConnection models.UserConnection
type UsersConnectionResponse server.Response[UserConnection]
