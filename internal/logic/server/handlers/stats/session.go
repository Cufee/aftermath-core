package stats

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/cufee/aftermath-core/dataprep/session"
	"github.com/cufee/aftermath-core/internal/core/database"
	"github.com/cufee/aftermath-core/internal/core/database/models"
	"github.com/cufee/aftermath-core/types"

	"github.com/cufee/aftermath-core/dataprep"
	"github.com/cufee/aftermath-core/internal/core/localization"
	"github.com/cufee/aftermath-core/internal/core/server"
	"github.com/cufee/aftermath-core/internal/logic/cache"
	"github.com/cufee/aftermath-core/internal/logic/stats"
	"github.com/cufee/aftermath-core/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func RecordPlayerSession(c *fiber.Ctx) error {
	account := c.Params("account")
	accountId, err := strconv.Atoi(account)
	if err != nil {
		return c.Status(400).JSON(server.NewErrorResponseFromError(err, "strconv.Atoi"))
	}

	sessionType := models.ParseSessionType(c.Query("type"))
	accountErrs, err := cache.RefreshSessionsAndAccounts(sessionType, utils.RealmFromAccountID(accountId), accountId)
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

	stats, err := getSessionStats(utils.RealmFromAccountID(accountId), accountId, types.RenderRequestPayload{Presets: strings.Split(c.Query("blocks"), ",")})
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

	stats, err := getSessionStats(utils.RealmFromAccountID(accountId), accountId, types.RenderRequestPayload{Presets: strings.Split(c.Query("blocks"), ",")})
	if err != nil {
		return c.Status(500).JSON(server.NewErrorResponseFromError(err, "getSessionStats"))
	}

	return c.JSON(server.NewResponse(stats))
}

func getSessionStats(realm string, accountId int, opts types.RenderRequestPayload) (*session.SessionStats, error) {
	blocks, err := dataprep.ParseTags(opts.Presets...)
	if err != nil {
		blocks = session.DefaultSessionBlocks
	}

	now := int(time.Now().Unix())
	playerSession, err := stats.GetCurrentPlayerSession(realm, accountId, database.SessionGetOptions{LastBattleBefore: &now})
	if err != nil {
		return nil, err
	}

	if playerSession.Account.ClanID != 0 {
		go func() {
			err := cache.CacheAllNewClanMembers(realm, playerSession.Account.ClanID)
			if err != nil {
				log.Err(err).Msg("failed to cache new clan members")
			}
		}()
	}
	averages, err := stats.GetVehicleAverages(playerSession.Diff.Vehicles)
	if err != nil {
		return nil, err
	}

	sortOptions := stats.SortOptions{
		By:    stats.SortByLastBattle,
		Limit: 5,
	}
	statsCards, err := session.SnapshotToSession(session.ExportInput{
		SessionStats:          playerSession.Diff,
		CareerStats:           playerSession.Selected,
		SessionVehicles:       stats.SortVehicles(playerSession.Diff.Vehicles, averages, sortOptions),
		GlobalVehicleAverages: averages,
	}, session.ExportOptions{
		Locale: localization.LanguageEN,
		Blocks: blocks,
	})
	if err != nil {
		return nil, err
	}

	return &session.SessionStats{
		Realm:      realm,
		Locale:     localization.LanguageEN.WargamingCode,
		LastBattle: playerSession.Account.LastBattleTime,
		Clan:       playerSession.Account.ClanMember.Clan,
		Account:    playerSession.Account.Account,
		Cards:      statsCards,
	}, nil
}
