package container

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/xmopen/godocker/utils/processutils"
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
func NewParentProcess(tty bool, containerName, volume, imageName string, env []string) (*exec.Cmd, *os.File) {
	readPipe, writePipe, err := NewPipe()
	if err != nil {
		return nil, nil
	}
	// 自己调用自己
	xlog.Infof("new parent process containerName:[%+v] volume:[%+v] imageName:[%+v] env:[%+v] tty:[%+v]",
		containerName, volume, imageName, env, tty)
	// 当前进程自己调用自己会带上init函数
	// args := []string{"init", command}
	initCmd, err := os.Readlink(processutils.ProcessSelfExe)
	if err != nil {
		xlog.Errorf("read link err:[%+v] path:[%+v]", err, processutils.ProcessSelfExe)
		return nil, nil
	}
	// /proc/self/exe 链接到当前正在运行的进程
	// /proc/self/exe 链接到当前运行进程的执行命令文件
	// TODO: 自己调用自己的时候对当前进程进行指定Flag资源隔离
	// godocker init run -ti /bin/sh
	cmd := exec.Command(initCmd, "init")
	// 对init进程进行资源隔离，这个init进程也就是pid=0的进程。
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWNET |
			syscall.CLONE_NEWIPC,
	}
	// 设置控制台输出
	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	} else {
		// 如果当前输出不是控制台输出，那么输出存入到指定默认的文件中。
		dirURL := fmt.Sprintf(DefaultInfoLocation, containerName)
		if err := os.MkdirAll(dirURL, 0622); err != nil {
			xlog.Errorf("mkdir err:[%+v] path:[%+v]", err, dirURL)
			return nil, nil
		}
		stdLogFile, err := os.Create(dirURL + ContainerLogFile)
		if err != nil {
			xlog.Errorf("create stdlog file err:[%+v] path:[%+v]", err, dirURL+ContainerLogFile)
			return nil, nil
		}
		cmd.Stdout = stdLogFile
	}

	cmd.ExtraFiles = []*os.File{readPipe}
	cmd.Env = append(cmd.Env, env...)
	if err := NewWorkSpace(imageName, volume, containerName).Init(); err != nil {
		xlog.Errorf("new work space err:[%+v] image:[%+v] volume:[%+v] container:[%+v]", err, imageName, volume,
			containerName)
		return nil, nil
	}
	// 设置进程的工作目录为挂载点。
	cmd.Dir = fmt.Sprintf(MountBasePath,containerName)
	return cmd, writePipe
}
