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

	user, err := database.FindUserByID(nonce.ReferenceID)
	if err != nil {
		return c.Status(404).JSON(server.NewErrorResponseFromError(err, "users.FindUserByID"))
	}

	connection := models.UserConnection{
		UserID:         user.ID,
		ExternalID:     payload.AccountID,
		ConnectionType: models.ConnectionTypeWargaming,
		Metadata:       map[string]interface{}{"verified": true},
	}

	err = database.UpdateUserConnection(user.ID, connection.ConnectionType, connection, true)
	if err != nil {
		return c.Status(500).JSON(server.NewErrorResponseFromError(err, "users.UpdateUserConnection"))
	}

	_, err = database.UpdateUserContentReferenceID[string](user.ID, models.UserContentTypePersonalBackground, payload.AccountID)
	if err != nil && !errors.Is(database.ErrUserContentNotFound, err) {
		return c.Status(500).JSON(server.NewErrorResponseFromError(err, "database.UpdateUserContent"))
	}

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

	details, err := database.FindUserByID(userId)
	if err != nil {
		if !errors.Is(err, database.ErrUserNotFound) {
			return c.Status(404).JSON(server.NewErrorResponseFromError(err, "users.FindUserByID"))
		}
		details, err = database.CreateUser(userId)
		if err != nil {
			return c.Status(500).JSON(server.NewErrorResponseFromError(err, "users.CreateUser"))
		}
		// User is created so we can continue
	}

	nonce, err := database.NewNonce(details.ID, time.Duration(5*time.Minute))
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