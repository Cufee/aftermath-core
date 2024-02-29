package wotinspector

import (
	"time"

	"github.com/cufee/aftermath-core/internal/core/stats"
)

type battleType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (bt battleType) String() string {
	return bt.Name
}

var (
	BattleTypeUnknown   = battleType{-1, "Unknown"}
	BattleTypeRandom    = battleType{0, "Regular Battle"}
	BattleTypeSupremacy = battleType{1, "Supremacy"}
)

var battleTypes = map[int]battleType{
	0: BattleTypeRandom,
	1: BattleTypeSupremacy,
}

type gameMode struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Special bool   `json:"special"`
}

func (gm gameMode) String() string {
	return gm.Name
}

var (
	GameModeUnknown    = gameMode{-1, "Unknown", false}
	GameModeRegular    = gameMode{1, "Regular Battle", false}
	GameModeTraining   = gameMode{2, "Training Room", true}
	GameModeTournament = gameMode{4, "Tournament Battle", true}
	GameModeRating     = gameMode{7, "Rating Battle", false}
	GameModeMadGames   = gameMode{8, "Mad Games", true}
	GameModeRealistic  = gameMode{22, "Realistic Battle", true}
	GameModeUprising   = gameMode{23, "Uprising", true}
	GameModeGravity    = gameMode{24, "Gravity Force", true}
)

var gameModes = map[int]gameMode{
	1:  GameModeRegular,
	2:  GameModeTraining,
	4:  GameModeTournament,
	7:  GameModeRating,
	8:  GameModeMadGames,
	22: GameModeRealistic,
	23: GameModeUprising,
	24: GameModeGravity,
}

type Replay struct {
	ID string `json:"id"`

	MapID      int        `json:"map_id"`
	GameMode   gameMode   `json:"game_mode"`
	BattleType battleType `json:"battle_type"`

	Victory        bool      `json:"victory"`
	BattleTime     time.Time `json:"battle_time"`
	BattleDuration int       `json:"battle_duration"`

	Spoils      Spoils `json:"spoils"`
	Protagonist Player `json:"protagonist"`

	Teams Teams `json:"teams"`
}

