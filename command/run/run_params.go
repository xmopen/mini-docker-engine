package runcommand

import "github.com/urfave/cli"

// runCommandParameter  run 命令解析参数，用以对docker初始化进程的资源隔离和限制。
type runCommandParameter struct {
	ctx           cli.Context
	createTty     bool     // createTty 是否启用命令行伪终端。
	detach        bool     // detach 后台运行。
	imageName     string   // imageName 镜像名称。
	containerName string   // containerName 容器名称。
	volume        string   // volume 容器挂载点。
	network       string   // network 网络配置。
	env           []string // env 环境变量。
	port          []string // port 端口映射。
}

func nweRunCommandParameter(ctx *cli.Context) *runCommandParameter {
	return &runCommandParameter{
		createTty:     ctx.Bool("ti"),
		detach:        ctx.Bool("d"),
		containerName: ctx.String("name"),
		volume:        ctx.String("v"),
		network:       ctx.String("net"),
		env:           ctx.StringSlice("e"),
		port:          ctx.StringSlice("p"),
	}
}

func (r *runCommandParameter) setImageName(name string) {
	r.imageName = name
}

func (r *runCommandParameter) setContainerName(name string) {
	r.containerName = name
}
