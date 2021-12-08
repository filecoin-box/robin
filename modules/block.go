package modules

import (
	"context"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/lotus/chain/types"
	logging "github.com/ipfs/go-log/v2"
	"time"
)

var blockLog = logging.Logger("block")

func (m *Monitor) getLatestBlock() (*types.TipSet, error) {
	tipSet, err := m.api.ChainHead(context.Background())
	if err != nil {
		blockLog.Errorw("getLatestBlock:ChainHead", "err", err)
		return nil, err
	}

	return tipSet, nil
}

func (m *Monitor) getTipSetByHeight(height abi.ChainEpoch) (*types.TipSet, error) {
	tipSet, err := m.api.ChainGetTipSetByHeight(context.Background(), height, types.NewTipSetKey())
	if err != nil {
		blockLog.Warnw("getTipSetByHeight:ChainGetTipSetByHeight", "err", err)
		return nil, err
	}

	if tipSet.Height() != height {
		return nil, nil
	}

	return tipSet, nil
}

func (m *Monitor) blockMonitor() {
	var latestHeight abi.ChainEpoch
	for {
		tipSet, err := m.getLatestBlock()
		if err != nil {
			blockLog.Errorw("getLatestBlock", "err", err)
			wait(10 * time.Second)
			continue
		}

		latestHeight = tipSet.Height()
		break
	}

	sleepTime := 30 * time.Second
	for {
		wait(sleepTime)

		newTipSet, err := m.getTipSetByHeight(latestHeight + 1)
		if newTipSet == nil || err != nil { // null block
			latestHeight++
			continue
		}

		latestHeight++
		m.send(latestHeight)

		latestTipSet, err := m.getLatestBlock()
		if err != nil {
			continue
		}

		if latestTipSet.Height() >= latestHeight+3 {
			sleepTime = 10 * time.Second
		} else {
			sleepTime = 30 * time.Second
		}
	}
}

func wait(sleepTime time.Duration) {
	time.Sleep(sleepTime)
}

func (m *Monitor) send(epoch abi.ChainEpoch) {
	m.lk.RLock()
	defer m.lk.RUnlock()

	for minerId, ch := range m.heightEvent {
		select {
		case ch <- epoch:
			blockLog.Warnw("send epoch success", "minerId", minerId)
		default:
			blockLog.Warnw("channel blocked", "minerId", minerId)
		}
	}
}
