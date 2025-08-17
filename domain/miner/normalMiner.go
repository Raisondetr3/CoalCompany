package miner

import (
	"context"
	"sync"
	"time"
)

type NormalMiner struct {
	info   *MinerInfo
	mtx    sync.RWMutex
	сoalCh chan Coal
}

func NewNormalMiner() *NormalMiner {
	return &NormalMiner{
		info: &MinerInfo{
			id:           generateMinerID(),
			Type:         SmallMinerType,
			Cost:         NormalMinerCost,
			energy:       45,
			workSchedule: 3,
			BreakTime:    time.Second * 2,
			isActive:     false,
		},
		сoalCh: make(chan Coal, 45),
	}
}

func (m *NormalMiner) Info() MinerInfo {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	return *m.info
}

func (m *NormalMiner) Run(ctx context.Context) <-chan Coal {
	m.mtx.RLock()
	breakDur := m.info.BreakTime
	m.info.isActive = true
	m.mtx.RUnlock()

	go func() {
		defer func() {
			close(m.сoalCh)
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

			// Сразу начинает работать
			coal := Coal{Count: m.info.workSchedule}
			m.info.energy--
			m.mtx.Unlock()

			// Отправка угля
			select {
			case <-ctx.Done():
				return
			case m.сoalCh <- coal:
			}

			// Перерыв
			select {
			case <-ticker.C:
			case <-ctx.Done():
				return
			}
		}
	}()

	return m.сoalCh
}
