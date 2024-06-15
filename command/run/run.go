package runcommand

import (
	"fmt"
	"os"

	"github.com/xmopen/godocker/cgroups"
	"github.com/xmopen/godocker/cgroups/subsystem"
	"github.com/xmopen/godocker/container"
	"github.com/xmopen/godocker/utils/random"
)

// run 真正执行run命令
func run(tti bool, cmdArray []string, resource *subsystem.ResourceConfig, runParameter *runCommandParameter) error {
	containerID := random.RandNumberToString(container.ContainerIDLength)
	if runParameter.containerName == "" {
		xlog.Warnf("container name is empty, use name:[%+v]", containerID)
		runParameter.setContainerName(containerID)
	}
	// PID为1的进程为/bin/sh(参数指定的cmd为/bin/sh)
	// parentProcess 已经是资源隔离的/bin/sh进程了
	parentProcess, writePipe := container.NewParentProcess(tti, runParameter.containerName, runParameter.volume, runParameter.imageName, runParameter.env)
	xlog.Infof("%v", writePipe)
	if parentProcess == nil {
		return fmt.Errorf("parent process is nil")
	}
	// 执行刚刚初始化的父进程，也就是init进程。
	if err := parentProcess.Start(); err != nil {
		// 执行失败退出就好了，没必要继续执行。这里的失败也就意味着命名空间限制失败。
		xlog.Errorf("parent process start err:[%+v]", err)
		return err
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
