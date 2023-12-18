package runcommand

import (
	"encoding/json"
	"fmt"

	"github.com/urfave/cli"
	"github.com/xmopen/golib/pkg/xlogging"
)

var xlog = xlogging.Tag("command.run")

// Command run command
var Command = cli.Command{
	Name: "run",
	Usage: "Create a container with namespace and cgroups limit",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name: "ti",
			Usage: "enable tti",
		},
	},
	Action: runCommandAction,
}

// runCommandAction run 命令执行函数
func runCommandAction(ctx *cli.Context) error {
	if len(ctx.Args()) < 1 {
		return fmt.Errorf("run command missiong args")
	}
	d,_ := json.Marshal(ctx.Args())
	xlog.Debugf("run command args:[%+v]",string(d))
	return run(ctx.Bool("ti"),ctx.Args().Get(0))
}