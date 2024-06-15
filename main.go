package main

import (
	"os"

	initcommand "github.com/xmopen/godocker/command/init"
	runcommand "github.com/xmopen/godocker/command/run"
	"github.com/urfave/cli"
	"github.com/xmopen/golib/pkg/xlogging"
)

// TODO: cli 还需要抽时间看下.
// NameSpace 资源隔离。
// Cgroups   资源限制。

var xlog = xlogging.Tag("godocker.main")

// 1、通过 docker run 命令启动docker进程
// 2、在Run命令处理函数中， 通过系统调用自己执行自己，并且带上init参数执行init回调。
// 3、在init对资源进行初始化完之后，
func main() {
	app := cli.NewApp()
	app.Name = "godocker"
	app.Usage = ""
	app.Commands = []cli.Command{
		runcommand.Command,
		initcommand.Command,
	}
	app.Before = func(ctx *cli.Context) error {
		xlog.Infof("docker run before.")
		return nil
	}
	if err := app.Run(os.Args); err != nil {
		xlog.Errorf("run err:[%+v]", err)
	}
}
