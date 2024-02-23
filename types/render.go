package types

import (
	"github.com/cufee/aftermath-core/internal/core/database/models"
)

type SessionRequestPayload struct {
	ReferenceID *string `json:"referenceId"`

	LastBattleBefore *int `json:"last_battle_before"`
	LastBattleAfter  *int `json:"last_battle_after"`

	TankLimit int    `json:"tank_limit"`
	SortBy    string `json:"sort_by"`

	Presets []string `json:"presets"`
	TypeStr string   `json:"type"`
}

func (p SessionRequestPayload) Type() models.SessionType {
	return models.ParseSessionType(p.TypeStr)
}
