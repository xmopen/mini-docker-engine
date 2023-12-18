package subsystem

import (
	"os"
	"path"
	"strconv"
)

// MemorySubSystem 内存子系统
type MemorySubSystem struct {
}

// Set MemorySubSytem设置内存限制值
func (m *MemorySubSystem) Set(cgroupPath string, res *ResourceConfig) error {
	xlog.Debugf("memory subsystem set cgroupPath:[%+v] res:[%+v]", cgroupPath, res)
	subSysCgroupPath, err := GetCgroupPath(m.Name(), cgroupPath, true)
	if err != nil {
		return err
	}
	if res.MemoryLimit == "" {
		return nil
	}
	// FIXME: 1.16之后如果文件不存在则不进行创建
	// 设置当前Cgrop的内存限制，值为配置的值
	err = IOWriteFile(path.Join(subSysCgroupPath, "memory.limit_in_bytes"), []byte(res.MemoryLimit), 0644)
	if err != nil {
		return err
	}
	return nil
}

// Name MemorySubSystem 名称
func (m *MemorySubSystem) Name() string {
	return "memory"
}

// Apply 将一个进程添加到path所对应的cgroup
func (m *MemorySubSystem) Apply(cgroupPath string, pid int) error {
	subSysCgroupPath, err := GetCgroupPath(m.Name(), cgroupPath, false)
	if err != nil {
		return err
	}
	// 将进程的PID写入到cgroup的虚拟文件系统对应的task文件中
	return IOWriteFile(path.Join(subSysCgroupPath, "tasks"), []byte(strconv.Itoa(pid)), 0644)
}

// Remove 删除cgroupPath对应的cgroup
func (m *MemorySubSystem) Remove(cgroupPath string) error {
	subSysCgroupPath, err := GetCgroupPath(m.Name(), cgroupPath, false)
	if err != nil {
		return err
	}
	return os.Remove(subSysCgroupPath)
}
