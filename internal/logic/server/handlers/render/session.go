package render

import (
	"errors"
	"fmt"
	"image"
	"strconv"
	"sync"

	"github.com/cufee/aftermath-core/dataprep"
	"github.com/cufee/aftermath-core/dataprep/session"
	"github.com/cufee/aftermath-core/internal/core/database"
	"github.com/cufee/aftermath-core/internal/core/database/models"
	"github.com/cufee/aftermath-core/internal/core/localization"
	"github.com/cufee/aftermath-core/internal/core/server"
	core "github.com/cufee/aftermath-core/internal/core/utils"
	"github.com/cufee/aftermath-core/internal/logic/cache"
	"github.com/cufee/aftermath-core/internal/logic/content"
	renderCore "github.com/cufee/aftermath-core/internal/logic/render"
	"github.com/cufee/aftermath-core/internal/logic/render/assets"
	render "github.com/cufee/aftermath-core/internal/logic/render/session"
	"github.com/cufee/aftermath-core/internal/logic/stats"
	"github.com/cufee/aftermath-core/types"
	"github.com/cufee/aftermath-core/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func SessionFromIDHandler(c *fiber.Ctx) error {
	account := c.Params("account")
	accountId, err := strconv.Atoi(account)
	if err != nil {
		return c.Status(400).JSON(server.NewErrorResponseFromError(err, "strconv.Atoi"))
	}

	var opts types.SessionRequestPayload
	err = c.BodyParser(&opts)
	if err != nil {
		return c.Status(400).JSON(server.NewErrorResponseFromError(err, "c.BodyParser"))
	}

	imageData, err := getEncodedSessionImage(utils.RealmFromAccountID(accountId), accountId, opts)
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

	var opts types.SessionRequestPayload
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

	imageData, err := getEncodedSessionImage(utils.RealmFromAccountID(accountId), accountId, opts)
	if err != nil {
		return c.Status(500).JSON(server.NewErrorResponseFromError(err, "getEncodedSessionImage"))
	}

	return c.JSON(server.NewResponse(imageData))
}

func getEncodedSessionImage(realm string, accountId int, options types.SessionRequestPayload) (string, error) {
	blocks, err := dataprep.ParseTags(options.Presets...)
	if err != nil {
		blocks = session.DefaultSessionBlocks
	}

	sessionData, err := stats.GetCurrentPlayerSession(realm, accountId, database.SessionGetOptions{Type: options.Type(), ReferenceID: options.ReferenceID})
	if err != nil {
		return "", err
	}

	if sessionData.Account.ClanID != 0 {
		go func() {
			err := cache.CacheAllNewClanMembers(realm, sessionData.Account.ClanID)
			if err != nil {
				log.Err(err).Msg("failed to cache new clan members")
			}
		}()
	}

	// Fetch the background image in a separate goroutine
	var wait sync.WaitGroup
	backgroundChan := make(chan image.Image, 1)
	cardsChan := make(chan core.DataWithError[image.Image], 1)

	wait.Add(1)
	go func() {
		defer wait.Done()

		referenceIDs := []string{fmt.Sprint(sessionData.Account.ID), fmt.Sprint(sessionData.Account.ClanID)}
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

		averages, err := stats.GetVehicleAverages(sessionData.Diff.Vehicles)
		if err != nil {
			cardsChan <- core.DataWithError[image.Image]{Err: err}
			return
		}

		// Find a user who has a verified connection for this account
		var referenceIds []string = []string{fmt.Sprint(sessionData.Account.ID), fmt.Sprint(sessionData.Account.ClanID)}
		connections, err := database.FindConnectionsByReferenceID(fmt.Sprint(sessionData.Account.ID), models.ConnectionTypeWargaming)
		if err != nil && !errors.Is(err, database.ErrConnectionNotFound) {
			log.Warn().Err(err).Msg("failed to get connection")
			// We can continue without connections
		}
		for _, connection := range connections {
			if connection.Metadata["verified"] == true {
				referenceIds = append([]string{connection.UserID}, referenceIds...)
			}
		}

		subscriptions, err := database.FindActiveSubscriptionsByReferenceIDs(referenceIds...)
		if err != nil && !errors.Is(err, database.ErrSubscriptionNotFound) {
			log.Warn().Err(err).Msg("failed to get subscriptions")
			// We can continue without subscriptions
		}

		sortOptions := stats.SortOptions{
			Limit: options.TankLimit,
			By:    stats.ParseSortOptions(options.SortBy),
		}
		if sortOptions.Limit == 0 {
			sortOptions.Limit = 5
		}

		statsCards, err := session.SnapshotToSession(session.ExportInput{
			SessionStats:          sessionData.Diff,
			CareerStats:           sessionData.Selected,
			SessionVehicles:       stats.SortVehicles(sessionData.Diff.Vehicles, averages, sortOptions),
			GlobalVehicleAverages: averages,
		}, session.ExportOptions{
			Locale: localization.LanguageEN,
			Blocks: blocks,
		})
		if err != nil {
			cardsChan <- core.DataWithError[image.Image]{Err: err}
			return
		}

		player := render.PlayerData{
			Clan:          &sessionData.Account.Clan,
			Account:       &sessionData.Account.Account,
			Subscriptions: subscriptions,
			Cards:         statsCards,
		}

		renderOptions := render.RenderOptions{
			PromoText: []string{"Aftermath is back!", "amth.one/join  |  amth.one/invite"},
			CardStyle: render.DefaultCardStyle(nil),
		}

		cards, err := render.RenderStatsImage(player, renderOptions)
		cardsChan <- core.DataWithError[image.Image]{Data: cards, Err: err}
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
