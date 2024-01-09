package cache

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/cufee/aftermath-core/internal/core/database"
	"github.com/cufee/aftermath-core/internal/core/utils"
	wg "github.com/cufee/am-wg-proxy-next/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type DatabaseAccount struct {
	ID       int    `json:"id" bson:"_id"`
	Realm    string `json:"realm" bson:"realm"`
	Nickname string `json:"nickname" bson:"nickname"`

	Clan *DatabaseAccountClan `json:"clan" bson:"clan"`

	LastBattleTime time.Time `json:"lastBattleTime" bson:"lastBattleTime"` // This will probably end up not being updated too often

	LastUpdated time.Time `json:"lastUpdated" bson:"lastUpdated"` // This will be updated every time the account is updated
}

type DatabaseAccountClan struct {
	ID       int       `json:"id" bson:"_id"`
	Role     string    `json:"role" bson:"role"`
	JoinedAt time.Time `json:"joinedAt" bson:"joinedAt"`
}

func GetPlayerAccount(id int) (*DatabaseAccount, error) {
	ctx, cancel := database.DefaultClient.Ctx()
	defer cancel()

	var account DatabaseAccount
	err := database.DefaultClient.Collection(database.CollectionAccounts).FindOne(ctx, bson.M{"_id": id}).Decode(&account)
	if err != nil {
		return nil, err
	}

	return &account, nil
}

func UpdatePlayerAccounts(accounts ...DatabaseAccount) error {
	var writes []mongo.WriteModel
	for _, account := range accounts {
		model := mongo.NewUpdateOneModel()
		model.SetFilter(bson.M{"_id": account.ID})
		model.SetUpdate(bson.M{"$set": account})
		model.SetUpsert(true)
		writes = append(writes, model)
	}

	ctx, cancel := database.DefaultClient.Ctx()
	defer cancel()

	_, err := database.DefaultClient.Collection(database.CollectionAccounts).BulkWrite(ctx, writes)
	if err != nil {
		return err
	}

	return nil
}

/*
UpdateRealmAccountsCache updates all active accounts for a realm in the cache.
*/
func UpdateRealmAccountsCache(realm string) error {
	ctx, cancel := database.DefaultClient.Ctx()
	defer cancel()

	cur, err := database.DefaultClient.Collection(database.CollectionAccounts).Find(ctx, bson.M{"realm": realm})
	if err != nil {
		return err
	}

	var accounts []DatabaseAccount
	err = cur.All(ctx, &accounts)
	if err != nil {
		return err
	}

	var accountIDs []int
	for _, account := range accounts {
		accountIDs = append(accountIDs, account.ID)
	}

	return UpdateAccountCache(realm, accountIDs)
}

/*
UpdateAccountCache updates active accounts in the cache.
*/
func UpdateAccountCache(realm string, accountIDs []int) error {
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

			accounts, err := cacheClient.BulkGetAccountsByID(accountIDsString, realm)
			accountsChan <- utils.DataWithError[map[string]wg.ExtendedAccount]{Data: accounts, Err: err}
		}(batch)
	}

	waitGroup.Wait()
	close(accountsChan)

	var documents []DatabaseAccount
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

	return UpdatePlayerAccounts(documents...)
}

func accountToDatabaseAccount(realm string, acc wg.ExtendedAccount) *DatabaseAccount {
	return &DatabaseAccount{
		ID:       acc.ID,
		Realm:    strings.ToUpper(realm),
		Nickname: acc.Nickname,
		// acc.LastBattleTime is unix timestamp
		LastBattleTime: time.Unix(int64(acc.LastBattleTime), 0),
		LastUpdated:    time.Now(),
	}
}
