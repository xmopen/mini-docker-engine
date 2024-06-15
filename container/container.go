package container

import "github.com/xmopen/golib/pkg/xlogging"

// Container 常量
const (
	// ContainerIDLength 容器ID长度
	ContainerIDLength = 10
	// DefaultInfoLocation 默认的输出路径。
	DefaultInfoLocation string = "/var/run/mydocker/%s/"
	// ContainerLogFile 容器默认输出路径。
	ContainerLogFile string = "container.log"
	// RootURL 容器根路径。
	RootURL string = "/root"
	// WriteLayerUrl .
	WriteLayerUrl string = "/root/writeLayer/%s"
	// MountBasePath 挂载点基础路径。
	MountBasePath string = "/root/mnt/%s"
)

var xlog = xlogging.Tag("container")
