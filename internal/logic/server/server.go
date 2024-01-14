package server

import (
	"os"

	"github.com/cufee/aftermath-core/internal/logic/server/handlers/render"
	"github.com/labstack/echo/v4"

	_ "github.com/joho/godotenv/autoload"
)

func NewServer() func() {
	e := echo.New()

	v1 := e.Group("/v1")

	renderV1 := v1.Group("/render")
	renderV1.GET("/session/:account", render.SessionHandler)

	return func() {
		e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
	}
}
