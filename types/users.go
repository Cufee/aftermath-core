package types

import (
	"strconv"
	"time"

	"github.com/cufee/aftermath-core/internal/core/database/models"
	"github.com/cufee/aftermath-core/internal/core/server"
	"github.com/cufee/aftermath-core/utils"
)

type UsersResponse server.Response[User]
type UserConnection models.UserConnection
type UsersConnectionResponse server.Response[UserConnection]

type User struct {
	models.CompleteUser `json:",inline"`
	IsBanned            bool `json:"is_banned"`
}

func (u User) WargamingConnection() (*wargamingConnection, bool) {
	connection := u.Connection(models.ConnectionTypeWargaming)
	if connection == nil {
		return nil, false
	}
	id, err := strconv.Atoi(connection.ExternalID)
	if err != nil {
		return nil, false
	}

	var wargamingConnection wargamingConnection
	wargamingConnection.AccountID = id
	wargamingConnection.Verified, _ = connection.Metadata["verified"].(bool)
	wargamingConnection.Realm = utils.RealmFromAccountID(id)
	return &wargamingConnection, true
}

type wargamingConnection struct {
	AccountID int    `json:"account_id"`
	Verified  bool   `json:"verified"`
	Realm     string `json:"realm"`
}

type UserContentPayload[T any] struct {
	Data T `json:"data"`
}

type UserVerificationPayload struct {
	AccessTokenExpiresAt int64  `json:"access_token_expires_at"` // Not used
	AccessToken          string `json:"access_token"`            // Not used

	AccountID string `json:"account_id"`
}

type UserSubscriptionPayload struct {
	UserID      string        `json:"user_id"`
	ReferenceID string        `json:"reference_id"`
	Duration    time.Duration `json:"duration"`
	Type        string        `json:"type"`
}
