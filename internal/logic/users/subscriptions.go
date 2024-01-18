package users

import (
	"errors"
	"time"

	"github.com/cufee/aftermath-core/internal/core/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type subscriptionType string

const (
	SubscriptionTypePro          = subscriptionType("aftermath-pro")
	SubscriptionTypeProClan      = subscriptionType("aftermath-pro-clsn")
	SubscriptionTypePlus         = subscriptionType("aftermath-plus")
	SubscriptionTypeSupporter    = subscriptionType("supporter")
	SubscriptionTypeVerifiedClan = subscriptionType("verified-clan")
)

var ErrSubscriptionNotFound = errors.New("subscription not found")

type UserSubscription struct {
	id          primitive.ObjectID `bson:"_id,omitempty" json:"-"`
	UserID      string             `bson:"userID" json:"userID"`
	ReferenceID *string            `bson:"referenceID" json:"referenceID"`

	Type         subscriptionType `bson:"subscriptionType" json:"subscriptionType"`
	ExpiryDate   time.Time        `bson:"expiryDate" json:"expiryDate"`
	CreationDate time.Time        `bson:"creationDate" json:"creationDate"`
}

func (s *UserSubscription) ID() primitive.ObjectID {
	return s.id
}

func (s *UserSubscription) IsExpired() bool {
	return s.ExpiryDate.Before(time.Now())
}

func FindActiveSubscriptionsByUserID(userId string) ([]UserSubscription, error) {
	subscriptions, err := FindSubscriptionsByUserID(userId)
	if err != nil {
		return nil, err
	}

	var activeSubscriptions []UserSubscription
	for _, subscription := range subscriptions {
		if !subscription.IsExpired() {
			activeSubscriptions = append(activeSubscriptions, subscription)
		}
	}

	return activeSubscriptions, nil
}

func FindSubscriptionsByUserID(userId string) ([]UserSubscription, error) {
	ctx, cancel := database.DefaultClient.Ctx()
	defer cancel()

	var subscriptions []UserSubscription
	opts := options.Find().SetSort(bson.M{"creationDate": -1})
	cur, err := database.DefaultClient.Collection(database.CollectionUserSubscriptions).Find(ctx, bson.M{"userID": userId}, opts)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrSubscriptionNotFound
		}
		return nil, err
	}

	return subscriptions, cur.All(ctx, &subscriptions)
}

func FindSubscriptionsByReferenceIDs(referenceIDs ...string) ([]UserSubscription, error) {
	if len(referenceIDs) == 0 {
		return nil, nil
	}

	ctx, cancel := database.DefaultClient.Ctx()
	defer cancel()

	var subscriptions []UserSubscription
	opts := options.Find().SetSort(bson.M{"creationDate": -1})
	cur, err := database.DefaultClient.Collection(database.CollectionUserSubscriptions).Find(ctx, bson.M{"referenceID": bson.M{"$in": referenceIDs}}, opts)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrSubscriptionNotFound
		}
		return nil, err
	}

	return subscriptions, cur.All(ctx, &subscriptions)
}

func FindActiveSubscriptionsByReferenceIDs(referenceIDs ...string) ([]UserSubscription, error) {
	subscriptions, err := FindSubscriptionsByReferenceIDs(referenceIDs...)
	if err != nil {
		return nil, err
	}

	var activeSubscriptions []UserSubscription
	for _, subscription := range subscriptions {
		if !subscription.IsExpired() {
			activeSubscriptions = append(activeSubscriptions, subscription)
		}
	}

	return activeSubscriptions, nil
}

func AddNewUserSubscription(userId string, payload UserSubscription) error {
	ctx, cancel := database.DefaultClient.Ctx()
	defer cancel()

	_, err := database.DefaultClient.Collection(database.CollectionUserSubscriptions).InsertOne(ctx, payload)
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
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrSubscriptionNotFound
		}
		return err
	}
	return nil
}
