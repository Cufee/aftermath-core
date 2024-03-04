package replay

import (
	"io"

	"github.com/nlpodyssey/gopickle/pickle"
	"github.com/nlpodyssey/gopickle/types"
	"go.dedis.ch/protobuf"
)

type battleResults struct {
	ModeAndMap uint32 `protobuf:"1" json:"modeAndMap"`
	Timestamp  uint64 `protobuf:"2" json:"timestamp"`
	WinnerTeam uint32 `protobuf:"3" json:"winnerTeam"`

	EnemiesDestroyed uint32 `protobuf:"4" json:"enemiesKilled"`
	TimeAlive        uint32 `protobuf:"5" json:"timeAlive"`

	Author   author `protobuf:"8,required" json:"protagonist"`
	RoomType uint32 `protobuf:"9" json:"roomType"`

	RepairCost uint32 `protobuf:"136" json:"repairCost"`
	FreeXP     uint32 `protobuf:"137" json:"freeXp"`

	TotalXP     uint32 `protobuf:"181" json:"totalXp"`
	BaseFreeXP  uint32 `protobuf:"182" json:"freeXpBase"`
	CreditsBase uint32 `protobuf:"183" json:"creditsBase"`

	// 185 {
	// 	1: "None"
	// 	2 {
	// 		27 {
	// 			1: 2020
	// 			2: 63
	// 		}
	// 	}
	// }
	// 185 {
	// 	1: "booster"
	// 	2 {
	// 		1: 21959
	// 		2: 10
	// 	}
	// }

	Players       []player        `protobuf:"201,repeated" json:"players"`
	PlayerResults []playerResults `protobuf:"301,repeated" json:"playerResults"`
}

func (result *battleResults) MapID() int {
	return int(result.ModeAndMap & 0xFFFF)
}

func (result *battleResults) GameMode() int {
	/// - `0x0xxxx` is «encounter»
	/// - `0x1xxxx` flag is «supremacy»
	return int((result.ModeAndMap >> 16) & 0xFF)
}

type player struct {
	AccountID uint32     `protobuf:"1" json:"accountId"`
	Info      playerInfo `protobuf:"2,required" json:"info"`
}

type playerInfo struct {
	Nickname  string  `protobuf:"1" json:"nickname"`
	PlatoonID *uint32 `protobuf:"2,optional" json:"platoon"`
	Team      uint32  `protobuf:"3" json:"team"`
	ClanID    *uint32 `protobuf:"4,optional" json:"clanId"`
	ClanTag   *string `protobuf:"5,optional" json:"clanTag"`
	// 6: "\000\000"
	Avatar avatar `protobuf:"7,required" json:"avatar"`
	// 8: 463102
	Rank *uint32 `protobuf:"9,optional" json:"rank"`
}

type playerResults struct {
	ResultID uint32            `protobuf:"1" json:"id"`
	Info     playerResultsInfo `protobuf:"2,required" json:"result"`
}

