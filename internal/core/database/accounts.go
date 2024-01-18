package database

import (
	"github.com/cufee/aftermath-core/internal/core/database/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func UpdatePlayerAccounts(accounts ...models.Account) error {
	var writes []mongo.WriteModel
	for _, account := range accounts {
		model := mongo.NewUpdateOneModel()
		model.SetFilter(bson.M{"_id": account.ID})
		model.SetUpdate(bson.M{"$set": account})
		model.SetUpsert(true)
		writes = append(writes, model)
	}

	ctx, cancel := DefaultClient.Ctx()
	defer cancel()

	_, err := DefaultClient.Collection(CollectionAccounts).BulkWrite(ctx, writes)
	if err != nil {
		return err
	}

	return nil
}

func GetPlayerAccount(id int) (*models.Account, error) {
	ctx, cancel := DefaultClient.Ctx()
	defer cancel()

	var account models.Account
	err := DefaultClient.Collection(CollectionAccounts).FindOne(ctx, bson.M{"_id": id}).Decode(&account)
	if err != nil {
		return nil, err
	}

	return &account, nil
}

func GetRealmAccountIDs(realm string) ([]int, error) {
	ctx, cancel := DefaultClient.Ctx()
	defer cancel()

	cur, err := DefaultClient.Collection(CollectionAccounts).Find(ctx, bson.M{"realm": realm})
	if err != nil {
		return nil, err
	}

	var accounts []models.Account
	err = cur.All(ctx, &accounts)
	if err != nil {
		return nil, err
	}

	var accountIDs []int
	for _, account := range accounts {
		accountIDs = append(accountIDs, account.ID)
	}

	return accountIDs, nil
}
