package container

import (
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

// defaultMoutFlags 默认挂载权限标签
// MS_NOEXEC        在本文件系统中不允许运行其他程序
// MS_NOSUID  		在本文件系统中不允许进行SetUserID
// MS_NODEV    		默认都会带入的参数
// MS_PRIVATE  		Namespace 私有Mount
// MS_REC
var defaultMoutFlags = syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV

// RunContainerInitProcess 运行容器内Init进程
// TODO: InitCommand 什么时候调用？
// FIXME: 参考： MountNamespace:
// https://www.sohu.com/a/260181668_467784
func RunContainerInitProcess(cmd string, obj any) error {
	rootPwd, err := os.Getwd()
	if err != nil {
		return err
	}
	xlog.Debugf("private root start")
	if err := privateRoot(rootPwd); err != nil {
		xlog.Errorf("private root err:[%+v] root:[%+v]", err, rootPwd)
	}
	xlog.Debugf("run container init process command:[%+v]", cmd)
	if err := syscall.Mount("", "/", "", syscall.MS_PRIVATE|syscall.MS_REC, ""); err != nil {
		xlog.Errorf("clear old mount info")
		return err
	}
	// systemd.Mount加入Linux之后，Mount Namespace就变成shared by default,所以必须要显示声明新的Mount namespace是独立的
	syscall.Mount("proc", "/proc", "proc", uintptr(defaultMoutFlags), "")
	argv := []string{cmd}
	// syscall.Exec: 执行当前filename对应的程序，覆盖当前进程的镜像、数据和堆栈等信息，包括PID，这些都会被将要运行的进程替换掉，
	// 这样当我们进入容器内进行就会发现第一个进程就是我们自己的要运行的进程
	if err := syscall.Exec(cmd, argv, os.Environ()); err != nil {
		xlog.Errorf("syscall.Exec err:[%+v] command:[%+v]", err, cmd)
	}
	return nil
}

// privateRoot
func privateRoot(root string) error {
	// 使当前老的Root和新的Root不在同一个文件系统下，重新把root挂载一次，bind mount就是把相同的内容换一个挂载点的挂载函数
	if err := syscall.Mount(root, root, "bind", syscall.MS_BIND|syscall.MS_REC, ""); err != nil {
		return err
	}
	// 创建rootfs/.pivot_root存储old_root
	pivotDir := filepath.Join(root, ".pivot_root")
	xlog.Debugf("privotDir:[%+v]", pivotDir)
	if err := os.Mkdir(pivotDir, 0777); err != nil {
		// file exists
		if !strings.Contains(err.Error(), "file exists") {
			return err
		}
	}

	if err := syscall.PivotRoot(root, pivotDir); err != nil {
		return err
	}
	// 修改当前的工作目录到根目录
	if err := syscall.Chdir("/"); err != nil {
		return err
	}
	pivotDir = filepath.Join("/", ".pivot_root")
	xlog.Debugf("privotDir:[%+v]", pivotDir)
	// unmount rootfs/.pivot_root
	if err := syscall.Unmount(pivotDir, syscall.MNT_DETACH); err != nil {
		return err
	}

	return os.Remove(pivotDir)
}
