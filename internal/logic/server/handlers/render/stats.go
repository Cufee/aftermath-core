package render

import (
	"bytes"
	"encoding/base64"
	"errors"
	"image/png"
	"strconv"

	"github.com/cufee/aftermath-core/internal/core/localization"
	"github.com/cufee/aftermath-core/internal/core/server"
	"github.com/cufee/aftermath-core/internal/logic/render/assets"
	render "github.com/cufee/aftermath-core/internal/logic/render/session"
	"github.com/cufee/aftermath-core/internal/logic/stats"
	"github.com/cufee/aftermath-core/internal/logic/users"
	"github.com/gofiber/fiber/v2"
)

func SessionFromIDHandler(c *fiber.Ctx) error {
	account := c.Params("account")
	accountId, err := strconv.Atoi(account)
	if err != nil {
		return c.Status(400).JSON(server.NewErrorResponseFromError(err, "strconv.Atoi"))
	}
	realm := c.Query("realm")
	if realm == "" {
		return c.Status(400).JSON(server.NewErrorResponse("realm query parameter is required", "c.QueryParam"))
	}

	imageData, err := getEncodedSessionImage(realm, accountId)
	if err != nil {
		return c.Status(500).JSON(server.NewErrorResponseFromError(err, "getEncodedSessionImage"))
	}

	return c.JSON(server.NewResponse(imageData))
}

func SessionFromUserHandler(c *fiber.Ctx) error {
	user := c.Params("id")
	if user == "" {
		return c.Status(400).JSON(server.NewErrorResponse("id path parameter is required", "c.Param"))
	}

	connection, err := users.FindUserConnection(user, users.ConnectionTypeWargaming)
	if err != nil {
		if errors.Is(err, users.ErrConnectionNotFound) {
			return c.Status(404).JSON(server.NewErrorResponseFromError(err, "users.FindUserConnection"))
		}
		return c.Status(500).JSON(server.NewErrorResponseFromError(err, "users.FindUserConnection"))
	}

	realm, ok := connection.Metadata["realm"].(string)
	if realm == "" || !ok {
		return c.Status(500).JSON(server.NewErrorResponse("invalid connection", "connection.Metadata"))
	}

	accountId, err := strconv.Atoi(connection.ExternalID)
	if err != nil {
		return c.Status(500).JSON(server.NewErrorResponse("invalid connection", "strconv.Atoi"))
	}

	imageData, err := getEncodedSessionImage(realm, accountId)
	if err != nil {
		return c.Status(500).JSON(server.NewErrorResponseFromError(err, "getEncodedSessionImage"))
	}

	return c.JSON(server.NewResponse(imageData))
}

func getEncodedSessionImage(realm string, accountId int) (string, error) {
	session, err := stats.GetCurrentPlayerSession(realm, accountId)
	if err != nil {
		return "", err
	}

	averages, err := stats.GetVehicleAverages(session.Diff.Vehicles)
	if err != nil {
		return "", err
	}

	opts := stats.SortOptions{
		By:    stats.SortByLastBattle,
		Limit: 5,
	}
	player := render.PlayerData{
		Snapshot: session,
		Averages: averages,
		Vehicles: stats.SortVehicles(session.Diff.Vehicles, averages, opts),
	}

	bgImage, _ := assets.GetImage("images/backgrounds/default")
	options := render.RenderOptions{
		PromoText:       []string{"Aftermath is back!", "amth.one/join"},
		Locale:          localization.LanguageEN,
		CardStyle:       render.DefaultCardStyle(nil),
		BackgroundImage: bgImage,
	}

	img, err := render.RenderStatsImage(player, options)
	if err != nil {
		return "", err
	}

	encoded := new(bytes.Buffer)
	err = png.Encode(encoded, img)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(encoded.Bytes()), nil
}
