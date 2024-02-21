package types

import (
	"time"

	"github.com/cufee/aftermath-core/internal/core/database/models"
)

type RenderRequestPayload struct {
	ReferenceID *string `json:"referenceId"`

	BattlesAfter *time.Time `json:"battles_after"`
	TankLimit    int        `json:"tank_limit"`
	SortBy       string     `json:"sort_by"`

	Presets []string `json:"presets"`
	TypeStr string   `json:"type"`
}

func (p RenderRequestPayload) Type() models.SessionType {
	return models.ParseSessionType(p.TypeStr)
}
