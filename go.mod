module github.com/luluup777/robin

go 1.16

require (
	github.com/filecoin-project/go-address v0.0.6
	github.com/filecoin-project/go-state-types v0.1.1
	github.com/filecoin-project/lotus v1.13.1
	github.com/ipfs/go-log/v2 v2.3.0
	github.com/spf13/viper v1.3.2
	github.com/urfave/cli/v2 v2.3.0
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1
)

replace github.com/filecoin-project/lotus => ./extern/lotus

replace github.com/filecoin-project/filecoin-ffi => ./extern/lotus/extern/filecoin-ffi
