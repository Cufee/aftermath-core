package render

import (
	"errors"
	"fmt"
	"image"
	"strconv"
	"sync"

	dataprep "github.com/cufee/aftermath-core/dataprep/period"
	"golang.org/x/text/language"

	"github.com/cufee/aftermath-core/internal/core/database"
	"github.com/cufee/aftermath-core/internal/core/database/models"
	"github.com/cufee/aftermath-core/internal/core/localization"
	"github.com/cufee/aftermath-core/internal/core/server"
	core "github.com/cufee/aftermath-core/internal/core/utils"
	"github.com/cufee/aftermath-core/internal/logic/content"
	renderCore "github.com/cufee/aftermath-core/internal/logic/render"
	"github.com/cufee/aftermath-core/internal/logic/render/assets"
	render "github.com/cufee/aftermath-core/internal/logic/render/period"
	"github.com/cufee/aftermath-core/internal/logic/stats/period"
	"github.com/cufee/aftermath-core/types"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func PeriodFromIDHandler(c *fiber.Ctx) error {
	account := c.Params("account")
	accountId, err := strconv.Atoi(account)
	if err != nil {
		return c.Status(400).JSON(server.NewErrorResponseFromError(err, "strconv.Atoi"))
	}

	var opts types.PeriodRequestPayload
	err = c.BodyParser(&opts)
	if err != nil {
		return c.Status(400).JSON(server.NewErrorResponseFromError(err, "c.BodyParser"))
	}

	imageData, err := getEncodedPeriodImage(accountId, opts)
	if err != nil {
		return c.Status(500).JSON(server.NewErrorResponseFromError(err, "getEncodedSessionImage"))
	}

	return c.JSON(server.NewResponse(imageData))
}

func PeriodFromUserHandler(c *fiber.Ctx) error {
	user := c.Params("id")
	if user == "" {
		return c.Status(400).JSON(server.NewErrorResponse("id path parameter is required", "c.Param"))
	}

	var opts types.PeriodRequestPayload
	err := c.BodyParser(&opts)
	if err != nil {
		return c.Status(400).JSON(server.NewErrorResponseFromError(err, "c.BodyParser"))
	}

	connection, err := database.FindUserConnection(user, models.ConnectionTypeWargaming)
	if err != nil {
		if errors.Is(err, database.ErrConnectionNotFound) {
			return c.Status(404).JSON(server.NewErrorResponseFromError(err, "models.FindUserConnection"))
		}
		return c.Status(500).JSON(server.NewErrorResponseFromError(err, "models.FindUserConnection"))
	}

	accountId, err := strconv.Atoi(connection.ExternalID)
	if err != nil {
		return c.Status(500).JSON(server.NewErrorResponse("invalid connection", "strconv.Atoi"))
	}

	imageData, err := getEncodedPeriodImage(accountId, opts)
	if err != nil {
		return c.Status(500).JSON(server.NewErrorResponseFromError(err, "getEncodedSessionImage"))
	}

	return c.JSON(server.NewResponse(imageData))
}

func getEncodedPeriodImage(accountId int, options types.PeriodRequestPayload) (string, error) {
	stats, err := period.GetPlayerStats(accountId, options.Days)
	if err != nil {
		return "", err
	}

	// Fetch the background image in a separate goroutine
	var wait sync.WaitGroup
	backgroundChan := make(chan image.Image, 1)
	cardsChan := make(chan core.DataWithError[image.Image], 1)

	wait.Add(1)
	go func() {
		defer wait.Done()

		referenceIDs := []string{fmt.Sprint(stats.Account.ID), fmt.Sprint(stats.Clan.ID)}
		backgrounds, err := database.GetContentByReferenceIDs[string](referenceIDs, models.UserContentTypePersonalBackground, models.UserContentTypeClanBackground)
		if err != nil {
			log.Warn().Err(err).Msg("failed to get backgrounds")
			bgImage, _ := assets.GetImage("images/backgrounds/default")
			backgroundChan <- bgImage
			return
		}

		// We should get personal image over clan image when possible, fallback to default
		for _, id := range referenceIDs {
			for _, c := range backgrounds {
				if c.Data != "" && c.ReferenceID == id {
					image, _, err := content.LoadRemoteImage(c.Data)
					if err == nil && image != nil {
						backgroundChan <- image
						return
					}
				}
			}
		}
		// fallback
		bgImage, _ := assets.GetImage("images/backgrounds/default")
		backgroundChan <- bgImage
	}()

	wait.Add(1)
	go func() {
		defer wait.Done()

		// Find a user who has a verified connection for this account
		var referenceIds []string = []string{fmt.Sprint(stats.Account.ID), fmt.Sprint(stats.Clan.ID)}
		connections, err := database.FindConnectionsByReferenceID(fmt.Sprint(stats.Account.ID), models.ConnectionTypeWargaming)
		if err != nil && !errors.Is(err, database.ErrConnectionNotFound) {
			log.Warn().Err(err).Msg("failed to get connection")
			// We can continue without connections
		}
		for _, connection := range connections {
			if connection.Metadata["verified"] == true {
				referenceIds = append([]string{connection.UserID}, referenceIds...)
			}
		}

		var vehicleIDs []int
		for _, vehicle := range stats.Vehicles {
			vehicleIDs = append(vehicleIDs, vehicle.VehicleID)
		}
		vehiclesGlossary, err := database.GetGlossaryVehicles(vehicleIDs...)
		if err != nil {
			// This is definitely not fatal, but will look ugly
			log.Warn().Err(err).Msg("failed to get vehicles glossary")
		}

		subscriptions, err := database.FindActiveSubscriptionsByReferenceIDs(referenceIds...)
		if err != nil && !errors.Is(err, database.ErrSubscriptionNotFound) {
			log.Warn().Err(err).Msg("failed to get subscriptions")
			// We can continue without subscriptions
		}

		cards, err := dataprep.SnapshotToSession(
			dataprep.ExportInput{
				Stats:           stats,
				VehicleGlossary: vehiclesGlossary,
			}, dataprep.ExportOptions{
				Blocks:        dataprep.DefaultBlocks,
				Locale:        language.English,
				LocalePrinter: localization.GetPrinter(language.English),
				Highlights:    dataprep.DefaultHighlights,
			})
		if err != nil {
			cardsChan <- core.DataWithError[image.Image]{Err: err}
		}

		renderOptions := render.RenderOptions{}

		img, err := render.RenderImage(render.PlayerData{
			Stats:         stats,
			Cards:         cards,
			Subscriptions: subscriptions,
		}, renderOptions)
		cardsChan <- core.DataWithError[image.Image]{Data: img, Err: err}
	}()

	wait.Wait()
	close(cardsChan)
	close(backgroundChan)

	cards := <-cardsChan
	if cards.Err != nil {
		return "", cards.Err
	}

	bgImage := <-backgroundChan
	img := renderCore.AddBackground(cards.Data, bgImage, renderCore.Style{Blur: 10, BorderRadius: 30})
	return core.EncodeImage(img)
}
