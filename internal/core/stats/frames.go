package stats

import (
	"math"
)

const InvalidValue = -1

type ReducedStatsFrame struct {
	Battles         int `json:"battles" bson:"battles"`
	BattlesWon      int `json:"battlesWon" bson:"battlesWon"`
	BattlesSurvived int `json:"battlesSurvived" bson:"battlesSurvived"`
	DamageDealt     int `json:"damageDealt" bson:"damageDealt"`
	DamageReceived  int `json:"damageReceived" bson:"damageReceived"`

	ShotsHit   int `json:"shotsHit" bson:"shotsHit"`
	ShotsFired int `json:"shotsFired" bson:"shotsFired"`

	Frags                int `json:"frags" bson:"frags"`
	MaxFrags             int `json:"maxFrags" bson:"maxFrags"`
	EnemiesSpotted       int `json:"enemiesSpotted" bson:"enemiesSpotted"`
	CapturePoints        int `json:"capturePoints" bson:"capturePoints"`
	DroppedCapturePoints int `json:"droppedCapturePoints" bson:"droppedCapturePoints"`

	wn8         int     `json:"-" bson:"-"`
	winrate     float64 `json:"-" bson:"-"`
	accuracy    float64 `json:"-" bson:"-"`
	avgDamage   float64 `json:"-" bson:"-"` // float64 will be rendered as a percentage with 2 decimal places
	damageRatio float32 `json:"-" bson:"-"` // float32 will be rendered with 2 decimal places
}

func (r *ReducedStatsFrame) Valid(value interface{}) bool {
	return ValueValid(value)
}

func ValueValid(value interface{}) bool {
	switch cast := value.(type) {
	case int:
		return cast != InvalidValue
	case float64:
		return int(cast) != InvalidValue
	case float32:
		return int(cast) != InvalidValue
	default:
		return false
	}
}

func (r *ReducedStatsFrame) AvgDamage() float64 {
	if r.Battles == 0 {
		return InvalidValue
	}
	if r.avgDamage == 0 {
		r.avgDamage = float64(r.DamageDealt) / float64(r.Battles)
	}
	return r.avgDamage
}

func (r *ReducedStatsFrame) DamageRatio() float32 {
	if r.Battles == 0 || r.DamageReceived == 0 {
		return InvalidValue
	}
	if r.damageRatio == 0 {
		r.damageRatio = float32(r.DamageDealt) / float32(r.DamageReceived)
	}
	return r.damageRatio
}

func (r *ReducedStatsFrame) Winrate() float64 {
	if r.Battles == 0 {
		return InvalidValue
	}
	if r.winrate == 0 {
		r.winrate = float64(r.BattlesWon) / float64(r.Battles) * 100
	}
	return r.winrate
}

func (r *ReducedStatsFrame) Accuracy() float64 {
	if r.Battles == 0 || r.ShotsFired == 0 {
		return InvalidValue
	}
	if r.accuracy == 0 {
		r.accuracy = float64(r.ShotsHit) / float64(r.ShotsFired) * 100
	}
	return r.accuracy
}

