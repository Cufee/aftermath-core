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

	user, err := database.FindUserByID(userId)
	if err != nil {
		if !errors.Is(err, database.ErrUserNotFound) {
			return c.Status(404).JSON(server.NewErrorResponseFromError(err, "users.FindUserByID"))
		}
		user, err = database.CreateUser(userId)
		if err != nil {
			return c.Status(500).JSON(server.NewErrorResponseFromError(err, "users.CreateUser"))
		}
		// User is created so we can continue
	}

	connection := models.UserConnection{
		UserID:         user.ID,
		ExternalID:     account,
		ConnectionType: models.ConnectionTypeWargaming,
		Metadata:       map[string]interface{}{"verified": false},
	}

	err = database.UpdateUserConnection(user.ID, connection.ConnectionType, connection, true)
	if err != nil {
		return c.Status(500).JSON(server.NewErrorResponseFromError(err, "users.UpdateUserConnection"))
	}

	return c.JSON(server.NewResponse(connection))
}
