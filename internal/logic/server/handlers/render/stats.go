package render

import (
	"image/png"
	"strconv"

	"github.com/cufee/aftermath-core/internal/core/localization"
	"github.com/cufee/aftermath-core/internal/core/server"
	"github.com/cufee/aftermath-core/internal/logic/render"
	"github.com/cufee/aftermath-core/internal/logic/stats"
	"github.com/labstack/echo/v4"
)

func SessionHandler(c echo.Context) error {
	account := c.Param("account")
	accountId, err := strconv.Atoi(account)
	if err != nil {
		return c.JSON(400, server.NewErrorResponseFromError(err, "strconv.Atoi"))
	}
	realm := c.QueryParam("realm")
	if realm == "" {
		return c.JSON(400, server.NewErrorResponse("realm query parameter is required", "c.QueryParam"))
	}

	session, err := stats.GetCurrentPlayerSession(realm, accountId)
	if err != nil {
		return c.JSON(500, server.NewErrorResponseFromError(err, "stats.GetCurrentPlayerSession"))
	}

	img, err := render.RenderStatsImage(session, nil, localization.LanguageEN)
	if err != nil {
		return c.JSON(500, server.NewErrorResponseFromError(err, "render.RenderStatsImage"))
	}

	c.Response().Header().Set("Content-Type", "image/png")
	return png.Encode(c.Response(), img)
}
