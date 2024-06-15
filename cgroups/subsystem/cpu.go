package subsystem

import (
	"io/ioutil"
	"os"
	"path"
	"strconv"

	"github.com/xmopen/godocker/utils/fileutils"
	"github.com/xmopen/golib/pkg/xlogging"
)

// CPUSubSystem CPU System
type CPUSubSystem struct {
	xlog *xlogging.Entry
}

// Name CPUSubSystem 子系统名称
func (c *CPUSubSystem) Name() string {
	return "cpu"
}

// Set 对cgroup设置CPU资源限制
func (c *CPUSubSystem) Set(cgroupPath string, res *ResourceConfig) error {
	c.xlog.Debugf("set cgrouppath:[%+v] resource:[%+v]", cgroupPath, res)
	subSystemPath, err := fileutils.GetCgroupPath(c.Name(), cgroupPath, true)
	if err != nil {
		return err
	}
	if res.CPUShare == "" {
		return nil
	}
	cpuShareSystemFilePath := path.Join(subSystemPath, "cpu.shares")
	if err := ioutil.WriteFile(cpuShareSystemFilePath, []byte(res.CPUShare), 0644); err != nil {
		c.xlog.Debugf("write file err:[%+v] filePath:[%+v] res:[%+v]", err, cpuShareSystemFilePath, res)
		return err
	}
	return nil
}

// Remove 移除掉某个CgroupPath中的CPUSubSystem限制
func (c *CPUSubSystem) Remove(cgroupPath string) error {
	c.xlog.Debugf("remove group:[%+v]", cgroupPath)
	subSystemPath, err := fileutils.GetCgroupPath(c.Name(), cgroupPath, true)
	if err != nil {
		return err
	}
	return os.Remove(subSystemPath)
}

// Apply 将PID对应的进程ID添加到对应的cgrouppath中
func (c *CPUSubSystem) Apply(cgroupPath string, pid int) error {
	c.xlog.Debugf("apply group:[%+v] pid:[%+v]", cgroupPath, pid)
	subSysGroupPath, err := fileutils.GetCgroupPath(c.Name(), cgroupPath, true)
	if err != nil {
		return err
	}
	cpuSubSystemInGroupPath := path.Join(subSysGroupPath, "tasks")
	if err := fileutils.IOWriteFile(cpuSubSystemInGroupPath, []byte(strconv.Itoa(pid)), 0644); err != nil {
		c.xlog.Errorf("writefile err:[%+v] path:[%+v] pid:[%+v]", err, cpuSubSystemInGroupPath, pid)
		return err
	}
	return nil
}
