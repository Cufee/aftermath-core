package wotblitz

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/cufee/aftermath-core/internal/core/wargaming"
)

type RatingSeason struct {
	SeasonID  int    `json:"current_season"`
	StartAt   string `json:"start_at"`
	FinishAt  string `json:"finish_at"`
	UpdatedAt string `json:"updated_at"`

	Rewards []Reward `json:"rewards"`
	Leagues []League `json:"leagues"`

	TotalPlayers int `json:"count"`
}

type League struct {
	Title      string  `json:"title"`
	SmallIcon  string  `json:"small_icon"`
	BigIcon    string  `json:"big_icon"`
	Background string  `json:"background"`
	Index      int     `json:"index"`
	Percentile float64 `json:"percentile"`
}
type Reward struct {
	Type         string         `json:"type"`
	Vehicle      *VehicleReward `json:"vehicle"`
	Item         *ItemReward    `json:"stuff"`
	FromPosition int            `json:"from_position"`
	ToPosition   int            `json:"to_position"`
	Count        int            `json:"count"`
}

type ItemReward struct {
	Tag      string `json:"name"`
	Type     string `json:"type"`
	Name     int    `json:"title"`
	Count    int    `json:"count"`
	ImageURL string `json:"image_url"`
}

type VehicleReward struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	Nation          string `json:"nation"`
	Level           int    `json:"level"`
	UserString      string `json:"user_string"`
	ImageURL        string `json:"image_url"`
	PreviewImageURL string `json:"preview_image_url"`
	IsPremium       bool   `json:"is_premium"`
	IsCollectible   bool   `json:"is_collectible"`
}

func GetCurrentRatingSeason(realm string) (*RatingSeason, error) {
	res, err := client.Get(apiUrl(realm, "/rating-leaderboards/season/"))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status code: %d", res.StatusCode)
	}

	var season RatingSeason
	err = json.NewDecoder(res.Body).Decode(&season)
	if err != nil {
		return nil, err
	}

	return &season, nil
}

type PlayerLeaderboard struct {
	PlayerPosition
	Neighbors []PlayerPosition `json:"neighbors"`
}

type PlayerPosition struct {
	UpdatedAt string `json:"updated_at"`

	SeasonID  int    `json:"season_number"`
	AccountID int    `json:"spa_id"`
	Nickname  string `json:"nickname"`
	ClanTag   string `json:"clan_tag"`

	Score      int     `json:"score"`
	Position   int     `json:"number"`
	Percentile float64 `json:"percentile"`

	RawRating              float64 `json:"mmr"`
	LeagueIndex            int     `json:"league_index"`
	CalibrationBattlesLeft int     `json:"calibrationBattlesLeft"`

	Skip bool `json:"skip"`
}

func GetPlayerRatingPosition(accountID int, neighbors int) (*PlayerLeaderboard, error) {
	res, err := client.Get(apiUrl(wargaming.Clients.Live.RealmFromAccountID(strconv.Itoa(accountID)), fmt.Sprintf("/rating-leaderboards/user/%d/?neighbors=%d", accountID, neighbors)))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status code: %d", res.StatusCode)
	}

	var leaderboard PlayerLeaderboard
	err = json.NewDecoder(res.Body).Decode(&leaderboard)
	if err != nil {
		return nil, err
	}

	return &leaderboard, nil
}
