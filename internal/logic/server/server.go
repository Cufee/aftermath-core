package server

import (
	"os"

	"github.com/cufee/aftermath-core/internal/logic/server/handlers/accounts"
	"github.com/cufee/aftermath-core/internal/logic/server/handlers/render"
	"github.com/cufee/aftermath-core/internal/logic/server/handlers/users"

	"github.com/gofiber/fiber/v2"

	"github.com/gofiber/contrib/fiberzerolog"
	_ "github.com/joho/godotenv/autoload"
)

func Start() {
	app := fiber.New(fiber.Config{
		Network: os.Getenv("NETWORK"),
	})

	app.Use(fiberzerolog.New())

	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	v1 := app.Group("/v1")

	renderV1 := v1.Group("/render")
	renderV1.Get("/session/user/:id", render.SessionFromUserHandler)
	renderV1.Get("/session/account/:account", render.SessionFromIDHandler)

	accountsV1 := v1.Group("/accounts")
	accountsV1.Get("/search", accounts.SearchAccountsHandler)

	usersV1 := v1.Group("/users")
	usersV1.Post("/:id/connections/wargaming/:account", users.UpdateWargamingConnectionHandler)

	panic(app.Listen(":" + os.Getenv("PORT")))
}
