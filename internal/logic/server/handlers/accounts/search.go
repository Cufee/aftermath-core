package accounts

import (
	"github.com/cufee/aftermath-core/internal/core/server"
	"github.com/cufee/aftermath-core/internal/core/wargaming"
	"github.com/gofiber/fiber/v2"
)

func SearchAccountsHandler(c *fiber.Ctx) error {
	realm := c.Query("realm")
	searchQuery := c.Query("search")
	if realm == "" || searchQuery == "" {
		return c.Status(400).JSON(server.NewErrorResponse("realm and search query parameters are required", "c.QueryParam"))
	}

	accounts, err := wargaming.Clients.Live.SearchAccounts(realm, searchQuery)
	if err != nil {
		return c.Status(500).JSON(server.NewErrorResponseFromError(err, "wargaming.Clients.Live.SearchAccounts"))
	}

	return c.JSON(server.NewResponse(accounts))
}
