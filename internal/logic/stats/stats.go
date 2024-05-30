package stats

import (
	"errors"
	"fmt"
	"strconv"
	"sync"

	"github.com/cufee/aftermath-core/internal/core/stats"
	"github.com/cufee/aftermath-core/internal/core/utils"
	"github.com/cufee/aftermath-core/internal/core/wargaming"
	"github.com/rs/zerolog/log"

	"github.com/cufee/am-wg-proxy-next/v2/client"
	wg "github.com/cufee/am-wg-proxy-next/v2/types"
)

type CompleteStats struct {
	Session  stats.SessionSnapshot
	Account  wg.ExtendedAccount
	Clan     wg.ClanMember
	Vehicles []wg.VehicleStatsFrame
}

type vehiclesWithAccount struct {
	accountID int
	vehicles  []wg.VehicleStatsFrame
}

var (
	ErrBlankResponse     = errors.New("blank response from wargaming")
	ErrTooManyAccountIDs = errors.New("too many account IDs")
)

func GetLastBattleTimes(realm string, accountIDs ...int) (map[int]int, error) {
	if len(accountIDs) == 0 {
		return make(map[int]int), nil
	}
	if len(accountIDs) > 100 {
		return nil, ErrTooManyAccountIDs
	}

	ids := make([]string, len(accountIDs))
	for i, id := range accountIDs {
		ids[i] = fmt.Sprintf("%d", id)
	}

	players, err := wargaming.Clients.Live.BatchAccountByID(realm, ids, "account_id", "last_battle_time")
	if err != nil {
		return nil, err
	}

	lastBattleTimes := make(map[int]int, len(players))
	for _, player := range players {
		lastBattleTimes[player.ID] = player.LastBattleTime
	}

	return lastBattleTimes, nil
}

func GetCompleteStatsWithClient(client client.Client, realm string, accountIDs ...int) (map[int]utils.DataWithError[*CompleteStats], error) {
	if len(accountIDs) > 100 {
		return nil, ErrTooManyAccountIDs
	}

	var waitGroup sync.WaitGroup

	accountClansChan := make(chan utils.DataWithError[map[string]wg.ClanMember], 1)
	accountChan := make(chan utils.DataWithError[map[string]wg.ExtendedAccount], 1)
	vehiclesChan := make(chan utils.DataWithError[vehiclesWithAccount], len(accountIDs))

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()

		// Convert ints to strings
		accountIDsString := make([]string, len(accountIDs))
		for i, accountID := range accountIDs {
			accountIDsString[i] = fmt.Sprintf("%d", accountID)
		}

		accounts, err := client.BatchAccountByID(realm, accountIDsString)
		if err != nil {
			log.Err(err).Msg("failed to get accounts")
		}
		accountChan <- utils.DataWithError[map[string]wg.ExtendedAccount]{Data: accounts, Err: err}
	}()

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()

		// Convert ints to strings
		accountIDsString := make([]string, len(accountIDs))
		for i, accountID := range accountIDs {
			accountIDsString[i] = fmt.Sprintf("%d", accountID)
		}

		clans, err := client.BatchAccountClan(realm, accountIDsString)
		if err != nil {
			// This is not a critical error, so we don't return it
			log.Err(err).Msg("failed to get accounts clans")
		}
		accountClansChan <- utils.DataWithError[map[string]wg.ClanMember]{Data: clans, Err: nil}
	}()

	// There is not endpoint to get vehicles for multiple accounts, so we have to do it one by one
	for _, accountID := range accountIDs {
		waitGroup.Add(1)

		go func(id int) {
			defer waitGroup.Done()

			accountVehicles, err := client.AccountVehicles(realm, strconv.Itoa(id))
			if err != nil {
				log.Err(err).Msg("failed to get account vehicles")
				vehiclesChan <- utils.DataWithError[vehiclesWithAccount]{Err: err}
				return
			}
			data := vehiclesWithAccount{accountID: id, vehicles: accountVehicles}
			vehiclesChan <- utils.DataWithError[vehiclesWithAccount]{Data: data}
		}(accountID)
	}

	waitGroup.Wait()
	close(accountChan)
	close(vehiclesChan)
	close(accountClansChan)

	accounts := <-accountChan
	if accounts.Err != nil {
		return nil, accounts.Err
	}

	accountClans := <-accountClansChan
	if accountClans.Err != nil {
		return nil, accountClans.Err
	}

	sessions := make(map[int]utils.DataWithError[*CompleteStats], len(accountIDs))
	for vehicle := range vehiclesChan {
		if vehicle.Err != nil {
			sessions[vehicle.Data.accountID] = utils.DataWithError[*CompleteStats]{Err: vehicle.Err}
			continue
		}

		account, ok := accounts.Data[fmt.Sprintf("%d", vehicle.Data.accountID)]
		if !ok || account.ID == 0 {
			sessions[vehicle.Data.accountID] = utils.DataWithError[*CompleteStats]{
				Err: errors.New("account invalid or private"),
			}
			continue
		}

		clan := accountClans.Data[fmt.Sprintf("%d", vehicle.Data.accountID)]
		session := CompleteStatsFromWargaming(account, vehicle.Data.vehicles)

		sessions[vehicle.Data.accountID] = utils.DataWithError[*CompleteStats]{
			Data: &CompleteStats{
				Clan:     clan,
				Session:  session,
				Account:  account,
				Vehicles: vehicle.Data.vehicles,
			},
		}
	}

	return sessions, nil
}
