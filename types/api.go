package types

import (
	"github.com/cufee/aftermath-core/internal/core/database/models"
	"github.com/cufee/aftermath-core/internal/core/server"
	"github.com/cufee/aftermath-core/internal/logic/dataprep"
	wg "github.com/cufee/am-wg-proxy-next/types"
)

type SessionStatsResponse struct {
	Realm      string                `json:"realm"`
	Locale     string                `json:"locale"`
	LastBattle int                   `json:"last_battle"`
	Clan       wg.Clan               `json:"clan"`
	Account    wg.Account            `json:"account"`
	Cards      dataprep.SessionCards `json:"cards"`
}

type Account wg.Account
type AccountsSearchResponse server.Response[Account]

type EncodedImage string
type RenderSessionResponse server.Response[EncodedImage]

type StatsSessionResponse server.Response[SessionStatsResponse]

type UserConnection models.UserConnection
type UsersConnectionResponse server.Response[UserConnection]
