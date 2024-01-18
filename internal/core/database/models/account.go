package models

import "time"

type Account struct {
	ID       int    `json:"id" bson:"_id"`
	Realm    string `json:"realm" bson:"realm"`
	Nickname string `json:"nickname" bson:"nickname"`
	Private  bool   `json:"private" bson:"private"` // Some accounts have stats API disabled by Wargaming

	Clan *AccountClan `json:"clan" bson:"clan"`

	LastBattleTime time.Time `json:"lastBattleTime" bson:"lastBattleTime"` // This will probably end up not being updated too often

	LastUpdated time.Time `json:"lastUpdated" bson:"lastUpdated"`
}

type AccountClan struct {
	ID       int       `json:"id" bson:"_id"`
	Role     string    `json:"role" bson:"role"`
	JoinedAt time.Time `json:"joinedAt" bson:"joinedAt"`
}
