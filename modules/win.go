package modules

import (
	"context"
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/lotus/chain/gen"
	"github.com/filecoin-project/lotus/chain/types"
	logging "github.com/ipfs/go-log/v2"
)

var winLog = logging.Logger("win")

func (m *Monitor) isWin(epoch abi.ChainEpoch, minerAddr address.Address) bool {
	ctx := context.Background()
	round := epoch
	tp, err := m.api.ChainGetTipSetByHeight(ctx, epoch-1, types.NewTipSetKey())
	if err != nil {
		winLog.Errorw("IsWin:ChainGetTipSetByHeight", "err", err)
		return false
	}

	mbi, err := m.api.MinerGetBaseInfo(ctx, minerAddr, round, tp.Key())
	if err != nil {
		winLog.Errorw("IsWin:MinerGetBaseInfo", "err", err)
		return false
	}

	if mbi == nil {
		return false
	}
	if !mbi.EligibleForMining {
		// slashed or just have no power yet
		return false
	}

	beaconPrev := mbi.PrevBeaconEntry
	bvals := mbi.BeaconEntries

	rbase := beaconPrev
	if len(bvals) > 0 {
		rbase = bvals[len(bvals)-1]
	}

	winner, err := gen.IsRoundWinner(ctx, tp, round, minerAddr, rbase, mbi, m.api)
	if err != nil {
		return false
	}

	if winner == nil {
		return false
	}

	return true
}

func isSuccessWin(tipSet *types.TipSet, minerAddr address.Address) bool {
	for _, block := range tipSet.Blocks() {
		if block.Miner == minerAddr {
			return true
		}
	}

	return false
}
