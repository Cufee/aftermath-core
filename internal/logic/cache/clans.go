package cache

import "time"

type DatabaseClan struct {
	ID       int    `json:"id" bson:"_id"`
	Tag      string `json:"tag" bson:"tag"`
	Name     string `json:"name" bson:"name"`
	EmblemID string `json:"emblem" bson:"emblem"`

	Members   []int     `json:"members" bson:"members"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`

	LastUpdated time.Time `json:"lastUpdated" bson:"lastUpdated"`
}
