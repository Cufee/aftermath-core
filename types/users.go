package types

import (
	"strconv"

	"github.com/cufee/aftermath-core/internal/core/database/models"
	"github.com/cufee/aftermath-core/internal/core/server"
	"github.com/rs/zerolog/log"
)

type UsersResponse server.Response[User]
type UserConnection models.UserConnection
type UsersConnectionResponse server.Response[UserConnection]

type User struct {
	models.User   `json:",inline"`
	IsBanned      bool                      `json:"is_banned"`
	Connections   []models.UserConnection   `json:"connections"`
	Subscriptions []models.SubscriptionType `json:"subscriptions"`
}

func (u User) WargamingConnection() (int, bool) {
	for _, connection := range u.Connections {
		if connection.ConnectionType == models.ConnectionTypeWargaming {
			accountID, err := strconv.Atoi(connection.ExternalID)
			if err != nil {
				log.Warn().Err(err).Msg("failed to parse account id")
				return 0, false
			}
			return accountID, true
		}
	}
	return 0, false
}

type UserContentPayload[T any] struct {
	Type models.UserContentType `json:"type"`
	Data T                      `json:"data"`
}
