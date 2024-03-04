package replay

import (
	"strconv"
	"time"

	"github.com/cufee/aftermath-core/internal/core/stats"
)

type battleType struct {
	ID  int    `json:"id"`
	Tag string `json:"tag"`
}

func (bt battleType) String() string {
	return bt.Tag
}

var (
	BattleTypeUnknown   = battleType{-1, "battle_type_unknown"}
	BattleTypeRandom    = battleType{0, "battle_type_regular"}
	BattleTypeSupremacy = battleType{1, "battle_type_supremacy"}
)

var battleTypes = map[int]battleType{
	0: BattleTypeRandom,
	1: BattleTypeSupremacy,
}

type gameMode struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Special bool   `json:"special"` // Signifies if WN8 should be calculated
}

func (gm gameMode) String() string {
	return gm.Name
}

var (
	GameModeUnknown         = gameMode{-1, "game_mode_unknown", false}
	GameModeRegular         = gameMode{1, "game_mode_regular", false}
	GameModeTraining        = gameMode{2, "game_mode_training", true}
	GameModeTournament      = gameMode{4, "game_mode_tournament", true}
	GameModeQuickTournament = gameMode{5, "game_mode_quick_tournament", true}
	GameModeRating          = gameMode{7, "game_mode_rating", false}
	GameModeMadGames        = gameMode{8, "game_mode_mad_games", true}
	GameModeRealistic       = gameMode{22, "game_mode_realistic", false}
	GameModeUprising        = gameMode{23, "game_mode_uprising", true}
	GameModeGravity         = gameMode{24, "game_mode_gravity", true}
	GameModeSkirmish        = gameMode{25, "game_mode_skirmish", false}
	GameModeBurningGames    = gameMode{26, "game_mode_burning_games", true}
)

var gameModes = map[int]gameMode{
	1:  GameModeRegular,
	2:  GameModeTraining,
	4:  GameModeTournament,
	5:  GameModeQuickTournament,
	7:  GameModeRating,
	8:  GameModeMadGames,
	22: GameModeRealistic,
	23: GameModeUprising,
	24: GameModeGravity,
	25: GameModeSkirmish,
	26: GameModeBurningGames,
}

type Replay struct {
	MapID      int        `json:"mapId"`
	GameMode   gameMode   `json:"gameMode"`
	BattleType battleType `json:"battleType"`

	Victory        bool      `json:"victory"`
	BattleTime     time.Time `json:"battleTime"`
	BattleDuration int       `json:"battleDuration"`

	Spoils      Spoils `json:"spoils"`
	Protagonist Player `json:"protagonist"`

	Teams Teams `json:"teams"`
}

func Prettify(battle battleResults, meta replayMeta) *Replay {
	var replay Replay

	replay.GameMode = GameModeUnknown
	if gm, ok := gameModes[int(battle.RoomType)]; ok {
		replay.GameMode = gm
	}

	// ModeAndMap
	replay.BattleType = BattleTypeUnknown
	if bt, ok := battleTypes[battle.GameMode()]; ok {
		replay.BattleType = bt
	}

	replay.MapID = battle.MapID()
	ts, _ := strconv.ParseInt(meta.BattleStartTime, 10, 64)
	replay.BattleTime = time.Unix(ts, 0)
	replay.BattleDuration = int(meta.BattleDuration)

	replay.Spoils = Spoils{
		Exp:     int(battle.Author.TotalXP),
		Credits: int(battle.Author.TotalCredits),
		// TODO: Find where mastery is set
		// MasteryBadge: data.MasteryBadge,
	}

	players := make(map[int]playerInfo)
	for _, p := range battle.Players {
		players[int(p.AccountID)] = p.Info
	}

	for _, result := range battle.PlayerResults {
		info, ok := players[int(result.Info.AccountID)]
		if !ok {
			continue
		}
		player := playerFromData(battle, info, result.Info)
		if player.ID == int(battle.Author.AccountID) {
			replay.Protagonist = player
		}
		if info.Team == TeamAlly {
			replay.Teams.Allies = append(replay.Teams.Allies, player)
		} else {
			replay.Teams.Enemies = append(replay.Teams.Enemies, player)
		}
	}

	return &replay
}

type Teams struct {
	Allies  []Player `json:"allies"`
	Enemies []Player `json:"enemies"`
}

type Player struct {
	ID       int    `json:"id"`
	ClanID   int    `json:"clanId"`
	ClanTag  string `json:"clanTag"`
	Nickname string `json:"nickname"`

	VehicleID int  `json:"vehicleId"`
	PlatoonID *int `json:"platoonId"`
	TimeAlive int  `json:"timeAlive"`
	HPLeft    int  `json:"hpLeft"`

	Performance  Performance `json:"performance"`
	Achievements map[int]int `json:"achievements"`
}

func playerFromData(battle battleResults, info playerInfo, result playerResultsInfo) Player {
	var player Player
	player.ID = int(result.AccountID)
	player.Nickname = info.Nickname
	player.VehicleID = int(result.TankID)
	if info.ClanTag != nil && info.ClanID != nil {
		player.ClanTag = *info.ClanTag
		player.ClanID = int(*info.ClanID)
	}

	if info.PlatoonID != nil {
		id := int(*info.PlatoonID)
		player.PlatoonID = &id
	}

	var frame stats.ReducedStatsFrame
	frame.Battles = 1
	if info.Team == battle.WinnerTeam {
		frame.BattlesWon = 1
	}

	if result.HitpointsLeft != nil {
		player.HPLeft = int(*result.HitpointsLeft)
		println(player.HPLeft)
	}
	if player.HPLeft > 0 {
		frame.BattlesSurvived = 1
	}

	frame.DamageDealt = int(result.DamageDealt)
	frame.DamageReceived = int(result.DamageReceived)
	frame.ShotsHit = int(result.ShotsHit)
	frame.ShotsFired = int(result.ShotsFired)
	frame.Frags = int(result.EnemiesDestroyed)
	frame.MaxFrags = frame.Frags
	frame.EnemiesSpotted = int(result.DamageAssisted)
	// TODO: Parse this from replays, it seems that those fields are only present when a battle was won by cap
	// frame.CapturePoints =
	// frame.DroppedCapturePoints =
	player.Performance = Performance{
		DamageBlocked:     int(result.DamageBlocked),
		DamageReceived:    int(result.DamageReceived),
		DamageAssisted:    int(result.DamageAssisted + result.DamageAssistedTrack),
		DistanceTraveled:  int(result.DistanceTraveled),
		ReducedStatsFrame: frame,
	}

	player.Achievements = make(map[int]int)
	for _, a := range append(result.Achievements, result.AchievementsOther...) {
		player.Achievements[int(a.Tag)] = int(a.Value)
	}

	return player
}

type Performance struct {
	DamageBlocked         int `json:"damageBlocked"`
	DamageReceived        int `json:"damageReceived"`
	DamageAssisted        int `json:"damageAssisted"`
	DistanceTraveled      int `json:"distanceTraveled"`
	SupremacyPointsEarned int `json:"supremacyPointsEarned"`
	SupremacyPointsStolen int `json:"supremacyPointsStolen"`

	stats.ReducedStatsFrame `json:",inline"`
}

type Spoils struct {
	Exp          int `json:"exp"`
	Credits      int `json:"credits"`
	MasteryBadge int `json:"masteryBadge"`
}
