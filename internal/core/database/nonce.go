package database

import (
	"errors"
	"time"

	"github.com/cufee/aftermath-core/internal/core/database/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrNonceNotFound = errors.New("nonce not found")
)

func NewNonce(referenceID string, duration time.Duration) (string, error) {
	ctx, cancel := DefaultClient.Ctx()
	defer cancel()

	var nonce models.Nonce
	nonce.ReferenceID = referenceID
	nonce.ExpiresAt = time.Now().Add(duration)
	nonce.CreatedAt = time.Now()

	res, err := DefaultClient.Collection(CollectionNonce).InsertOne(ctx, nonce)
	if err != nil {
		return "", err
	}

	id, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", errors.New("invalid inserted id")
	}

	return id.Hex(), nil
}

func GetNonceByID(id string) (*models.Nonce, error) {
	ctx, cancel := DefaultClient.Ctx()
	defer cancel()

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var nonce models.Nonce
	err = DefaultClient.Collection(CollectionNonce).FindOne(ctx, bson.M{"_id": oid, "expiresAt": bson.M{"$gt": time.Now()}}).Decode(&nonce)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNonceNotFound
		}
		return nil, err
	}

	return &nonce, nil
}

func ExpireNonceByID(id string) error {
	ctx, cancel := DefaultClient.Ctx()
	defer cancel()

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = DefaultClient.Collection(CollectionNonce).UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"$set": bson.M{"expiresAt": time.Now()}})
	if err != nil {
		return err
	}

	return nil
}
