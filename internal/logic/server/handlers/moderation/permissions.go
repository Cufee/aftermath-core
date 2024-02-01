package moderation

import (
	"github.com/cufee/aftermath-core/internal/core/server"
	"github.com/cufee/aftermath-core/permissions/v1"
	"github.com/gofiber/fiber/v2"
)

func GetPermissionsMapHandler(c *fiber.Ctx) error {
	return c.JSON(server.NewResponse(permissions.PermissionsMap))
}