type playerResultsInfo struct {
	HitpointsLeft       *uint64 `protobuf:"1,optional" json:"hitpointsLeft"`
	CreditsEarned       uint32  `protobuf:"2" json:"creditsEarned"`
	BaseXP              uint32  `protobuf:"3" json:"baseXp"`
	ShotsFired          uint32  `protobuf:"4" json:"shotsFired"`
	ShotsHit            uint32  `protobuf:"5" json:"shotsHit"`
	ShotsPenetrated     uint32  `protobuf:"7" json:"shotsPenetrated"`
	DamageDealt         uint32  `protobuf:"8" json:"damageDealt"`
	DamageAssisted      uint32  `protobuf:"9" json:"damageAssisted"`
	DamageAssistedTrack uint32  `protobuf:"10" json:"damageAssistedTrack"`
	DamageReceived      uint32  `protobuf:"11" json:"damageReceived"`
	HitsReceived        uint32  `protobuf:"12" json:"hitsReceived"`
	HitsBlocked         uint32  `protobuf:"13" json:"hitsBlocked"`
	HitsPenetrated      uint32  `protobuf:"15" json:"hitsPenetrated"`
	EnemiesSpotted      uint32  `protobuf:"16" json:"enemiesSpotted"`
	EnemiesDamaged      uint32  `protobuf:"17" json:"enemiesDamaged"`
	EnemiesDestroyed    uint32  `protobuf:"18" json:"enemiesDestroyed"`
	DistanceTraveled    uint32  `protobuf:"23" json:"distanceTraveled"`
	TimeAlive           uint32  `protobuf:"24" json:"timeAlive"`
	KilledByAccountID   uint32  `protobuf:"25" json:"killedBy"`
	//  26: "\323\003"
	Achievements          []Achievement `protobuf:"27,repeated" json:"achievements"`
	AchievementsOther     []Achievement `protobuf:"28,repeated" json:"achievementsOther"`
	DamageXP              uint32        `protobuf:"29" json:"damageXp"`
	AssistXP              uint32        `protobuf:"30" json:"assistXp"`
	TeamBonusXP           uint32        `protobuf:"31" json:"teamBonusXp"`
	SupremacyPointsEarned uint32        `protobuf:"32" json:"supremacyPointsEarned"`
	SupremacyPointsStolen uint32        `protobuf:"33" json:"supremacyPointsStolen"`

	AccountID uint32 `protobuf:"101" json:"accountId"`
	Team      uint32 `protobuf:"102" json:"team"`
	TankID    uint32 `protobuf:"103" json:"tankId"`
	//  105: 18446744073709551615
	// 	106: 80120
	MMRating      float32 `protobuf:"107" json:"mmRating"`
	DamageBlocked uint32  `protobuf:"117" json:"damageBlocked"`
}

func (results *playerResultsInfo) DisplayRating() float32 {
	return results.MMRating*10 + 3000
}

type Achievement struct {
	Tag   uint32 `protobuf:"1" json:"tag"`
	Value uint32 `protobuf:"2" json:"value"`
}

const PlayerAfkHitpointsLeft = -2

type author struct {
	// The field is PlayerAfkHitpointsLeft if the player is auto-destroyed by inactivity.
	HitpointsLeft   int32  `protobuf:"1" json:"hitpointsLeft"`
	TotalCredits    uint32 `protobuf:"2" json:"creditsEarned"`
	TotalXP         uint32 `protobuf:"3" json:"xpEarned"`
	ShotsFired      uint32 `protobuf:"4" json:"shotsFired"`
	ShotsHit        uint32 `protobuf:"5" json:"shotsHit"`
	ShotsSplashHit  uint32 `protobuf:"6" json:"shotsSplashHit"`
	ShotsPenetrated uint32 `protobuf:"7" json:"shotsPenetrated"`
	DamageDealt     uint32 `protobuf:"8" json:"damageDealt"`

	AccountID uint32 `protobuf:"101" json:"accountId"`
	Team      uint32 `protobuf:"102" json:"team"`
}

type avatar struct {
	// 1: 312062 id?
	Info avatarInfo `protobuf:"2,required" json:"data"`
}

type avatarInfo struct {
	// 1: 312062 id?
	GfxURL  string `protobuf:"2" json:"gfxUrl"`
	Gfx2URL string `protobuf:"3" json:"gfx2Url"`
	Kind    string `protobuf:"4" json:"kind"`
}

// Un-pickled `battle_results.dat`.
//
// # `battle_results.dat` structure
//
// Entire file is a [pickled](https://docs.python.org/3/library/pickle.html) 2-tuple:
//
// - Arena unique ID
// - Battle results serialized with [Protocol Buffers](https://developers.google.com/protocol-buffers)
func decodeBattleResults(reader io.Reader) (*battleResults, error) {
	unpickler := pickle.NewUnpickler(reader)
	data, err := unpickler.Load()
	if err != nil {
		return nil, err
	}

	if tuple, ok := data.(*types.Tuple); ok && tuple.Len() == 2 {
		if data, ok := tuple.Get(1).(string); ok {
			var result battleResults
			return &result, protobuf.Decode([]byte(data), &result)
		}
		return nil, ErrInvalidReplayFile
	}
	return nil, ErrInvalidReplayFile
}
