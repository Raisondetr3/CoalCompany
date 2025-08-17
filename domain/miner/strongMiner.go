package miner

import (
	"context"
	"sync"
	"time"
)

type StrongMiner struct {
	info   *MinerInfo
	CoalCh chan Coal
	mtx    sync.RWMutex
}

func NewStrongMiner() *StrongMiner {
	return &StrongMiner{
		info: &MinerInfo{
			id:           generateMinerID(),
			Type:         StrongMinerType,
			Cost:         StrongMinerCost,
			energy:       60,
			workSchedule: 10,
			BreakTime:    time.Second,
			isActive:     false,
		},
		CoalCh: make(chan Coal, 60),
	}
}

func (m *StrongMiner) Info() MinerInfo {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	return *m.info
}

func (m *StrongMiner) Run(ctx context.Context) <-chan Coal {
	m.mtx.RLock()
	breakDur := m.info.BreakTime
	m.mtx.RUnlock()

	go func() {
		defer func() {
			close(m.CoalCh)
			m.mtx.Lock()
			m.info.isActive = false
			m.mtx.Unlock()
		}()

		ticker := time.NewTicker(breakDur)
		defer ticker.Stop()

		for {
			m.mtx.Lock()
			if m.info.energy <= 0 {
				m.mtx.Unlock()
				return
			}

			coal := Coal{Count: m.info.workSchedule}
			m.info.energy--
			m.info.workSchedule += 3
			m.mtx.Unlock()

			select {
			case <-ctx.Done():
				return
			case m.CoalCh <- coal:
			}

			select {
			case <-ticker.C:
			case <-ctx.Done():
				return
			}
		}
	}()

	return m.CoalCh
}
