package dto

import (
	"CoalCompany/domain/miner"
)

type MinerTypeInfo struct {
	Type         string `json:"type"`
	Cost         int    `json:"remuneration"`
	Energy       int    `json:"energy"`
	WorkSchedule int    `json:"work_schedule"`
	BreakSeconds int    `json:"break_seconds"`
}

type MinerTypesResponse struct {
	Types []MinerTypeInfo `json:"types"`
}

type MinerType struct {
	Type string `json:"type"`
}

type HiredMinerInfo struct {
	ID        int    `json:"id"`
	Type      string `json:"type"`
	CurEnergy int    `json:"cur_energy"`
}

type HiredMinersResponse struct {
	Miners []HiredMinerInfo `json:"types"`
}

func MapMinerToHiredMinerInfo(m miner.Miner) HiredMinerInfo {
	info := m.Info()
	var minerType string
	switch m.(type) {
	case *miner.SmallMiner:
		minerType = "small"
	case *miner.NormalMiner:
		minerType = "normal"
	case *miner.StrongMiner:
		minerType = "strong"
	}

	return HiredMinerInfo{
		ID:        info.GetID(),
		Type:      minerType,
		CurEnergy: info.GetEnergy(),
	}
}

func MapMinerToTypeInfo(minerType string, m miner.Miner) MinerTypeInfo {
	info := m.Info()
	return MinerTypeInfo{
		Type:         minerType,
		Cost:         int(info.Cost),
		Energy:       info.GetEnergy(),
		WorkSchedule: info.GetWorkSchedule(),
		BreakSeconds: int(info.BreakTime.Seconds()),
	}
}
