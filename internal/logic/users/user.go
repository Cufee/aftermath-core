package users

import (
	"errors"

	"github.com/cufee/aftermath-core/internal/core/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type User struct {
	ID string `bson:"_id" json:"id"`

	FeatureFlags []featureFlag `bson:"featureFlags" json:"featureFlags"`
}

func newUser(id string) User {
	return User{
		ID:           id,
		FeatureFlags: []featureFlag{},
	}
}

type featureFlag string

const (
	FeatureFlagCustomizationDisabled = featureFlag("customizationDisabled")
)

func (u *User) HasFeatureFlag(flag featureFlag) bool {
	for _, f := range u.FeatureFlags {
		if f == flag {
			return true
		}
	}
	return false
}

func FindUserByID(id string) (*User, error) {
	ctx, cancel := database.DefaultClient.Ctx()
	defer cancel()

	var user User
	err := database.DefaultClient.Collection(database.CollectionUsers).FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func FindUserByConnection(connectionType connectionType, externalID string) (*User, error) {
	ctx, cancel := database.DefaultClient.Ctx()
	defer cancel()

	var connection UserConnection
	err := database.DefaultClient.Collection(database.CollectionUsers).FindOne(ctx, bson.M{"connectionType": connectionType, "connectionID": externalID}).Decode(&connection)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return FindUserByID(connection.UserID)
}

func CreateUser(id string) (*User, error) {
	ctx, cancel := database.DefaultClient.Ctx()
	defer cancel()

	user := newUser(id)

	_, err := database.DefaultClient.Collection(database.CollectionUsers).InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func UpdateUser(id string, update User) (*User, error) {
	update.ID = id

	ctx, cancel := database.DefaultClient.Ctx()
	defer cancel()

	_, err := database.DefaultClient.Collection(database.CollectionUsers).UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": update})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &update, nil
}
