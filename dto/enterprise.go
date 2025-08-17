package dto

type EnterpriseStats struct {
	Balance            int             `json:"balance"`
	ActiveMiners       int             `json:"active_miners"`
	TotalMinersHired   int             `json:"total_miners_hired"`
	PurchasedEquipment []EquipmentInfo `json:"purchased_equipment"`
	GameProgress       GameProgress    `json:"game_progress"`
}

type GameProgress struct {
	TotalEquipmentTypes     int  `json:"total_equipment_types"`
	PurchasedEquipmentCount int  `json:"purchased_equipment_count"`
	IsCompleted             bool `json:"is_completed"`
}

type GameShutdownResponse struct {
	FinalStats  EnterpriseStats `json:"final_stats"`
	GameSummary GameSummary     `json:"game_summary"`
}

type GameSummary struct {
	IsGameOver       bool `json:"is_game_over"`
	FinalBalance     int  `json:"final_balance"`
	TotalMinersHired int  `json:"total_miners_hired"`
	TotalCoalMined   int  `json:"total_coal_mined"`
}
