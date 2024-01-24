package users

import (
	"strconv"

	"github.com/cufee/aftermath-core/internal/core/database"
	"github.com/cufee/aftermath-core/internal/core/database/models"
	"github.com/cufee/aftermath-core/internal/core/server"
	"github.com/cufee/aftermath-core/internal/logic/content"
	"github.com/cufee/aftermath-core/permissions/v1"
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
	if body.Data == "" {
		return c.Status(400).JSON(server.NewErrorResponse("body data invalid", ""))
	}

	user, err := database.GetOrCreateUserByID(userId)
	if err != nil {
		return c.Status(500).JSON(server.NewErrorResponseFromError(err, "database.GetOrCreateUserByID"))
	}
	if !user.Permissions().Has(permissions.UploadPersonalBackground) {
		return c.Status(403).JSON(server.NewErrorResponse("user has no permissions", ""))
	}

	connection := user.Connection(models.ConnectionTypeWargaming)
	if connection == nil {
		return c.Status(404).JSON(server.NewErrorResponse("user has no wargaming connection", ""))
	}

	link, err := content.UploadUserImage(user.ID, body.Data)
	if err != nil {
		return c.Status(500).JSON(server.NewErrorResponseFromError(err, "content.UploadUserImage"))
	}

	err = database.UpdateUserContent(user.ID, connection.ExternalID, models.UserContentTypePersonalBackground, link, nil, true)
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

	user, err := database.GetOrCreateUserByID(userId)
	if err != nil {
		return c.Status(500).JSON(server.NewErrorResponseFromError(err, "users.CreateUser"))
	}
	if !user.Permissions().Has(permissions.SelectPersonalBackgroundPreset) {
		return c.Status(403).JSON(server.NewErrorResponse("user has no permissions", ""))
	}

	connection := user.Connection(models.ConnectionTypeWargaming)
	if connection == nil {
		return c.Status(404).JSON(server.NewErrorResponse("user has no wargaming connection", ""))
	}
	if verified, ok := connection.Metadata["verified"].(bool); !ok || !verified {
		return c.Status(400).JSON(server.NewErrorResponse("user wargaming connection not verified", ""))
	}

	err = database.UpdateUserContent(user.ID, connection.ExternalID, models.UserContentTypePersonalBackground, data.Value[i], nil, true)
	if err != nil {
		return c.Status(500).JSON(server.NewErrorResponseFromError(err, "database.UpdateUserContent"))
	}

	return c.JSON(server.NewResponse(data.Value[i]))
}
