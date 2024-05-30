package types

import (
	"strconv"
	"time"

	"github.com/cufee/aftermath-core/internal/core/database/models"
	"github.com/cufee/aftermath-core/internal/core/server"
	"github.com/cufee/aftermath-core/internal/core/wargaming"
)

type UsersResponse server.Response[User]
type UserConnection models.UserConnection
type UsersConnectionResponse server.Response[UserConnection]

type User struct {
	models.CompleteUser `json:",inline"`
	Restrictions        []UserRestriction `json:"restrictions"`
}

type UserRestriction struct {
	CreatedAt   time.Time `json:"createdAt"`
	ExpiresAt   time.Time `json:"expiresAt"`
	UserMessage string    `json:"userMessage"`
	Comment     string    `json:"comment"`
	Scopes      []string  `json:"scopes"`
}

func (u User) ActiveRestriction(scopes ...string) *UserRestriction {
	now := time.Now()
	for _, restriction := range u.Restrictions {
		if restriction.ExpiresAt.After(now) {
			if len(restriction.Scopes) == 0 {
				return &restriction
			}

			for _, providedScope := range scopes {
				for _, restrictionScope := range restriction.Scopes {
					if restrictionScope == providedScope {
						return &restriction
					}
				}
			}

		}
	}
	return nil
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
	wargamingConnection.Realm = wargaming.Clients.Live.RealmFromAccountID(strconv.Itoa(id))
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