func (data *replayData) Replay() *Replay {
	var replay Replay
	replay.ID = data.ID

	replay.GameMode = GameModeUnknown
	if gm, ok := gameModes[data.RoomType]; ok {
		replay.GameMode = gm
	}
	replay.BattleType = BattleTypeUnknown
	if bt, ok := battleTypes[data.BattleType]; ok {
		replay.BattleType = bt
	}

	replay.MapID = data.MapID
	replay.BattleTime = data.BattleStartTime
	replay.BattleDuration = int(data.BattleDuration)

	replay.Spoils = Spoils{
		Exp:          data.ExpTotal,
		Credits:      data.CreditsTotal,
		MasteryBadge: data.MasteryBadge,
	}

	for _, p := range data.PlayersData {
		player := p.Player(data.WinnerTeam)
		if p.Dbid == data.Protagonist {
			replay.Protagonist = player
		}

		if p.Team == data.ProtagonistTeam {
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
	ClanID   int    `json:"clan_id"`
	ClanTag  string `json:"clan_tag"`
	Nickname string `json:"nickname"`

	VehicleID  int `json:"vehicle_id"`
	SquadIndex int `json:"squad_index"`
	TimeAlive  int `json:"time_alive"`
	HPLeft     int `json:"hp_left"`

	Performance  Performance `json:"performance"`
	Achievements map[int]int `json:"achievements"`
}

func (data *replayPlayer) Player(winningTeam int) Player {
	var player Player
	player.ID = data.EntityID
	player.ClanID = data.Clanid
	player.ClanTag = data.ClanTag
	player.Nickname = data.Name

	if data.SquadIndex != nil {
		player.SquadIndex = *data.SquadIndex
	}
	player.VehicleID = data.VehicleDescr
	player.TimeAlive = data.TimeAlive
	player.HPLeft = data.HitpointsLeft

	var frame stats.ReducedStatsFrame
	frame.Battles = 1
	if data.Team == winningTeam {
		frame.BattlesWon = 1
	}
	if data.HitpointsLeft > 0 {
		frame.BattlesSurvived = 1
	}
	frame.DamageDealt = data.DamageMade
	frame.DamageReceived = data.DamageReceived
	frame.ShotsHit = data.ShotsHit
	frame.ShotsFired = data.ShotsMade
	frame.Frags = data.EnemiesDestroyed
	frame.MaxFrags = data.EnemiesDestroyed
	frame.EnemiesSpotted = data.EnemiesSpotted
	frame.CapturePoints = data.BaseCapturePoints
	frame.DroppedCapturePoints = data.BaseDefendPoints
	player.Performance = Performance{
		DamageBlocked:     data.DamageBlocked,
		DamageReceived:    data.DamageReceived,
		DamageAssisted:    data.DamageAssisted,
		DistanceTraveled:  data.DistanceTravelled,
		ReducedStatsFrame: frame,
	}

	player.Achievements = make(map[int]int)
	for _, a := range data.Achievements {
		player.Achievements[a.T] = a.V
	}

	return player

}

type Performance struct {
	DamageBlocked         int `json:"damage_blocked"`
	DamageReceived        int `json:"damage_received"`
	DamageAssisted        int `json:"damage_assisted"`
	DistanceTraveled      int `json:"distance_traveled"`
	SupremacyPointsEarned int `json:"supremacy_points_earned"`
	SupremacyPointsStolen int `json:"supremacy_points_stolen"`

	stats.ReducedStatsFrame `json:",inline"`
}

type Spoils struct {
	Exp          int `json:"exp"`
	Credits      int `json:"credits"`
	MasteryBadge int `json:"mastery_badge"`
}

type replayData struct {
	ID                     string         `json:"id"`
	MapID                  int            `json:"map_id"`
	BattleDuration         float64        `json:"battle_duration"`
	Title                  string         `json:"title"`
	PlayerName             string         `json:"player_name"`
	Protagonist            int            `json:"protagonist"`
	VehicleDescr           int            `json:"vehicle_descr"`
	MasteryBadge           int            `json:"mastery_badge"`
	ExpBase                int            `json:"exp_base"`
	EnemiesSpotted         int            `json:"enemies_spotted"`
	EnemiesDestroyed       int            `json:"enemies_destroyed"`
	DamageAssisted         int            `json:"damage_assisted"`
	DamageMade             int            `json:"damage_made"`
	DetailsURL             string         `json:"details_url"`
	DownloadURL            string         `json:"download_url"`
	ArenaUniqueID          string         `json:"arena_unique_id"`
	DownloadCount          int            `json:"download_count"`
	DataVersion            int            `json:"data_version"`
	Private                bool           `json:"private"`
	PrivateClan            bool           `json:"private_clan"`
	BattleStartTime        time.Time      `json:"battle_start_time"`
	UploadTime             time.Time      `json:"upload_time"`
	Allies                 []int          `json:"allies"`
	Enemies                []int          `json:"enemies"`
	ProtagonistClan        int            `json:"protagonist_clan"`
	ProtagonistTeam        int            `json:"protagonist_team"`
	BattleResult           int            `json:"battle_result"`
	CreditsBase            int            `json:"credits_base"`
	Tags                   []int          `json:"tags"`
	BattleType             int            `json:"battle_type"`
	RoomType               int            `json:"room_type"`
	LastAccessedTime       time.Time      `json:"last_accessed_time"`
	WinnerTeam             int            `json:"winner_team"`
	FinishReason           int            `json:"finish_reason"`
	PlayersData            []replayPlayer `json:"players_data"`
	ExpTotal               int            `json:"exp_total"`
	CreditsTotal           int            `json:"credits_total"`
	RepairCost             int            `json:"repair_cost"`
	ExpFree                int            `json:"exp_free"`
	ExpFreeBase            int            `json:"exp_free_base"`
	ExpPenalty             int            `json:"exp_penalty"`
	CreditsPenalty         int            `json:"credits_penalty"`
	CreditsContributionIn  int            `json:"credits_contribution_in"`
	CreditsContributionOut int            `json:"credits_contribution_out"`
	CamouflageID           int            `json:"camouflage_id"`
}

type replayPlayer struct {
	Team                int                 `json:"team"`
	Name                string              `json:"name"`
	EntityID            int                 `json:"entity_id"`
	Dbid                int                 `json:"dbid"`
	Clanid              int                 `json:"clanid"`
	ClanTag             string              `json:"clan_tag"`
	HitpointsLeft       int                 `json:"hitpoints_left"`
	Credits             int                 `json:"credits"`
	Exp                 int                 `json:"exp"`
	ShotsMade           int                 `json:"shots_made"`
	ShotsHit            int                 `json:"shots_hit"`
	ShotsSplash         int                 `json:"shots_splash"`
	ShotsPen            int                 `json:"shots_pen"`
	DamageMade          int                 `json:"damage_made"`
	DamageReceived      int                 `json:"damage_received"`
	DamageAssisted      int                 `json:"damage_assisted"`
	DamageAssistedTrack int                 `json:"damage_assisted_track"`
	HitsReceived        int                 `json:"hits_received"`
	HitsBounced         int                 `json:"hits_bounced"`
	HitsSplash          int                 `json:"hits_splash"`
	HitsPen             int                 `json:"hits_pen"`
	EnemiesSpotted      int                 `json:"enemies_spotted"`
	EnemiesDamaged      int                 `json:"enemies_damaged"`
	EnemiesDestroyed    int                 `json:"enemies_destroyed"`
	TimeAlive           int                 `json:"time_alive"`
	DistanceTravelled   int                 `json:"distance_travelled"`
	KilledBy            int                 `json:"killed_by"`
	BaseCapturePoints   int                 `json:"base_capture_points"`
	BaseDefendPoints    int                 `json:"base_defend_points"`
	ExpForDamage        int                 `json:"exp_for_damage"`
	ExpForAssist        int                 `json:"exp_for_assist"`
	ExpTeamBonus        int                 `json:"exp_team_bonus"`
	WpPointsEarned      int                 `json:"wp_points_earned"`
	WpPointsStolen      int                 `json:"wp_points_stolen"`
	HeroBonusCredits    int                 `json:"hero_bonus_credits"`
	HeroBonusExp        int                 `json:"hero_bonus_exp"`
	DeathReason         int                 `json:"death_reason"`
	Achievements        []replayAchievement `json:"achievements"`
	VehicleDescr        int                 `json:"vehicle_descr"`
	TurretID            int                 `json:"turret_id"`
	GunID               int                 `json:"gun_id"`
	ChassisID           int                 `json:"chassis_id"`
	SquadIndex          *int                `json:"squad_index"`
	DamageBlocked       int                 `json:"damage_blocked"`
}

type replayAchievement struct {
	T int `json:"t"`
	V int `json:"v"`
}
