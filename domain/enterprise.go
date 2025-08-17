package domain

import (
	"CoalCompany/domain/miner"
	"CoalCompany/dto"
	appErrors "CoalCompany/errors"
	"context"
	"sync"
	"time"
)

type Enterprise struct {
	ctx         context.Context
	cancel      func()
	balance     int
	equipments  []Equipment
	hiredMiners []miner.Miner
	mtx         sync.RWMutex
	isGameOver  bool
}

func NewEnterprise(ctx context.Context, cancel func()) *Enterprise {
	return &Enterprise{
		ctx:         ctx,
		cancel:      cancel,
		balance:     0,
		equipments:  make([]Equipment, 0),
		hiredMiners: make([]miner.Miner, 0),
		mtx:         sync.RWMutex{},
	}
}

func (e *Enterprise) GetBalance() int {
	e.mtx.RLock()
	defer e.mtx.RUnlock()
	return e.balance
}

func (e *Enterprise) GetEquipments() []Equipment {
	e.mtx.RLock()
	defer e.mtx.RUnlock()

	equipments := make([]Equipment, len(e.equipments))
	copy(equipments, e.equipments)
	return equipments
}

func (e *Enterprise) GetHiredMinersCount() int {
	e.mtx.RLock()
	defer e.mtx.RUnlock()
	return len(e.hiredMiners)
}

func (e *Enterprise) StartPassiveIncome() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-e.ctx.Done():
			return
		case <-ticker.C:
			e.mtx.Lock()
			e.balance++
			e.mtx.Unlock()
		}
	}
}

func (e *Enterprise) HireMiner(minerType string) (*dto.HiredMinerInfo, error) {
	e.mtx.Lock()
	defer e.mtx.Unlock()

	cost, err := miner.GetMinerCostByType(minerType)
	if err != nil {
		return nil, appErrors.ErrMinerTypeNotFound
	}

	if e.balance < cost {
		return nil, appErrors.ErrInsufficientFunds
	}

	e.balance -= cost

	newMiner, err := miner.CreateMiner(minerType)
	if err != nil {
		return nil, appErrors.ErrMinerTypeNotFound
	}

	coalCh := newMiner.Run(e.ctx)
	go e.TakeCoal(coalCh)

	e.hiredMiners = append(e.hiredMiners, newMiner)

	return &dto.HiredMinerInfo{
		ID:        newMiner.Info().GetID(),
		Type:      minerType,
		CurEnergy: newMiner.Info().GetEnergy(),
	}, nil
}

func (e *Enterprise) TakeCoal(coalCh <-chan miner.Coal) {
	for {
		select {
		case <-e.ctx.Done():
			return
		case coal, ok := <-coalCh:
			if !ok {
				return
			}
			e.mtx.Lock()
			e.balance += coal.Count
			e.mtx.Unlock()
		}
	}
}

func (e *Enterprise) FindHiredMiners(isActive *bool, class string) dto.HiredMinersResponse {
	res := make([]dto.HiredMinerInfo, 0)
	e.mtx.RLock()
	defer e.mtx.RUnlock()

	for _, m := range e.hiredMiners {
		if isActive != nil && m.Info().GetIsActive() != *isActive {
			continue
		}

		if class != "" && string(m.Info().Type) != class {
			continue
		}

		res = append(res, dto.MapMinerToHiredMinerInfo(m))
	}

	return dto.HiredMinersResponse{
		Miners: res,
	}
}

func (e *Enterprise) GetAllEquipments() dto.EquipmentResponse {
	res := make([]dto.EquipmentInfo, 0, len(EquipmentCatalog))

	for _, eq := range EquipmentCatalog {
		res = append(res, dto.EquipmentInfo{
			Type:  string(eq.TypeEquipment),
			Price: eq.Price,
		})
	}

	return dto.EquipmentResponse{Equipments: res}
}

func (e *Enterprise) GetPurchasedEquipments() dto.EquipmentResponse {
	e.mtx.RLock()
	defer e.mtx.RUnlock()

	res := make([]dto.EquipmentInfo, 0, len(e.equipments))

	for _, eq := range e.equipments {
		res = append(res, dto.EquipmentInfo{
			Type:  string(eq.TypeEquipment),
			Price: eq.Price,
		})
	}

	return dto.EquipmentResponse{Equipments: res}
}

func (e *Enterprise) BuyEquipment(eqType string) error {
	equipmentType, err := ValidateEquipmentType(eqType)
	if err != nil {
		return err
	}

	e.mtx.Lock()
	defer e.mtx.Unlock()

	if e.isEquipmentPurchased(string(equipmentType)) {
		return appErrors.ErrEquipmentAlreadyPurchased
	}

	equipment := GetEquipmentByType(equipmentType)
	if equipment == nil {
		return appErrors.ErrEquipmentNotFound
	}

	if e.balance < equipment.Price {
		return appErrors.ErrInsufficientFunds
	}

	e.balance -= equipment.Price
	equipment.Purchased = true
	e.equipments = append(e.equipments, *equipment)

	return nil
}

func (e *Enterprise) isEquipmentPurchased(eqType string) bool {
	for _, eq := range e.equipments {
		if string(eq.TypeEquipment) == eqType {
			return true
		}
	}
	return false
}

func (e *Enterprise) GetEnterpriseStatsSafe() dto.EnterpriseStats {
	e.mtx.Lock()
	defer e.mtx.Unlock()

	return e.GetEnterpriseStatsTemplate()
}

func (e *Enterprise) GetEnterpriseStatsUnsafe() dto.EnterpriseStats {
	return e.GetEnterpriseStatsTemplate()
}

func (e *Enterprise) GetEnterpriseStatsTemplate() dto.EnterpriseStats {
	activeMiners := 0
	for _, miner := range e.hiredMiners {
		if miner.Info().GetIsActive() {
			activeMiners++
		}
	}

	purchasedEquipment := make([]dto.EquipmentInfo, 0, len(e.equipments))
	for _, eq := range e.equipments {
		purchasedEquipment = append(purchasedEquipment, dto.EquipmentInfo{
			Type:  string(eq.TypeEquipment),
			Price: eq.Price,
		})
	}

	totalEquipmentTypes := len(EquipmentCatalog)
	purchasedCount := len(e.equipments)

	return dto.EnterpriseStats{
		Balance:            e.balance,
		ActiveMiners:       activeMiners,
		TotalMinersHired:   len(e.hiredMiners),
		PurchasedEquipment: purchasedEquipment,
		GameProgress: dto.GameProgress{
			TotalEquipmentTypes:     totalEquipmentTypes,
			PurchasedEquipmentCount: purchasedCount,
			IsCompleted:             purchasedCount == totalEquipmentTypes,
		},
	}
}

func (e *Enterprise) ShutdownGame() (dto.GameShutdownResponse, error) {
	e.mtx.Lock()
	defer e.mtx.Unlock()

	if e.isGameOver {
		return dto.GameShutdownResponse{}, appErrors.ErrEnterpriseAlreadyShutdown
	}

	finalStats := e.GetEnterpriseStatsUnsafe()
	gameSummary := dto.GameSummary{
		IsGameOver:       e.isGameOver,
		FinalBalance:     e.balance,
		TotalMinersHired: len(e.hiredMiners),
	}

	e.cancel()
	e.isGameOver = true

	return dto.GameShutdownResponse{
		FinalStats:  finalStats,
		GameSummary: gameSummary,
	}, nil
}
