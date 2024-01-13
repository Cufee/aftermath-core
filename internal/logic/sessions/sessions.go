package sessions

import (
	"errors"
	"fmt"
	"sync"

	"github.com/cufee/aftermath-core/internal/core/stats"
	"github.com/cufee/aftermath-core/internal/core/utils"

	client "github.com/cufee/am-wg-proxy-next/client"
	wg "github.com/cufee/am-wg-proxy-next/types"
)

type SessionWithRawData struct {
	Session  *stats.SessionSnapshot
	Account  *wg.ExtendedAccount
	Clan     *wg.ClanMember
	Vehicles []wg.VehicleStatsFrame
}

type vehiclesWithAccount struct {
	accountID int
	vehicles  []wg.VehicleStatsFrame
}

var (
	ErrTooManyAccountIDs = errors.New("too many account IDs")
)

func GetLiveSessions(realm string, accountIDs ...int) (map[int]*SessionWithRawData, error) {
	return GetSessionsWithClient(liveClient, realm, accountIDs...)
}

func GetSessionsWithClient(client *client.Client, realm string, accountIDs ...int) (map[int]*SessionWithRawData, error) {
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

		accounts, err := client.BulkGetAccountsByID(accountIDsString, realm)
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

		clans, err := client.BulkGetAccountsClans(accountIDsString, realm)
		accountClansChan <- utils.DataWithError[map[string]wg.ClanMember]{Data: clans, Err: err}
	}()

	// There is not endpoint to get vehicles for multiple accounts, so we have to do it one by one
	for _, accountID := range accountIDs {
		waitGroup.Add(1)

		go func(id int) {
			defer waitGroup.Done()
			accountVehicles, err := client.GetAccountVehicles(id)
			if err != nil {
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

	sessions := make(map[int]*SessionWithRawData, len(accountIDs))
	for vehicle := range vehiclesChan {
		if vehicle.Err != nil {
			return nil, vehicle.Err
		}

		account, ok := accounts.Data[fmt.Sprintf("%d", vehicle.Data.accountID)]
		if !ok || account.ID == 0 {
			return nil, fmt.Errorf("account %d not found", vehicle.Data.accountID)
		}

		clan := accountClans.Data[fmt.Sprintf("%d", vehicle.Data.accountID)]
		session := AccountStatsToSession(account, vehicle.Data.vehicles)

		sessions[vehicle.Data.accountID] = &SessionWithRawData{
			Clan:     &clan,
			Session:  session,
			Account:  &account,
			Vehicles: vehicle.Data.vehicles,
		}
	}

	return sessions, nil
}
