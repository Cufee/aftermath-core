package stats

import (
	"errors"
	"strconv"
	"time"

	"github.com/cufee/aftermath-core/dataprep/session"
	"github.com/cufee/aftermath-core/internal/core/database"
	"github.com/cufee/aftermath-core/internal/core/database/models"
	"github.com/cufee/aftermath-core/internal/core/localization"
	"github.com/cufee/aftermath-core/types"
	"github.com/cufee/am-wg-proxy-next/v2/utils"
	"golang.org/x/text/language"

	"github.com/cufee/aftermath-core/dataprep"
	"github.com/cufee/aftermath-core/internal/core/server"
	"github.com/cufee/aftermath-core/internal/logic/cache"
	"github.com/cufee/aftermath-core/internal/logic/stats"
	"github.com/cufee/aftermath-core/internal/logic/stats/sessions"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func RecordPlayerSession(c *fiber.Ctx) error {
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

	accountErrs, err := cache.RefreshSessionsAndAccounts(opts.Type(), opts.ReferenceID, utils.RealmFromPlayerID(accountId), accountId)
	if err != nil {
		return c.Status(500).JSON(server.NewErrorResponseFromError(err, "cache.RefreshSessionsAndAccounts"))
	}
	if accountErrs[accountId] != nil {
		return c.Status(500).JSON(server.NewErrorResponseFromError(err, "cache.RefreshSessionsAndAccounts"))
	}

	return c.JSON(server.NewResponse(err))
}

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

	stats, err := getSessionStats(accountId, opts)
	if err != nil {
		return c.Status(500).JSON(server.NewErrorResponseFromError(err, "getEncodedSessionImage"))
	}

	return c.JSON(server.NewResponse(stats))
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
			return c.Status(404).JSON(server.NewErrorResponseFromError(err, "users.FindUserConnection"))
		}
		return c.Status(500).JSON(server.NewErrorResponseFromError(err, "users.FindUserConnection"))
	}

	accountId, err := strconv.Atoi(connection.ExternalID)
	if err != nil {
		return c.Status(500).JSON(server.NewErrorResponse("invalid connection", "strconv.Atoi"))
	}

	stats, err := getSessionStats(accountId, opts)
	if err != nil {
		return c.Status(500).JSON(server.NewErrorResponseFromError(err, "getSessionStats"))
	}

	return c.JSON(server.NewResponse(stats))
}

func getSessionStats(accountId int, opts types.SessionRequestPayload) (*session.SessionStats, error) {
	realm := utils.RealmFromPlayerID(accountId)

	blocks, err := dataprep.ParseTags(opts.Presets...)
	if err != nil {
		blocks = session.DefaultSessionBlocks
	}

	now := int(time.Now().Unix())
	playerSession, err := sessions.GetCurrentPlayerSession(accountId, database.SessionGetOptions{LastBattleBefore: &now, ReferenceID: opts.ReferenceID})
	if err != nil {
		if !errors.Is(err, sessions.ErrNoSessionCached) {
			return nil, err
		}
		// Refresh the session cache in the background
		go func(realm string, accountId int) {
			accountErrs, err := cache.RefreshSessionsAndAccounts(models.SessionTypeDaily, opts.ReferenceID, realm, accountId)
			if err != nil || len(accountErrs) > 0 {
				log.Err(err).Msg("failed to refresh session cache")
			}
		}(realm, accountId)
	}

	if playerSession.Account.ClanID != 0 {
		go func() {
			err := cache.CacheAllNewClanMembers(realm, playerSession.Account.ClanID)
			if err != nil {
				log.Err(err).Msg("failed to cache new clan members")
			}
		}()
	}

	var vehicleIDs []int
	for _, vehicle := range playerSession.Diff.Vehicles {
		vehicleIDs = append(vehicleIDs, vehicle.VehicleID)
	}
	for _, vehicle := range playerSession.Selected.Vehicles {
		vehicleIDs = append(vehicleIDs, vehicle.VehicleID)
	}

	averages, err := database.GetVehicleAverages(vehicleIDs...)
	if err != nil {
		return nil, err
	}

	vehiclesGlossary, err := database.GetGlossaryVehicles(vehicleIDs...)
	if err != nil {
		// This is definitely not fatal, but will look ugly
		log.Warn().Err(err).Msg("failed to get vehicles glossary")
	}

	unratedVehicles, ratingVehicles := stats.SortAndSplitVehicles(playerSession.Diff.Vehicles, averages, stats.SortOptions{By: stats.SortByLastBattle, Limit: 5}, stats.SortOptions{By: stats.SortByLastBattle, Limit: 3})

	statsCards, err := session.SnapshotToSession(session.ExportInput{
		SessionStats:           playerSession.Diff,
		CareerStats:            playerSession.Selected,
		SessionUnratedVehicles: unratedVehicles,
		SessionRatingVehicles:  ratingVehicles,

		VehicleGlossary:       vehiclesGlossary,
		GlobalVehicleAverages: averages,
	}, session.ExportOptions{
		Locale:                language.English,
		LocalePrinter:         localization.GetPrinter(language.English),
		Blocks:                blocks,
		IncludeRatingVehicles: true,
	})
	if err != nil {
		return nil, err
	}

	return &session.SessionStats{
		Realm:      realm,
		Locale:     language.English.String(),
		LastBattle: playerSession.Account.LastBattleTime,
		Clan:       playerSession.Account.ClanMember.Clan,
		Account:    playerSession.Account.Account,
		Cards:      statsCards,
	}, nil
}
