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
	verified := c.Query("verified") == "true"

	account := c.Params("account")
	_, err := strconv.Atoi(account)
	if err != nil {
		return c.Status(400).JSON(server.NewErrorResponseFromError(err, "strconv.Atoi"))
	}

	user, err := database.GetOrCreateUserByID(userId)
	if err != nil {
		return c.Status(500).JSON(server.NewErrorResponseFromError(err, "users.GetOrCreateUserByID"))
	}

	connection := models.UserConnection{
		UserID:         user.ID,
		ExternalID:     account,
		ConnectionType: models.ConnectionTypeWargaming,
		Metadata:       map[string]interface{}{"verified": verified},
	}

	err = database.UpdateUserConnection(user.ID, connection.ConnectionType, connection, true)
	if err != nil {
		return c.Status(500).JSON(server.NewErrorResponseFromError(err, "users.UpdateUserConnection"))
	}

	if verified {
		_, err = database.UpdateUserContentReferenceID[string](user.ID, models.UserContentTypePersonalBackground, account)
		if err != nil && !errors.Is(database.ErrUserContentNotFound, err) {
			return c.Status(500).JSON(server.NewErrorResponseFromError(err, "database.UpdateUserContent"))
		}
	}

	return c.JSON(server.NewResponse(connection))
}
