package runcommand

import (
	"fmt"

	"github.com/urfave/cli"
	"github.com/xmopen/godocker/cgroups/subsystem"
	"github.com/xmopen/golib/pkg/xlogging"
)

var xlog = xlogging.Tag("command.run")

// Command run command
var Command = cli.Command{
	Name:  "run",
	Usage: "Create a container with namespace and cgroups limit",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "ti",
			Usage: "enable tti",
		},
		cli.BoolFlag{
			Name:  "m",
			Usage: "memory limit",
		},
		cli.BoolFlag{
			Name:  "cpushare",
			Usage: "cpushare limit",
		},
		cli.BoolFlag{
			Name:  "cpuset",
			Usage: "cpuset limit",
		},
	},
	Action: runCommandAction,
}

// runCommandAction run 命令执行函数
// docker run imgae cmd...
// -d 后台运行，-t 伪终端 -i标准输入输出
func runCommandAction(ctx *cli.Context) error {
	if len(ctx.Args()) < 1 {
		return fmt.Errorf("run command missiong args")
	}
	cmdArray := make([]string, 0)
	for _, arg := range ctx.Args() {
		cmdArray = append(cmdArray, arg)
	}
	runParameter := nweRunCommandParameter(ctx)
	if runParameter.createTty && runParameter.detach {
		// 这两个命令不能同时存在。
		return fmt.Errorf("ti and d paramter can not both provided")
	}
	runParameter.setImageName(cmdArray[0])
	cmdArray = cmdArray[1:]
	resourceConfig := &subsystem.ResourceConfig{
		MemoryLimit: ctx.String("m"),
		CPUShare:    ctx.String("cpushare"),
		CPUSet:      ctx.String("cpuset"),
	}
	xlog.Debugf("run command resource:[%+v],cmd:[%+v],param:[%+v]", resourceConfig, cmdArray, runParameter)
	return run(ctx.Bool("ti"), cmdArray, resourceConfig, runParameter)
}
