package subsystem

import (
	"os"
	"path"
	"strconv"

	"github.com/xmopen/golib/pkg/xlogging"
)

// CPUSetSubSystem 设置容器进程CPU占用多少核
type CPUSetSubSystem struct {
	xlog *xlogging.Entry
}

// Name CPUSetSubSystem Name
func (c *CPUSetSubSystem) Name() string {
	return "cpuset"
}

// Set 设置CGroup对应的CPUSetSubSystem
func (c *CPUSetSubSystem) Set(cgroupPath string, res *ResourceConfig) error {
	c.xlog.Debugf("set cgrouppath:[%+v] res:[%+v]", cgroupPath, res)
	groupSubsytemPath, err := GetCgroupPath(c.Name(), cgroupPath, true)
	if err != nil {
		return err
	}
	if res.CPUSet == "" {
		return nil
	}
	cgroupSubsystemCPUSetPath := path.Join(groupSubsytemPath, "cpuset.cpus")
	if err := IOWriteFile(cgroupSubsystemCPUSetPath, []byte(res.CPUSet), 0644); err != nil {
		c.xlog.Errorf("iowritefile err:[%+v] path:[%+v] res:[%+v]", err, cgroupSubsystemCPUSetPath, res)
		return err
	}
	return nil
}

// Apply
func (c *CPUSetSubSystem) Apply(cgroupPath string, pid int) error {
	c.xlog.Debugf("apply cgrouppath:[%+v] pid:[%+v]", cgroupPath, pid)
	subSystemPath, err := GetCgroupPath(c.Name(), cgroupPath, true)
	if err != nil {
		return err
	}
	subSystemTasksPath := path.Join(subSystemPath, "tasks")
	if err := IOWriteFile(subSystemTasksPath, []byte(strconv.Itoa(pid)), 0644); err != nil {
		c.xlog.Debugf("write file err:[%+v] path:[%+v] pid:[%+v]", err, subSystemTasksPath, pid)
		return err
	}
	return nil
}

func (c *CPUSetSubSystem) Remove(cgroupPath string) error {
	c.xlog.Debugf("remove path:[%+v]", cgroupPath)
	subSysGroupPath, err := GetCgroupPath(c.Name(), cgroupPath, true)
	if err != nil {
		return err
	}
	return os.Remove(subSysGroupPath)
}
