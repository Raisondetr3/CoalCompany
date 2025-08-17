package miner

import (
	"context"
	"sync"
	"time"
)

type SmallMiner struct {
	info   *MinerInfo
	сoalCh chan Coal
	mtx    sync.RWMutex
}

func NewSmallMiner() *SmallMiner {
	return &SmallMiner{
		info: &MinerInfo{
			id:           generateMinerID(),
			Type:         NormalMinerType,
			Cost:         SmallMinerCost,
			energy:       30,
			workSchedule: 1,
			BreakTime:    time.Second * 3,
			isActive:     false,
		},
		сoalCh: make(chan Coal, 30),
	}
}

func (m *SmallMiner) Info() MinerInfo {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	return *m.info
}

func (m *SmallMiner) Run(ctx context.Context) <-chan Coal {
	m.mtx.Lock()
	breakDur := m.info.BreakTime
	m.info.isActive = true
	m.mtx.Unlock()

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
