package main

import (
	"os"

	initcommand "github.com/xmopen/godocker/command/init"
	runcommand "github.com/xmopen/godocker/command/run"
	"github.com/urfave/cli"
	"github.com/xmopen/golib/pkg/xlogging"
)

// TODO: cli 还需要抽时间看下.

var xlog = xlogging.Tag("godocker.main")

func main() {
	app := cli.NewApp()
	app.Name = "godocker"
	app.Usage = ""
	app.Commands = []cli.Command{
		runcommand.Command,
		initcommand.Command,
	}
	app.Before = func(ctx *cli.Context) error {
		return nil
	}
	if err := app.Run(os.Args); err != nil {
		xlog.Errorf("run err:[%+v]", err)
	}
}
