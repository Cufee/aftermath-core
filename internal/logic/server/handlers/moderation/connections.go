package moderation

import (
	"errors"

	"github.com/cufee/aftermath-core/internal/core/database"
	"github.com/cufee/aftermath-core/internal/core/database/models"
	"github.com/cufee/aftermath-core/internal/core/server"
	"github.com/cufee/aftermath-core/types"
	"github.com/gofiber/fiber/v2"
)

func ForceUpdateConnectionHandler(c *fiber.Ctx) error {
	var opts types.ForceUpdateConnectionPayload
	err := c.BodyParser(&opts)
	if err != nil {
		return c.Status(400).JSON(server.NewErrorResponseFromError(err, "c.BodyParser"))
	}
	if opts.UserID == "" || opts.ConnectionID == "" || opts.ConnectionType == "" {
		return c.Status(400).JSON(server.NewErrorResponse("userId, connectionId, connectionType are all required", "payload"))
	}
	connectionType, valid := opts.Type()
	if !valid {
		return c.Status(400).JSON(server.NewErrorResponse("invalid connectionType provided", "payload"))
	}

	user, err := database.GetOrCreateUserByID(opts.UserID)
	if err != nil {
		return c.Status(500).JSON(server.NewErrorResponseFromError(err, "users.GetOrCreateUserByID"))
	}

	var update models.ConnectionUpdate
	update.Metadata = opts.Metadata
	update.ExternalID = &opts.ConnectionID

	connection, err := database.UpdateUserConnection(user.ID, connectionType, update)
	if err != nil {
		if !errors.Is(err, database.ErrConnectionNotFound) {
			return c.Status(500).JSON(server.NewErrorResponseFromError(err, "database.FindUserConnection"))
		}
		connection, err = database.AddUserConnection(user.ID, models.ConnectionTypeWargaming, opts.ConnectionID, opts.Metadata)
		if err != nil {
			return c.Status(500).JSON(server.NewErrorResponseFromError(err, "database.AddUserConnection"))
		}
	}

	return c.JSON(server.NewResponse(connection))
}
