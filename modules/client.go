package modules

import (
	"context"
	"github.com/filecoin-project/lotus/api/client"
	"github.com/filecoin-project/lotus/api/v1api"
	cliutil "github.com/filecoin-project/lotus/cli/util"
	logging "github.com/ipfs/go-log/v2"
)

var clientLog = logging.Logger("client")

func NewFullNodeApi(fullnode_api_info string) (v1api.FullNode, error) {
	ainfo := cliutil.ParseApiInfo(fullnode_api_info)
	addr, err := ainfo.DialArgs("v1")
	if err != nil {
		clientLog.Errorw("DialArgs", "err", err)
		return nil, err
	}

	fullNodeApi, _, err := client.NewFullNodeRPCV1(context.Background(), addr, ainfo.AuthHeader())
	if err != nil {
		clientLog.Errorw("NewFullNodeRPCV1", "err", err)
		return nil, err
	}

	return fullNodeApi, nil
}
