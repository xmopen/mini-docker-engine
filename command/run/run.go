package runcommand

import (
	"os"

	"github.com/xmopen/godocker/cgroups"
	"github.com/xmopen/godocker/cgroups/subsystem"
	"github.com/xmopen/godocker/container"
	"github.com/xmopen/godocker/utils/random"
)

// run 真正执行run命令
func run(tti bool, cmdArray []string, resource *subsystem.ResourceConfig) error {

	containerID := random.RandNumberToString(container.ContainerIDLength)
	// PID为1的进程为/bin/sh(参数指定的cmd为/bin/sh)
	// parentProcess 已经是资源隔离的/bin/sh进程了
	parentProcess := container.NewParentProcess(tti, cmdArray[0])
	if err := parentProcess.Start(); err != nil {
		xlog.Errorf("parent process start err:[%+v]", err)
	}

	// cgroup
	cgroupManager := cgroups.NewCGroupManager(containerID)
	cgroupManager.Set(resource)
	cgroupManager.Apply(parentProcess.Process.Pid)
	xlog.Infof("start wait")
	parentProcess.Wait()
	os.Exit(-1)
	xlog.Infof("run exit")
	return nil
}
