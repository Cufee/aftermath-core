package blitzstars

import (
	"encoding/json"
	"fmt"
	"net/http"

	wg "github.com/cufee/am-wg-proxy-next/v2/types"
)

type TopPlayersResponse struct {
	PlayerID   int                `json:"_id"`
	Career     ExtendedStatsFrame `json:"statistics"`
	Last30Days ExtendedStatsFrame `json:"period30d"`
	Last90Days ExtendedStatsFrame `json:"period90d"`
}

type ExtendedStatsFrame struct {
	All wg.StatsFrame `json:"all"`
}

func GetPlayerStats(accountId int) (*TopPlayersResponse, error) {
	res, err := insecureClient.Get(fmt.Sprintf("%s/top/player/%d", starsStatsApiURL, accountId))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status code: %d", res.StatusCode)
	}

	var stats TopPlayersResponse
	err = json.NewDecoder(res.Body).Decode(&stats)
	if err != nil {
		return nil, err
	}
	return &stats, err
}

type TankHistoryEntry struct {
	TankID          int           `json:"tank_id"`
	LastBattleTime  int           `json:"last_battle_time"`
	BattlesLifeTime int           `json:"battle_life_time"`
	MarkOfMastery   int           `json:"mark_of_mastery"`
	Stats           wg.StatsFrame `json:"all"`
}

func GetPlayerTankHistories(accountId int) (map[int][]TankHistoryEntry, error) {
	res, err := insecureClient.Get(fmt.Sprintf("%s/tankhistories/for/%d", starsStatsApiURL, accountId))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status code: %d", res.StatusCode)
	}

	var histories []TankHistoryEntry
	err = json.NewDecoder(res.Body).Decode(&histories)
	if err != nil {
		return nil, err
	}

	var historiesMap = make(map[int][]TankHistoryEntry)
	for _, entry := range histories {
		historiesMap[entry.TankID] = append(historiesMap[entry.TankID], entry)
	}

	return historiesMap, err
}
