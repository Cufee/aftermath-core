package users

import (
	"errors"
	"time"

	"github.com/cufee/aftermath-core/internal/core/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type subscriptionType string

const (
	SubscriptionTypePro       = subscriptionType("aftermath-pro")
	SubscriptionTypeProXL     = subscriptionType("aftermath-pro-xl")
	SubscriptionTypePlus      = subscriptionType("aftermath-plus")
	SubscriptionTypeSupporter = subscriptionType("supporter")
)

type UserSubscription struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"-"`
	UserID string             `bson:"userID" json:"userID"`

	SubscriptionType subscriptionType `bson:"subscriptionType" json:"subscriptionType"`
	CreationDate     time.Time        `bson:"creationDate" json:"creationDate"`
	ExpiryDate       time.Time        `bson:"expiryDate" json:"expiryDate"`
}

func FindUserSubscription(userId string, subscriptionType subscriptionType) (*UserSubscription, error) {
	ctx, cancel := database.DefaultClient.Ctx()
	defer cancel()

	var subscription UserSubscription
	err := database.DefaultClient.Collection(database.CollectionUserSubscriptions).FindOne(ctx, bson.M{"userID": userId, "subscriptionType": subscriptionType}).Decode(&subscription)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrConnectionNotFound
		}
		return nil, err
	}

	return &subscription, nil
}

func AddNewUserSubscription(userId string, subscriptionType subscriptionType, expiryDate time.Time) error {
	ctx, cancel := database.DefaultClient.Ctx()
	defer cancel()

	_, err := database.DefaultClient.Collection(database.CollectionUserSubscriptions).InsertOne(ctx, UserSubscription{
		UserID:           userId,
		SubscriptionType: subscriptionType,
		CreationDate:     time.Now(),
		ExpiryDate:       expiryDate,
	})
	if err != nil {
		return err
	}
	return nil
}

func UpdateUserSubscription(subscriptionId primitive.ObjectID, expiryDate time.Time) error {
	ctx, cancel := database.DefaultClient.Ctx()
	defer cancel()

	_, err := database.DefaultClient.Collection(database.CollectionUserSubscriptions).UpdateOne(ctx, bson.M{"_id": subscriptionId}, bson.M{"$set": bson.M{"expiryDate": expiryDate}})
	if err != nil {
		return err
	}
	return nil
}
