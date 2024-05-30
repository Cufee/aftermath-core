package cache

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/cufee/aftermath-core/internal/core/database"
	"github.com/cufee/aftermath-core/internal/core/database/models"
	"github.com/cufee/aftermath-core/internal/core/utils"
	"github.com/cufee/aftermath-core/internal/core/wargaming"
	wg "github.com/cufee/am-wg-proxy-next/v2/types"
)

func CacheAllNewClanMembers(realm string, clanId int) error {
	clan, err := wargaming.Clients.Cache.ClanByID(realm, strconv.Itoa(clanId))
	if err != nil {
		return err
	}

	lastBattles, err := database.GetLastBattleTimes(models.SessionTypeDaily, nil, clan.MembersIDS...)
	if err != nil {
		return err
	}

	var newAccounts []int
	for _, member := range clan.MembersIDS {
		if _, ok := lastBattles[member]; !ok {
			newAccounts = append(newAccounts, member)
		}
	}
	if len(newAccounts) == 0 {
		return nil
	}

	_, err = RefreshSessionsAndAccounts(models.SessionTypeDaily, nil, realm, newAccounts...)
	if err != nil {
		return err
	}
	return nil
}

func UpdatePlayerAccountsFromWG(realm string, accounts ...wg.ExtendedAccount) error {
	var converted []models.Account
	for _, account := range accounts {
		if account.ID == 0 {
			continue
		}
		converted = append(converted, *accountToDatabaseAccount(realm, account))
	}
	return database.UpdatePlayerAccounts(converted...)
}

/*
UpdateRealmAccountsCache updates all active accounts for a realm in the cache.
*/
func UpdateRealmAccountsCache(realm string) error {
	accountIDs, err := database.GetRealmAccountIDs(realm)
	if err != nil {
		return err
	}
	return UpdateAccountsCache(realm, accountIDs)
}

/*
UpdateAccountCache updates active accounts in the cache.
*/
func UpdateAccountsCache(realm string, accountIDs []int) error {
	var waitGroup sync.WaitGroup

	batches := utils.BatchAccountIDs(accountIDs, 100)
	accountsChan := make(chan utils.DataWithError[map[string]wg.ExtendedAccount], len(batches))

	for _, batch := range batches {
		waitGroup.Add(1)

		go func(ids []int) {
			defer waitGroup.Done()

			// Convert ints to strings
			accountIDsString := make([]string, len(ids))
			for i, accountID := range ids {
				accountIDsString[i] = fmt.Sprintf("%d", accountID)
			}

			accounts, err := wargaming.Clients.Cache.BatchAccountByID(realm, accountIDsString)
			accountsChan <- utils.DataWithError[map[string]wg.ExtendedAccount]{Data: accounts, Err: err}
		}(batch)
	}

	waitGroup.Wait()
	close(accountsChan)

	var documents []models.Account
	for accounts := range accountsChan {
		if accounts.Err != nil {
			return accounts.Err
		}

		for _, account := range accounts.Data {
			// Skip accounts that haven't played in the last 24 hours to avoid pointless db writes
			lastBattle := time.Unix(int64(account.LastBattleTime), 0)
			if lastBattle.Before(time.Now().Add(-time.Hour * 24)) {
				continue
			}

			documents = append(documents, *accountToDatabaseAccount(realm, account))
		}
	}

	return database.UpdatePlayerAccounts(documents...)
}

func accountToDatabaseAccount(realm string, acc wg.ExtendedAccount) *models.Account {
	return &models.Account{
		ID:       acc.ID,
		Realm:    strings.ToUpper(realm),
		Nickname: acc.Nickname,
		// acc.LastBattleTime is unix timestamp
		LastBattleTime: time.Unix(int64(acc.LastBattleTime), 0),
		LastUpdated:    time.Now(),
	}
}
