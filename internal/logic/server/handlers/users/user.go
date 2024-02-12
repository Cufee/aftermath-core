package users

import (
	"github.com/cufee/aftermath-core/internal/core/database"
	"github.com/cufee/aftermath-core/internal/core/server"
	"github.com/cufee/aftermath-core/permissions/v2"
	"github.com/cufee/aftermath-core/types"

	"github.com/gofiber/fiber/v2"
)

func GetUserHandler(c *fiber.Ctx) error {
	userId := c.Params("id")
	if userId == "" {
		return c.Status(400).JSON(server.NewErrorResponse("id path parameter is required", "c.Param"))
	}

	user, err := database.GetOrCreateUserByID(userId)
	if err != nil {
		return c.Status(500).JSON(server.NewErrorResponseFromError(err, "users.CreateUser"))
	}
	if user.User.Permissions == permissions.Blank {
		user.User.Permissions = permissions.User
	}

	var extended types.User
	extended.CompleteUser = *user
	extended.IsBanned = false // TODO: Find ban records

	return c.JSON(server.NewResponse(user))
}
