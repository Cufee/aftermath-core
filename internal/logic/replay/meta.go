package replay

type replayMeta struct {
	Version               string  `json:"version"`
	Title                 string  `json:"title"`
	Dbid                  string  `json:"dbid"`
	PlayerName            string  `json:"playerName"`
	BattleStartTime       string  `json:"battleStartTime"`
	PlayerVehicleName     string  `json:"playerVehicleName"`
	MapName               string  `json:"mapName"`
	ArenaUniqueID         string  `json:"arenaUniqueId"`
	BattleDuration        float64 `json:"battleDuration"`
	VehicleCompDescriptor int     `json:"vehicleCompDescriptor"`
	CamouflageID          int     `json:"camouflageId"`
	MapID                 int     `json:"mapId"`
	ArenaBonusType        int     `json:"arenaBonusType"`
}
