package types

type RenderRequestPayload struct {
	TankLimit int    `json:"tank_limit"`
	SortBy    string `json:"sort_by"`
}
