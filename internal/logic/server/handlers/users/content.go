package users

import (
	"errors"
	"strconv"

	"github.com/cufee/aftermath-core/internal/core/database"
	"github.com/cufee/aftermath-core/internal/core/database/models"
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

	connection, err := database.FindUserConnection(details.ID, models.ConnectionTypeWargaming)
	if err != nil {
		return c.Status(404).JSON(server.NewErrorResponseFromError(err, "models.FindUserConnection"))
	}
	if connection.Metadata["verified"] != true {
		return c.Status(400).JSON(server.NewErrorResponse("user account is not verified", ""))
	}

	link, err := content.UploadUserImage(details.ID, body.Data)
	if err != nil {
		return c.Status(500).JSON(server.NewErrorResponseFromError(err, "content.UploadUserImage"))
	}

	err = database.UpdateUserContent(details.ID, connection.ExternalID, body.Type, link, nil, true)
	if err != nil {
		return c.Status(500).JSON(server.NewErrorResponseFromError(err, "database.UpdateUserContent"))
	}

	return c.JSON(server.NewResponse(link))
}

func SelectBackgroundPresetHandler(c *fiber.Ctx) error {
	userId := c.Params("id")
	if userId == "" {
		return c.Status(400).JSON(server.NewErrorResponse("id path parameter is required", "c.Param"))
	}
	index := c.Params("index")
	if index == "" {
		return c.Status(400).JSON(server.NewErrorResponse("index path parameter is required", "c.Param"))
	}
	i, err := strconv.Atoi(index)
	if err != nil {
		return c.Status(400).JSON(server.NewErrorResponseFromError(err, "strconv.Atoi"))
	}

	data, err := database.GetAppConfiguration[[]string]("backgroundImagesSelection")
	if err != nil {
		return c.Status(500).JSON(server.NewErrorResponseFromError(err, "database.GetAppConfiguration"))
	}

	if i < 0 || i >= len(data.Value) {
		return c.Status(400).JSON(server.NewErrorResponse("index out of range", ""))
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

	connection, err := database.FindUserConnection(details.ID, models.ConnectionTypeWargaming)
	if err != nil {
		return c.Status(404).JSON(server.NewErrorResponseFromError(err, "models.FindUserConnection"))
	}

	err = database.UpdateUserContent(details.ID, connection.ExternalID, models.UserContentTypePersonalBackground, data.Value[i], nil, true)
	if err != nil {
		return c.Status(500).JSON(server.NewErrorResponseFromError(err, "database.UpdateUserContent"))
	}

	return c.JSON(server.NewResponse(data.Value[i]))
}
