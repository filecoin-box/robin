package modules

import (
	"github.com/filecoin-project/go-state-types/abi"
)

var delayEpoch = abi.ChainEpoch(10)

type WinBlockInfo struct {
	WinEpoch   abi.ChainEpoch
	CheckEpoch abi.ChainEpoch
}

type Queue []*WinBlockInfo

func (q *Queue) NoEmpty() bool {
	return len(*q) != 0
}

func (q *Queue) GetHead() *WinBlockInfo {
	return (*q)[0]
}

func (q *Queue) Push(w *WinBlockInfo) {
	*q = append(*q, w)
}

func (q *Queue) Pop() *WinBlockInfo {
	w := (*q)[0]
	*q = (*q)[1:]
	return w
}