/*
	 Calculate WN8 Rating for a tank using the following formula:
		(980*rDAMAGEc + 210*rDAMAGEc*rFRAGc + 155*rFRAGc*rSPOTc + 75*rDEFc*rFRAGc + 145*MIN(1.8,rWINc))/EXPc
*/
func (r *ReducedStatsFrame) WN8(average *ReducedStatsFrame) int {
	if r.wn8 > 0 {
		return r.wn8
	}
	if average == nil || r.Battles == 0 || average.Battles == 0 {
		return InvalidValue
	}

	// Expected values for WN8
	expDef := float64(average.DroppedCapturePoints) / float64(average.Battles)
	expFrag := float64(average.Frags) / float64(average.Battles)
	expSpot := float64(average.EnemiesSpotted) / float64(average.Battles)
	expDmg := average.AvgDamage()
	expWr := average.Winrate() / 100

	// Actual performance
	pDef := float64(r.DroppedCapturePoints) / float64(r.Battles)
	pFrag := float64(r.Frags) / float64(r.Battles)
	pSpot := float64(r.EnemiesSpotted) / float64(r.Battles)
	pDmg := r.AvgDamage()
	pWr := r.Winrate() / 100

	// Calculate WN8 metrics
	rDef := pDef / expDef
	rFrag := pFrag / expFrag
	rSpot := pSpot / expSpot
	rDmg := pDmg / expDmg
	rWr := pWr / expWr

	adjustedWr := math.Max(0, ((rWr - 0.71) / (1 - 0.71)))
	adjustedDmg := math.Max(0, ((rDmg - 0.22) / (1 - 0.22)))
	adjustedDef := math.Max(0, (math.Min(adjustedDmg+0.1, (rDef-0.10)/(1-0.10))))
	adjustedSpot := math.Max(0, (math.Min(adjustedDmg+0.1, (rSpot-0.38)/(1-0.38))))
	adjustedFrag := math.Max(0, (math.Min(adjustedDmg+0.2, (rFrag-0.12)/(1-0.12))))

	r.wn8 = int(math.Round(((980 * adjustedDmg) + (210 * adjustedDmg * adjustedFrag) + (155 * adjustedFrag * adjustedSpot) + (75 * adjustedDef * adjustedFrag) + (145 * math.Min(1.8, adjustedWr)))))

	return r.wn8
}

func (r *ReducedStatsFrame) Add(other *ReducedStatsFrame) {
	r.Battles += other.Battles
	r.BattlesWon += other.BattlesWon
	r.BattlesSurvived += other.BattlesSurvived
	r.DamageDealt += other.DamageDealt
	r.DamageReceived += other.DamageReceived

	r.ShotsHit += other.ShotsHit
	r.ShotsFired += other.ShotsFired

	r.Frags += other.Frags
	if r.MaxFrags < other.MaxFrags {
		r.MaxFrags = other.MaxFrags
	}

	r.EnemiesSpotted += other.EnemiesSpotted
	r.CapturePoints += other.CapturePoints
	r.DroppedCapturePoints += other.DroppedCapturePoints
}

func (r *ReducedStatsFrame) Subtract(other *ReducedStatsFrame) {
	r.Battles -= other.Battles
	r.BattlesWon -= other.BattlesWon
	r.BattlesSurvived -= other.BattlesSurvived
	r.DamageDealt -= other.DamageDealt
	r.DamageReceived -= other.DamageReceived

	r.ShotsHit -= other.ShotsHit
	r.ShotsFired -= other.ShotsFired

	r.Frags -= other.Frags
	if r.MaxFrags > other.MaxFrags {
		r.MaxFrags = other.MaxFrags
	}

	r.EnemiesSpotted -= other.EnemiesSpotted
	r.CapturePoints -= other.CapturePoints
	r.DroppedCapturePoints -= other.DroppedCapturePoints
}

type ReducedVehicleStats struct {
	VehicleID          int `json:"vehicleId" bson:"vehicleId"`
	*ReducedStatsFrame `bson:",inline"`

	MarkOfMastery  int `json:"markOfMastery" bson:"markOfMastery"`
	LastBattleTime int `json:"lastBattleTime" bson:"lastBattleTime"`
}

func (r *ReducedVehicleStats) Add(other *ReducedVehicleStats) {
	r.ReducedStatsFrame.Add(other.ReducedStatsFrame)
	if other.MarkOfMastery > r.MarkOfMastery {
		r.MarkOfMastery = other.MarkOfMastery
	}
	if other.LastBattleTime > r.LastBattleTime {
		r.LastBattleTime = other.LastBattleTime
	}
}

func (r *ReducedVehicleStats) Subtract(other *ReducedVehicleStats) {
	r.ReducedStatsFrame.Subtract(other.ReducedStatsFrame)
	if other.MarkOfMastery > r.MarkOfMastery {
		r.MarkOfMastery = other.MarkOfMastery
	}
	if other.LastBattleTime > r.LastBattleTime {
		r.LastBattleTime = other.LastBattleTime
	}
}
