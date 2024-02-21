package server

import (
	"os"

	"github.com/cufee/aftermath-core/internal/logic/server/handlers/accounts"
	"github.com/cufee/aftermath-core/internal/logic/server/handlers/content"
	"github.com/cufee/aftermath-core/internal/logic/server/handlers/moderation"
	"github.com/cufee/aftermath-core/internal/logic/server/handlers/render"
	"github.com/cufee/aftermath-core/internal/logic/server/handlers/stats"
	"github.com/cufee/aftermath-core/internal/logic/server/handlers/users"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/gofiber/contrib/fiberzerolog"
	_ "github.com/joho/godotenv/autoload"
)

func Start() {
	app := fiber.New(fiber.Config{
		Network: os.Getenv("NETWORK"),
	})
	app.Use(fiberzerolog.New())
	app.Use(recover.New())

	app.Get("/healthy", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	v1 := app.Group("/v1")

	renderV1 := v1.Group("/render")
	renderV1.Post("/session/user/:id", render.SessionFromUserHandler)
	renderV1.Post("/session/account/:account", render.SessionFromIDHandler)

	statsV1 := v1.Group("/stats")
	statsV1.Get("/session/user/:id", stats.SessionFromUserHandler)
	statsV1.Get("/session/account/:account", stats.SessionFromIDHandler)
	statsV1.Post("/session/account/:account/reset", stats.RecordPlayerSession)

	accountsV1 := v1.Group("/accounts")
	accountsV1.Get("/search", accounts.SearchAccountsHandler)

	usersV1 := v1.Group("/users")
	usersV1.Get("/:id", users.GetUserHandler)
	usersV1.Post("/:id/content", users.UploadUserContentHandler)
	usersV1.Get("/:id/content/select", content.PreviewCurrentBackgroundSelectionHandler)
	usersV1.Post("/:id/content/select/:index", users.SelectBackgroundPresetHandler)
	usersV1.Post("/:id/connections/wargaming/:account", users.UpdateWargamingConnectionHandler)

	connectionsV1 := v1.Group("/connections")
	connectionsV1.Get("/wargaming/verify/:id", users.StartUserVerificationHandler)
	connectionsV1.Post("/wargaming/verify/:nonce", users.CompleteUserVerificationHandler)

	moderationV1 := v1.Group("/moderation")
	moderationV1.Get("/permissions", moderation.GetPermissionsMapHandler)

	moderationV1.Get("/content/rotate", moderation.RotateBackgroundImagesHandler)
	moderationV1.Post("/content/upload", moderation.UploadBackgroundImageHandler)

	moderationV1.Get("/subscriptions/user/:userId", moderation.GetUserSubscriptionsHandler)
	moderationV1.Post("/subscriptions", moderation.CreateUserSubscriptionsHandler)
	moderationV1.Get("/subscriptions/:id", moderation.GetSubscriptionHandler)
	moderationV1.Patch("/subscriptions/:id", moderation.UpdateSubscriptionHandler)

	panic(app.Listen(":" + os.Getenv("PORT")))
}
