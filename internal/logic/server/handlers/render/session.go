package render

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"image/png"
	"strconv"
	"sync"

	"github.com/cufee/aftermath-core/dataprep"
	"github.com/cufee/aftermath-core/internal/core/database"
	"github.com/cufee/aftermath-core/internal/core/database/models"
	"github.com/cufee/aftermath-core/internal/core/localization"
	"github.com/cufee/aftermath-core/internal/core/server"
	"github.com/cufee/aftermath-core/internal/logic/cache"
	"github.com/cufee/aftermath-core/internal/logic/content"
	"github.com/cufee/aftermath-core/internal/logic/render"
	"github.com/cufee/aftermath-core/internal/logic/render/assets"
	"github.com/cufee/aftermath-core/internal/logic/render/session"
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

	imageData, err := getEncodedSessionImage(utils.RealmFromAccountID(accountId), accountId, types.RenderRequestPayload{TankLimit: 5})
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

	imageData, err := getEncodedSessionImage(utils.RealmFromAccountID(accountId), accountId, types.RenderRequestPayload{TankLimit: 5})
	if err != nil {
		return c.Status(500).JSON(server.NewErrorResponseFromError(err, "getEncodedSessionImage"))
	}

	return c.JSON(server.NewResponse(imageData))
}

func getEncodedSessionImage(realm string, accountId int, options types.RenderRequestPayload) (string, error) {
	sessionData, err := stats.GetCurrentPlayerSession(realm, accountId)
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

	averages, err := stats.GetVehicleAverages(sessionData.Diff.Vehicles)
	if err != nil {
		return "", err
	}

	subscriptions, err := database.FindActiveSubscriptionsByReferenceIDs(fmt.Sprint(sessionData.Account.ID), fmt.Sprint(sessionData.Account.ClanID))
	if err != nil && !errors.Is(err, database.ErrSubscriptionNotFound) {
		log.Warn().Err(err).Msg("failed to get subscriptions")
		// We can continue without subscriptions
	}

	// Fetch the background image in a separate goroutine
	var imageWg sync.WaitGroup
	backgroundChan := make(chan image.Image, 1)
	imageWg.Add(1)
	go func() {
		defer imageWg.Done()
		backgrounds, err := database.GetContentByReferenceIDs[[]string](models.UserContentTypeBackground, fmt.Sprint(sessionData.Account.ID), fmt.Sprint(sessionData.Account.ClanID))
		if err != nil {
			bgImage, _ := assets.GetImage("images/backgrounds/default")
			backgroundChan <- bgImage
			return
		}
		// Try to load the first background image
		for _, c := range backgrounds.Data {
			if c != "" {
				image, _, err := content.LoadRemoteImage(c)
				if err == nil && image != nil {
					backgroundChan <- image
					return
				}
			}
		}
		// fallback
		bgImage, _ := assets.GetImage("images/backgrounds/default")
		backgroundChan <- bgImage
	}()

	sortOptions := stats.SortOptions{
		Limit: options.TankLimit,
		By:    stats.ParseSortOptions(options.SortBy),
	}
	if sortOptions.Limit == 0 {
		sortOptions.Limit = 5
	}

	statsCards, err := dataprep.SnapshotToSession(dataprep.ExportInput{
		SessionStats:          sessionData.Diff,
		CareerStats:           sessionData.Selected,
		SessionVehicles:       stats.SortVehicles(sessionData.Diff.Vehicles, averages, sortOptions),
		GlobalVehicleAverages: averages,
	}, dataprep.ExportOptions{
		Blocks: dataprep.DefaultBlockPresets,
		Locale: localization.LanguageEN,
	})
	if err != nil {
		return "", err
	}

	player := session.PlayerData{
		Clan:          &sessionData.Account.Clan,
		Account:       &sessionData.Account.Account,
		Subscriptions: subscriptions,
		Cards:         statsCards,
	}

	renderOptions := session.RenderOptions{
		PromoText: []string{"Aftermath is back!", "amth.one/join  |  amth.one/invite"},
		CardStyle: session.DefaultCardStyle(nil),
	}

	cards, err := session.RenderStatsImage(player, renderOptions)
	if err != nil {
		return "", err
	}

	imageWg.Wait()
	close(backgroundChan)
	bgImage := <-backgroundChan

	img := render.AddBackground(cards, bgImage, render.Style{Blur: 10, BorderRadius: 30, BackgroundColor: render.DiscordBackgroundColor})

	encoded := new(bytes.Buffer)
	err = png.Encode(encoded, img)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(encoded.Bytes()), nil
}
