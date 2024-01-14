package render

import (
	"bytes"
	"encoding/base64"
	"errors"
	"image/png"
	"strconv"

	"github.com/cufee/aftermath-core/internal/core/localization"
	"github.com/cufee/aftermath-core/internal/core/server"
	"github.com/cufee/aftermath-core/internal/logic/render"
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

	// TODO: sorting options and limits
	opts := stats.SortOptions{
		By:    stats.SortByLastBattle,
		Limit: 7,
	}
	vehicles := stats.SortVehicles(session.Diff.Vehicles, opts)

	img, err := render.RenderStatsImage(session, vehicles, localization.LanguageEN)
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
