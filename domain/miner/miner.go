package miner

import (
	appErrors "CoalCompany/errors"
	"context"
	"math/rand"
	"time"
)

type Cost int
type MinerType string

const (
	SmallMinerCost  Cost = 5
	NormalMinerCost Cost = 50
	StrongMinerCost Cost = 450
)

const (
	SmallMinerType  MinerType = "small"
	NormalMinerType MinerType = "normal"
	StrongMinerType MinerType = "strong"
)

type MinerInfo struct {
	id           int
	Type         MinerType
	Cost         Cost
	energy       int
	workSchedule int
	BreakTime    time.Duration
	isActive     bool
}

func (m MinerInfo) GetID() int {
	return m.id
}

func (m MinerInfo) GetEnergy() int {
	return m.energy
}

func (m MinerInfo) GetIsActive() bool {
	return m.isActive
}

func (m MinerInfo) GetWorkSchedule() int {
	return m.workSchedule
}

type Miner interface {
	Run(ctx context.Context) <-chan Coal
	Info() MinerInfo
}

func GetMinerCostByType(minerType string) (int, error) {
	switch minerType {
	case "small":
		return int(SmallMinerCost), nil
	case "normal":
		return int(NormalMinerCost), nil
	case "strong":
		return int(StrongMinerCost), nil
	}

	return 0, appErrors.ErrMinerTypeNotFound
}

func CreateMiner(minerType string) (Miner, error) {
	switch minerType {
	case "small":
		return NewSmallMiner(), nil
	case "normal":
		return NewNormalMiner(), nil
	case "strong":
		return NewStrongMiner(), nil
	default:
		return nil, appErrors.ErrMinerTypeNotFound
	}
}

func generateMinerID() int {
	return rand.Intn(10000)
}
