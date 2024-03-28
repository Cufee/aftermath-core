package models

import (
	"time"

	"github.com/cufee/aftermath-core/internal/core/stats"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RatingSnapshot struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	CreatedAt time.Time          `bson:"createdAt"`

	SeasonID       int `bson:"seasonId"`
	AccountID      int `bson:"accountId"`
	LastBattleTime int `bson:"lastBattleTime"`

	Stats stats.ReducedStatsFrame `bson:",inline"`
}
