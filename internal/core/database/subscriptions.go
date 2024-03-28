package database

import (
	"errors"

	"github.com/cufee/aftermath-core/internal/core/database/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ErrSubscriptionNotFound = errors.New("subscription not found")

func GetSubscriptionByID(id primitive.ObjectID) (models.UserSubscription, error) {
	ctx, cancel := DefaultClient.Ctx()
	defer cancel()

	var subscription models.UserSubscription
	err := DefaultClient.Collection(CollectionUserSubscriptions).FindOne(ctx, bson.M{"_id": id}).Decode(&subscription)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return subscription, ErrSubscriptionNotFound
		}
		return subscription, err
	}

	return subscription, nil
}

func FindActiveSubscriptionsByUserID(userId string) ([]models.UserSubscription, error) {
	subscriptions, err := FindSubscriptionsByUserID(userId)
	if err != nil {
		return nil, err
	}

	var activeSubscriptions []models.UserSubscription
	for _, subscription := range subscriptions {
		if !subscription.IsExpired() {
			activeSubscriptions = append(activeSubscriptions, subscription)
		}
	}

	return activeSubscriptions, nil
}

func FindSubscriptionsByUserID(userId string) ([]models.UserSubscription, error) {
	ctx, cancel := DefaultClient.Ctx()
	defer cancel()

	var subscriptions []models.UserSubscription
	opts := options.Find().SetSort(bson.M{"creationDate": -1})
	cur, err := DefaultClient.Collection(CollectionUserSubscriptions).Find(ctx, bson.M{"userID": userId}, opts)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrSubscriptionNotFound
		}
		return nil, err
	}

	return subscriptions, cur.All(ctx, &subscriptions)
}

func FindSubscriptionsByReferenceIDs(referenceIDs ...string) ([]models.UserSubscription, error) {
	if len(referenceIDs) == 0 {
		return nil, nil
	}

	ctx, cancel := DefaultClient.Ctx()
	defer cancel()

	var subscriptions []models.UserSubscription
	opts := options.Find().SetSort(bson.M{"creationDate": 1})
	cur, err := DefaultClient.Collection(CollectionUserSubscriptions).Find(ctx, bson.M{"referenceID": bson.M{"$in": referenceIDs}}, opts)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrSubscriptionNotFound
		}
		return nil, err
	}

	return subscriptions, cur.All(ctx, &subscriptions)
}

func FindActiveSubscriptionsByReferenceIDs(referenceIDs ...string) ([]models.UserSubscription, error) {
	subscriptions, err := FindSubscriptionsByReferenceIDs(referenceIDs...)
	if err != nil {
		return nil, err
	}

	var activeSubscriptions []models.UserSubscription
	for _, subscription := range subscriptions {
		if !subscription.IsExpired() {
			activeSubscriptions = append(activeSubscriptions, subscription)
		}
	}

	return activeSubscriptions, nil
}

func AddNewUserSubscription(userId string, payload models.UserSubscription) (models.UserSubscription, error) {
	payload.ID = primitive.NilObjectID // Ensure ID is empty

	ctx, cancel := DefaultClient.Ctx()
	defer cancel()

	result, err := DefaultClient.Collection(CollectionUserSubscriptions).InsertOne(ctx, payload)
	if err != nil {
		return models.UserSubscription{}, err
	}

	id, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return models.UserSubscription{}, errors.New("invalid inserted id")
	}

	payload.ID = id
	return payload, nil
}

func UpdateUserSubscription(id primitive.ObjectID, payload models.SubscriptionUpdate) (models.UserSubscription, error) {
	ctx, cancel := DefaultClient.Ctx()
	defer cancel()

	_, err := DefaultClient.Collection(CollectionUserSubscriptions).UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": payload})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.UserSubscription{}, ErrSubscriptionNotFound
		}
		return models.UserSubscription{}, err
	}

	return GetSubscriptionByID(id)
}
