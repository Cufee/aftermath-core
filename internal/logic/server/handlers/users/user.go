package users

import (
	"errors"
	"strconv"

	"github.com/cufee/aftermath-core/internal/core/server"
	"github.com/cufee/aftermath-core/internal/logic/users"
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

	realm := c.Query("realm")
	if realm == "" {
		return c.Status(400).JSON(server.NewErrorResponse("realm query parameter is required", "c.QueryParam"))
	}

	user, err := users.FindUserByID(userId)
	if err != nil {
		if !errors.Is(err, users.ErrUserNotFound) {
			return c.Status(404).JSON(server.NewErrorResponseFromError(err, "users.FindUserByID"))
		}
		user, err = users.CreateUser(userId)
		if err != nil {
			return c.Status(500).JSON(server.NewErrorResponseFromError(err, "users.CreateUser"))
		}
		// User is created so we can continue
	}

	connection := users.UserConnection{
		UserID:         user.ID,
		ExternalID:     account,
		ConnectionType: users.ConnectionTypeWargaming,
		Metadata:       map[string]interface{}{"verified": false, "realm": realm},
	}

	err = users.UpdateUserConnection(user.ID, connection.ConnectionType, connection, true)
	if err != nil {
		return c.Status(500).JSON(server.NewErrorResponseFromError(err, "users.UpdateUserConnection"))
	}

	return c.JSON(server.NewResponse(connection))
}
