package modules

import (
	"context"
	"fmt"
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/lotus/api/v1api"
	logging "github.com/ipfs/go-log/v2"
	"github.com/luluup777/robin/parse"
	"os"
	"strings"
	"sync"
)

var monitorLog = logging.Logger("monitor")

type Monitor struct {
	lk sync.RWMutex

	minerIds    map[string]context.CancelFunc
	heightEvent map[string]chan abi.ChainEpoch

	config            parse.Config
	monitorConfChange chan struct{}
	notifyConfChange  chan struct{}

	notify chan string

	api v1api.FullNode
}

func StartMonitor() {
	c := parse.GetRobinConfig()
	api, err := NewFullNodeApi(c.Monitor.Fullnode_api_info)
	if err != nil {
		monitorLog.Errorw("NewFullNodeApi", "err", err)
		os.Exit(1)
	}

	m := &Monitor{
		minerIds:          map[string]context.CancelFunc{},
		heightEvent:       map[string]chan abi.ChainEpoch{},
		config:            c,
		monitorConfChange: make(chan struct{}, 1),
		notifyConfChange:  make(chan struct{}, 1),
		notify:            make(chan string, 100),
		api:               api,
	}

	go m.blockMonitor()
	go m.watchConfig()
	go m.robin()
	go m.run()
}

func (m *Monitor) run() {
	minerIds := strings.Split(m.config.Monitor.MinerId, ",")
	for _, minerId := range minerIds {
		epochCh := make(chan abi.ChainEpoch, 1)
		ctx, cancel := context.WithCancel(context.Background())

		m.addMonitor(minerId, cancel, epochCh)

		go m.monitor(ctx, epochCh, minerId)
	}

	for {
		select {
		case <-m.monitorConfChange:
			nowMinerIds := strings.Split(m.config.Monitor.MinerId, ",")
			for _, minerId := range minerIds {
				isCancel := true
				for _, nowMinerId := range nowMinerIds {
					if minerId == nowMinerId {
						isCancel = false
						break
					}
				}

				if isCancel {
					monitorLog.Infow("cancel monitor", "minerId", minerId)
					m.cancelMonitor(minerId)
				}
			}

			for _, nowMinerId := range nowMinerIds {
				isAdd := true
				for _, minerId := range minerIds {
					if minerId == nowMinerId {
						isAdd = false
						break
					}
				}

				if isAdd {
					monitorLog.Infow("add monitor", "minerId", nowMinerId)

					epochCh := make(chan abi.ChainEpoch, 1)
					ctx, cancel := context.WithCancel(context.Background())
					m.addMonitor(nowMinerId, cancel, epochCh)

					go m.monitor(ctx, epochCh, nowMinerId)
				}
			}

			newMinerIds := nowMinerIds
			minerIds = newMinerIds
		}
	}
}

func (m *Monitor) addMonitor(minerId string, cancel context.CancelFunc, epochCh chan abi.ChainEpoch) {
	m.lk.Lock()
	defer m.lk.Unlock()

	m.minerIds[minerId] = cancel
	m.heightEvent[minerId] = epochCh
}

func (m *Monitor) cancelMonitor(minerId string) {
	m.lk.Lock()
	defer m.lk.Unlock()
	if cancel, ok := m.minerIds[minerId]; ok {
		cancel()
	}

	delete(m.minerIds, minerId)
	delete(m.heightEvent, minerId)
}

func (m *Monitor) monitor(ctx context.Context, epochCh chan abi.ChainEpoch, minerId string) {
	monitorLog.Infow("start monitor", "minerId", minerId)

	minerAddr, err := address.NewFromString(minerId)
	if err != nil {
		monitorLog.Errorw("address.NewFromString", "err", err)
		m.cancelMonitor(minerId)
		return
	}

	var winBlocks = &Queue{}
	for {
		select {
		case <-ctx.Done():
			monitorLog.Infow("exit monitor", "minerId", minerId)
			return
		case epoch := <-epochCh:
			monitorLog.Infow("block received", "minerId", minerId)

			if m.isWin(epoch, minerAddr) {
				monitorLog.Infow("win block", "minerId", minerId, "epoch", epoch)
				winBlocks.Push(&WinBlockInfo{
					WinEpoch:   epoch,
					CheckEpoch: epoch + delayEpoch,
				})
			}

			if winBlocks.NoEmpty() {
				if winBlocks.GetHead().CheckEpoch == epoch {
					win := winBlocks.Pop()
					tipSet, err := m.getTipSetByHeight(win.WinEpoch)
					if err != nil {
						monitorLog.Errorw("getTipSetByHeight", "err", err, "minerId", minerId)
						continue
					}

					msg := ""
					if isSuccessWin(tipSet, minerAddr) {
						monitorLog.Infow("success win block", "minerId", minerId, "epoch", epoch)
						msg = fmt.Sprintf("GOOD!!! %s win block at %d height", minerId, win.WinEpoch)
					} else {
						monitorLog.Warnw("loss block", "minerId", minerId, "epoch", epoch)
						msg = fmt.Sprintf("BAD!!! %s loss block at %d height", minerId, win.WinEpoch)
					}

					m.notify <- msg
				}
			}
		}
	}
}
