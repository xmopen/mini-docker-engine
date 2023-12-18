package container

import (
	"os"
	"os/exec"
	"syscall"
)

// NewParentProcess 初始化父进程: Init 进程
// 创建NameSpace隔离的容器进程
// Linux Namespace are a feature of the Linux Kernel that partitions kernel resources such that one set of processes
// sees one set of resources while another set of processes sees a different set of resources
// 1、CLONE_NEWUTS 用来隔离Nodename(hostname)或者Domainname两个系统概念
// 2、CLONE_NEWPID 用来隔离进程PID
// 3、CLONE_NEWNS  用来隔离Mount隔离
// 4、CLONE_NEWNET 用来隔离Network
// 5、CLONE_NEWIPC 用来隔离IPC进程之间通信
func NewParentProcess(tty bool, command string) *exec.Cmd {
	// 自己调用自己
	xlog.Infof("new parent process command:[%+v] tty:[%+v]", command, tty)
	// 当前进程自己调用自己会带上init函数
	args := []string{"init", command}
	// /proc/self/exe 链接到当前正在运行的进程
	// /proc/self/exe 链接到当前运行进程的执行命令文件
	// TODO: 自己调用自己的时候对当前进程进行指定Flag资源隔离
	// godocker init run -ti /bin/sh
	cmd := exec.Command("/proc/self/exe", args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWNET |
			syscall.CLONE_NEWIPC,
	}
	// 设置控制台输出
	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	return cmd
}
