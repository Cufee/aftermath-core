package models

import "time"

type AppConfiguration[T any] struct {
	Key       string         `json:"key" bson:"_id"`
	Value     T              `json:"value" bson:"value"`
	Metadata  map[string]any `json:"metadata" bson:"metadata"`
	UpdatedAt time.Time      `json:"updated_at" bson:"updated_at"`
}
