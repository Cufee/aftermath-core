package users

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/cufee/aftermath-core/internal/core/database"
	"github.com/cufee/aftermath-core/internal/core/database/models"
	"github.com/cufee/aftermath-core/internal/core/localization"
	"github.com/cufee/aftermath-core/internal/core/server"
	"github.com/cufee/aftermath-core/internal/core/utils"
	"github.com/cufee/aftermath-core/types"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

var frontendURL = utils.MustGetEnv("FRONTEND_URL")
var authWargamingAppID = utils.MustGetEnv("AUTH_WARGAMING_APP_ID")

func CompleteUserVerificationHandler(c *fiber.Ctx) error {
	nonceID := c.Params("nonce")
	if nonceID == "" {
		return c.Status(400).JSON(server.NewErrorResponse("nonce path parameter is required", "c.Param"))
	}

	var payload types.UserVerificationPayload
	err := c.BodyParser(&payload)
	if err != nil {
		return c.Status(400).JSON(server.NewErrorResponseFromError(err, "c.BodyParser"))
	}
	if payload.AccountID == "" {
		return c.Status(400).JSON(server.NewErrorResponse("payload accountID is required", "c.BodyParser"))
	}

	nonce, err := database.GetNonceByID(nonceID)
	if err != nil {
		if !errors.Is(err, database.ErrNonceNotFound) {
			return c.Status(404).JSON(server.NewErrorResponseFromError(err, "database.GetNonceByID"))
		}
		return c.Status(400).JSON(server.NewErrorResponseFromError(err, "database.GetNonceByID"))
	}
	go func() {
		err := database.ExpireNonceByID(nonceID)
		if err != nil {
			log.Warn().Err(err).Msg("failed to expire nonce")
		}
	}()

	if nonce.ReferenceID == "" {
		return c.Status(400).JSON(server.NewErrorResponse("nonce referenceID is required", "c.Param"))
	}

	user, err := database.GetUserByID(nonce.ReferenceID)
	if err != nil {
		return c.Status(404).JSON(server.NewErrorResponseFromError(err, "users.FindUserByID"))
	}

	var update models.ConnectionUpdate
	update.Metadata = map[string]interface{}{"verified": true}
	update.ExternalID = &payload.AccountID

	connection, err := database.UpdateUserConnection(user.ID, models.ConnectionTypeWargaming, update)
	if err != nil {
		if !errors.Is(err, database.ErrConnectionNotFound) {
			return c.Status(500).JSON(server.NewErrorResponseFromError(err, "database.FindUserConnection"))
		}
		connection, err = database.AddUserConnection(user.ID, models.ConnectionTypeWargaming, payload.AccountID, map[string]interface{}{"verified": true})
		if err != nil {
			return c.Status(500).JSON(server.NewErrorResponseFromError(err, "database.AddUserConnection"))
		}
	}

	// Update user content reference ID
	_, err = database.UpdateUserContentReferenceID[string](user.ID, models.UserContentTypePersonalBackground, payload.AccountID)
	if err != nil && !errors.Is(database.ErrUserContentNotFound, err) {
		return c.Status(500).JSON(server.NewErrorResponseFromError(err, "database.UpdateUserContent"))
	}

	go func(externalId string) {
		// Mark all other connections for this account as unverified
		err := database.UpdateManyConnectionsByReferenceID(externalId, models.ConnectionTypeWargaming, models.ConnectionUpdate{Metadata: map[string]interface{}{"verified": false}})
		if err != nil {
			log.Err(err).Msg("failed to find connections by reference ID")
		}
	}(connection.ExternalID)

	return c.JSON(server.NewResponse(connection))
}

func StartUserVerificationHandler(c *fiber.Ctx) error {
	userId := c.Params("id")
	if userId == "" {
		return c.Status(400).JSON(server.NewErrorResponse("id path parameter is required", "c.Param"))
	}
	realm := c.Query("realm")
	if realm == "" {
		return c.Status(400).JSON(server.NewErrorResponse("realm query parameter is required", "c.Param"))
	}

	user, err := database.GetOrCreateUserByID(userId)
	if err != nil {
		return c.Status(500).JSON(server.NewErrorResponseFromError(err, "users.CreateUser"))
	}

	nonce, err := database.NewNonce(user.ID, time.Duration(5*time.Minute))
	if err != nil {
		return c.Status(500).JSON(server.NewErrorResponseFromError(err, "database.NewNonce"))
	}

	link, err := loginUrlFromRealm(realm, localization.LanguageEN.WargamingCode, nonce)
	if err != nil {
		return c.Status(500).JSON(server.NewErrorResponseFromError(err, "loginUrlFromRealm"))
	}

	return c.JSON(server.NewResponse(link))
}

func loginUrlFromRealm(realm, language, nonce string) (string, error) {
	var base string
	switch strings.ToUpper(realm) {
	case "EU":
		base = "https://api.worldoftanks.eu"
	case "NA":
		base = "https://api.worldoftanks.com"
	case "AS":
		base = "https://api.worldoftanks.asia"
	default:
		return "", errors.New("unknown realm")
	}

	return fmt.Sprintf("%s/wot/auth/login/?redirect_uri=%s&language=%s&application_id=%s", base, url.QueryEscape(fmt.Sprintf("%s/auth/wargaming/redirect/%s", frontendURL, nonce)), language, authWargamingAppID), nil
}
