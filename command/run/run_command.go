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
func runCommandAction(ctx *cli.Context) error {
	if len(ctx.Args()) < 1 {
		return fmt.Errorf("run command missiong args")
	}
	resource := &subsystem.ResourceConfig{
		MemoryLimit: ctx.String("m"),
		CPUShare:    ctx.String("cpushare"),
		CPUSet:      ctx.String("cpuset"),
	}
	cmdArray := make([]string, 0)
	for _, arg := range ctx.Args() {
		cmdArray = append(cmdArray, arg)
	}
	xlog.Debugf("run command resource:[%+v]", resource)
	return run(ctx.Bool("ti"), cmdArray, resource)
}
