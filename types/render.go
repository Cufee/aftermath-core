package types

import (
	"github.com/cufee/aftermath-core/internal/core/database/models"
)

type RenderRequestPayload struct {
	TankLimit int      `json:"tank_limit"`
	Presets   []string `json:"presets"`
	SortBy    string   `json:"sort_by"`
	// Days      int    `json:"days"`

	TypeStr string `json:"type"`
}

func (p RenderRequestPayload) Type() models.SessionType {
	return models.ParseSessionType(p.TypeStr)
}
