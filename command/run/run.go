package runcommand

import (
	"os"

	"github.com/xmopen/godocker/container"
)

// run 真正执行run命令
func run(tti bool, cmd string) error {
	// PID为1的进程为/bin/sh(参数指定的cmd为/bin/sh)
	// parentProcess 已经是资源隔离的/bin/sh进程了
	parentProcess := container.NewParentProcess(tti, cmd)
	if err := parentProcess.Start(); err != nil {
		xlog.Errorf("parent process start err:[%+v]", err)
	}
	xlog.Infof("start wait")
	parentProcess.Wait()
	os.Exit(-1)
	xlog.Infof("run exit")
	return nil
}
