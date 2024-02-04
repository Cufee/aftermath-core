package users

import (
	"errors"
	"strconv"

	"github.com/cufee/aftermath-core/internal/core/database"
	"github.com/cufee/aftermath-core/internal/core/database/models"
	"github.com/cufee/aftermath-core/internal/core/server"
	"github.com/gofiber/fiber/v2"
)

func UpdateWargamingConnectionHandler(c *fiber.Ctx) error {
	userId := c.Params("id")
	if userId == "" {
		return c.Status(400).JSON(server.NewErrorResponse("id path parameter is required", "c.Param"))
	}

	account := c.Params("account")
	_, err := strconv.Atoi(account)
	if err != nil {
		return c.Status(400).JSON(server.NewErrorResponseFromError(err, "strconv.Atoi"))
	}

	user, err := database.GetOrCreateUserByID(userId)
	if err != nil {
		return c.Status(500).JSON(server.NewErrorResponseFromError(err, "users.GetOrCreateUserByID"))
	}

	var update models.ConnectionUpdate
	update.Metadata = map[string]interface{}{"verified": false}
	update.ExternalID = &account

	connection, err := database.UpdateUserConnection(user.ID, models.ConnectionTypeWargaming, update)
	if err != nil {
		if !errors.Is(err, database.ErrConnectionNotFound) {
			return c.Status(500).JSON(server.NewErrorResponseFromError(err, "database.FindUserConnection"))
		}
		connection, err = database.AddUserConnection(user.ID, models.ConnectionTypeWargaming, account, map[string]interface{}{"verified": false})
		if err != nil {
			return c.Status(500).JSON(server.NewErrorResponseFromError(err, "database.AddUserConnection"))
		}
	}

	return c.JSON(server.NewResponse(connection))
}
