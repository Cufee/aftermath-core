package content

import (
	"github.com/cufee/aftermath-core/internal/core/server"
	"github.com/cufee/aftermath-core/internal/logic/preview"

	"github.com/gofiber/fiber/v2"
)

func PreviewCurrentBackgroundSelectionHandler(c *fiber.Ctx) error {
	preview, err := preview.CurrentBackgroundsPreview()
	if err != nil {
		return c.Status(500).JSON(server.NewErrorResponseFromError(err, "preview.CurrentBackgroundsPreview"))
	}

	return c.JSON(server.NewResponse(preview))
}
