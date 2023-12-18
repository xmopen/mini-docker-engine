package subsystem

import "github.com/xmopen/golib/pkg/xlogging"

var xlog = xlogging.Tag("subsytem")

var (
	SubSystemInstances = []Subsystem{
		&CPUSetSubSystem{
			xlog: xlogging.Tag("cpuset.subsystem"),
		},
		&CPUSubSystem{
			xlog: xlogging.Tag("cpu.subsystem"),
		},
		&MemorySubSystem{},
	}
)

// ResourceConfig 容器资源限制：内存限制、CCPU时间片权重、CPU核心数
type ResourceConfig struct {
	MemoryLimit string
	CPUShare    string
	CPUSet      string
}

// Subsystem Subsystem interface
// 将cgroup抽象成path,是因为cgroup在hierarchy中的路径，便是虚拟文件系统中的虚拟路径
type Subsystem interface {
	// Name 返回子系统的名称
	Name() string
	// Set 设置CGroup在Subsystem中的资源限制
	Set(path string, res *ResourceConfig) error
	// Apply 将进程添加到某个cgroup中
	Apply(path string, pid int) error
	// 移除掉某个cgroup
	Remove(path string) error
}
