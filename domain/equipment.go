package domain

import appErrors "CoalCompany/errors"

type TypeEquipment string

const (
	Pickaxes    TypeEquipment = "pickaxes"
	Ventilation TypeEquipment = "ventilation"
	Trolleys    TypeEquipment = "trolleys"
)

type Equipment struct {
	TypeEquipment TypeEquipment
	Price         int
	Purchased     bool
}

var EquipmentCatalog = []Equipment{
	{TypeEquipment: Pickaxes, Price: 3000, Purchased: false},
	{TypeEquipment: Ventilation, Price: 15000, Purchased: false},
	{TypeEquipment: Trolleys, Price: 50000, Purchased: false},
}

func GetEquipmentByType(eqType TypeEquipment) *Equipment {
	for _, eq := range EquipmentCatalog {
		if eq.TypeEquipment == eqType {
			eqCopy := eq
			return &eqCopy
		}
	}

	return nil
}

func ValidateEquipmentType(eqType string) (TypeEquipment, error) {
	switch eqType {
	case string(Pickaxes), string(Ventilation), string(Trolleys):
		return TypeEquipment(eqType), nil
	default:
		return "", appErrors.ErrEquipmentNotFound
	}

}
