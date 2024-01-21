package users

import (
	"errors"
	"fmt"

	"github.com/cufee/aftermath-core/internal/core/database"
	"github.com/cufee/aftermath-core/internal/core/server"
	"github.com/cufee/aftermath-core/types"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func GetUserHandler(c *fiber.Ctx) error {
	userId := c.Params("id")
	if userId == "" {
		return c.Status(400).JSON(server.NewErrorResponse("id path parameter is required", "c.Param"))
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

	var user types.User
	user.User.ID = details.ID
	user.User.FeatureFlags = details.FeatureFlags
	user.IsBanned = false // TODO: Find ban records

	connections, err := database.GetUserConnections(user.ID)
	if err != nil && errors.Is(err, database.ErrConnectionNotFound) {
		log.Warn().Err(err).Msg("failed to get connections")
		// We can continue without connections
	}
	user.Connections = connections

	subscriptions, err := database.FindActiveSubscriptionsByReferenceIDs(fmt.Sprint(user.ID))
	if err != nil && !errors.Is(err, database.ErrSubscriptionNotFound) {
		log.Warn().Err(err).Msg("failed to get subscriptions")
		// We can continue without subscriptions
	}
	for _, subscription := range subscriptions {
		user.Subscriptions = append(user.Subscriptions, subscription.Type)
	}

	return c.JSON(server.NewResponse(user))
}
