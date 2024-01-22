package database

import (
	"errors"
	"time"

	"github.com/cufee/aftermath-core/internal/core/database/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
