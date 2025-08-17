package errors

import "errors"

var (
	ErrInsufficientFunds         = errors.New("insufficient funds")
	ErrMinerTypeNotFound         = errors.New("miner type not found")
	ErrEquipmentNotFound         = errors.New("equipment not found")
	ErrEquipmentAlreadyPurchased = errors.New("equipment already purchased")
	ErrGameNotCompleted          = errors.New("game not completed")
	ErrEnterpriseAlreadyShutdown = errors.New("enterprise already shutdown")
)
