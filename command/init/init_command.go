package initcommand

import (
	"github.com/urfave/cli"
	"github.com/xmopen/godocker/container"
	"github.com/xmopen/golib/pkg/xlogging"
)

var xlog = xlogging.Tag("command.init")

// Command 定义InitCommand，内部操作，禁止外部调用
// 当前进程在调用自己的同时会带上init参数来执行当前函数
// 自己调用自己在run/ /proc/self/exe 自己调用自己
var Command = cli.Command{
	Name:   "init",
	Usage:  "Init container process run user's process in container. Do not call it outside",
	Action: initCommandAction,
}

// initCommandAction init进程执行，执行到这里已经是container中的init进程
func initCommandAction(ctx *cli.Context) error {
	xlog.Infof("init command action")
	// 真正本次容器要执行的命令cmd,比如/bin/sh
	cmd := ctx.Args().Get(0)
	return container.RunContainerInitProcess(cmd, nil)
}
