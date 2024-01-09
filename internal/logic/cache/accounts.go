package cache

import (
	"time"

	"github.com/cufee/aftermath-core/internal/core/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type DatabaseAccount struct {
	ID       int    `json:"id" bson:"_id"`
	Realm    string `json:"realm" bson:"realm"`
	Nickname string `json:"nickname" bson:"nickname"`

	Clan *DatabaseAccountClan `json:"clan" bson:"clan"`

	LastBattleTime time.Time `json:"lastBattleTime" bson:"lastBattleTime"` // This will probably end up not being updated too often
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
