package types

import (
	"strconv"

	"github.com/cufee/aftermath-core/internal/core/database/models"
	"github.com/cufee/aftermath-core/internal/core/server"
	"github.com/cufee/aftermath-core/utils"
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

type wargamingConnection struct {
	AccountID int    `json:"account_id"`
	Verified  bool   `json:"verified"`
	Realm     string `json:"realm"`
}

func (u User) WargamingConnection() (*wargamingConnection, bool) {
	for _, connection := range u.Connections {
		if connection.ConnectionType == models.ConnectionTypeWargaming {
			accountID, err := strconv.Atoi(connection.ExternalID)
			if err != nil {
				log.Warn().Err(err).Msg("failed to parse account id")
				return nil, false
			}
			verified, _ := connection.Metadata["verified"].(bool)
			return &wargamingConnection{
				AccountID: accountID,
				Verified:  verified,
				Realm:     utils.RealmFromAccountID(accountID),
			}, true
		}
	}
	return nil, false
}

type UserContentPayload[T any] struct {
	Type models.UserContentType `json:"type"`
	Data T                      `json:"data"`
}

type UserVerificationPayload struct {
	AccessTokenExpiresAt int64  `json:"access_token_expires_at"` // Not used
	AccessToken          string `json:"access_token"`            // Not used

	AccountID string `json:"account_id"`
}
