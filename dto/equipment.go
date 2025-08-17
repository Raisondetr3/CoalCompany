package dto

type EquipmentInfo struct {
	Type  string `json:"type"`
	Price int    `json:"price"`
}

type EquipmentResponse struct {
	Equipments []EquipmentInfo `json:"equipment"`
}

type PurchaseEquipmentRequest struct {
	Type string `json:"type"`
}

type PurchaseEquipmentResponse struct {
	Type    string `json:"type"`
	Balance int    `json:"balance"`
}
