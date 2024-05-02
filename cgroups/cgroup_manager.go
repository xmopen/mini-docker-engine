package cgroups

import (
	"github.com/xmopen/godocker/cgroups/subsystem"
	"github.com/xmopen/golib/pkg/xlogging"
)

var xlog = xlogging.Tag("cgroups")

// CGroupManager CGroup管理器
type CGroupManager struct {
	// Path CGroup在Hierarchy(分层)中的路径
	Path string
	// Resource 当前CGroup资源配置文件
	Resource *subsystem.ResourceConfig
}

// NewCGroupManager 构造CGroupManager管理器
func NewCGroupManager(path string) *CGroupManager {
	return &CGroupManager{
		Path: path,
	}
}

// Apply 将进程PID添加到对应的cgroup分组上
func (c *CGroupManager) Apply(pid int) {
	for _, item := range subsystem.SubSystemInstances {
		if err := item.Apply(c.Path, pid); err != nil {
			xlog.Errorf("apply pid err,pid:[%+v] cgroup path:[%+v]", pid, c.Path)
		}
	}
}

// Set 对已有的cgroup分组设置资源
func (c *CGroupManager) Set(res *subsystem.ResourceConfig) {
	for _, item := range subsystem.SubSystemInstances {
		if err := item.Set(c.Path, res); err != nil {
			xlog.Errorf("set err,cgroup path:[%+v] resource:[%+v]", c.Path, res)
		}
	}
}

// Destroy 销毁所有CGroups
func (c *CGroupManager) Destroy() {
	for _, item := range subsystem.SubSystemInstances {
		if err := item.Remove(c.Path); err != nil {
			xlog.Errorf("destroy err,cgroup path:[%+v]", err)
		}
	}
}
