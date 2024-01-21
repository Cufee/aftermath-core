package users

import (
	"errors"

	"github.com/cufee/aftermath-core/internal/core/database"
	"github.com/cufee/aftermath-core/internal/core/server"
	"github.com/cufee/aftermath-core/internal/logic/content"
	"github.com/cufee/aftermath-core/types"
	"github.com/gofiber/fiber/v2"
)

func UploadUserContentHandler(c *fiber.Ctx) error {
	userId := c.Params("id")
	if userId == "" {
		return c.Status(400).JSON(server.NewErrorResponse("id path parameter is required", "c.Param"))
	}

	var body types.UserContentPayload[string]
	err := c.BodyParser(&body)
	if err != nil {
		return c.Status(400).JSON(server.NewErrorResponseFromError(err, "c.BodyParser"))
	}
	if !body.Type.Valid() {
		return c.Status(400).JSON(server.NewErrorResponse("body type invalid", "c.BodyParser"))
	}
	if body.Data == "" {
		return c.Status(400).JSON(server.NewErrorResponse("body data invalid", ""))
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

	// TODO: ReferenceID should be a verified account ID
	// connections, err := database.FindUserConnection(details.ID, models.ConnectionTypeWargaming)
	// if err != nil {
	// 	if !errors.Is(err, database.ErrConnectionNotFound) {
	// 		return c.Status(404).JSON(server.NewErrorResponseFromError(err, "models.FindUserConnection"))
	// 	}
	// }

	link, err := content.UploadUserImage(details.ID, body.Data)
	if err != nil {
		return c.Status(500).JSON(server.NewErrorResponseFromError(err, "content.UploadUserImage"))
	}

	err = database.UpdateUserContent(details.ID, details.ID, body.Type, link, nil, true)
	if err != nil {
		return c.Status(500).JSON(server.NewErrorResponseFromError(err, "database.UpdateUserContent"))
	}

	return c.JSON(server.NewResponse(link))
}
