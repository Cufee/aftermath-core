package moderation

import (
	"github.com/cufee/aftermath-core/internal/core/cloudinary"
	"github.com/cufee/aftermath-core/internal/core/server"
	"github.com/cufee/aftermath-core/internal/logic/content"
	"github.com/cufee/aftermath-core/types"
	"github.com/gofiber/fiber/v2"
)

func UploadBackgroundImageHandler(c *fiber.Ctx) error {
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

	image, err := content.EncodeRemoteImage(body.Data)
	if err != nil {
		return c.Status(400).JSON(server.NewErrorResponseFromError(err, "content.EncodeRemoteImage"))
	}

	link, err := cloudinary.DefaultClient.ManualUpload(image)
	if err != nil {
		return c.Status(500).JSON(server.NewErrorResponseFromError(err, "content.UploadUserImage"))
	}

	return c.JSON(server.NewResponse(link))
}
