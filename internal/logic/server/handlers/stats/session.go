package stats

import (
	"errors"
	"strconv"

	"github.com/cufee/aftermath-core/internal/core/database"
	"github.com/cufee/aftermath-core/internal/core/database/models"

	"github.com/cufee/aftermath-core/dataprep"
	"github.com/cufee/aftermath-core/internal/core/localization"
	"github.com/cufee/aftermath-core/internal/core/server"
	"github.com/cufee/aftermath-core/internal/logic/cache"
	"github.com/cufee/aftermath-core/internal/logic/stats"
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

	stats, err := getSessionStats(utils.RealmFromAccountID(accountId), accountId)
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

	stats, err := getSessionStats(utils.RealmFromAccountID(accountId), accountId)
	if err != nil {
		return c.Status(500).JSON(server.NewErrorResponseFromError(err, "getSessionStats"))
	}

	return c.JSON(server.NewResponse(stats))
}

func getSessionStats(realm string, accountId int) (*dataprep.SessionStats, error) {
	session, err := stats.GetCurrentPlayerSession(realm, accountId)
	if err != nil {
		return nil, err
	}

	if session.Account.ClanID != 0 {
		go func() {
			err := cache.CacheAllNewClanMembers(realm, session.Account.ClanID)
			if err != nil {
				log.Err(err).Msg("failed to cache new clan members")
			}
		}()
	}
	averages, err := stats.GetVehicleAverages(session.Diff.Vehicles)
	if err != nil {
		return nil, err
	}

	sortOptions := stats.SortOptions{
		By:    stats.SortByLastBattle,
		Limit: 5,
	}
	statsCards, err := dataprep.SnapshotToSession(dataprep.ExportInput{
		SessionStats:          session.Diff,
		CareerStats:           session.Selected,
		SessionVehicles:       stats.SortVehicles(session.Diff.Vehicles, averages, sortOptions),
		GlobalVehicleAverages: averages,
	}, dataprep.ExportOptions{
		Blocks: dataprep.DefaultBlockPresets,
		Locale: localization.LanguageEN,
	})
	if err != nil {
		return nil, err
	}

	return &dataprep.SessionStats{
		Realm:      realm,
		Locale:     localization.LanguageEN.WargamingCode,
		LastBattle: session.Account.LastBattleTime,
		Clan:       session.Account.ClanMember.Clan,
		Account:    session.Account.Account,
		Cards:      statsCards,
	}, nil
}
