package main

import (
	lcli "github.com/filecoin-project/lotus/cli"
	logging "github.com/ipfs/go-log/v2"
	"github.com/luluup777/robin/modules"
	"github.com/luluup777/robin/parse"
	"github.com/urfave/cli/v2"
	"os"
	"os/signal"
	"syscall"
)

var log = logging.Logger("robin")

func main() {
	_ = logging.SetLogLevel("*", "INFO")

	app := &cli.App{
		Name:                 "robin",
		Usage:                "mining monitoring and alarm",
		Version:              "v0.1",
		EnableBashCompletion: true,
		Commands: []*cli.Command{
			run,
		},
	}

	app.Setup()
	lcli.RunApp(app)
}

var run = &cli.Command{
	Name:  "run",
	Usage: "start robin",
	Before: func(context *cli.Context) error {
		parse.Init()
		return nil
	},
	Action: func(cctx *cli.Context) error {
		log.Info("robin start")

		go modules.StartMonitor()

		quit := make(chan os.Signal)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		log.Info("robin shutting down")
		return nil
	},
}
