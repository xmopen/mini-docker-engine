package container

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/xmopen/godocker/utils/fileutils"
)

// WorkSpace 容器工作空间。
type WorkSpace struct {
	ImageName     string
	ContainerName string
	Volume        string
}

// NewWorkSpace 构造容器工作空间。
func NewWorkSpace(imageName, volume, containerName string) *WorkSpace {
	return &WorkSpace{
		ImageName:     imageName,
		ContainerName: containerName,
		Volume:        volume,
	}
}

// Init 初始化容器工作空间。
// 容器内的文件目录位置需要单独进行挂载。
func (w *WorkSpace) Init() error {
	w.createReadOnlyLayer()
	w.createWriteLayer()
	w.createMountPoint()
	// 设置挂载点。
	if w.Volume != "" {
		volumePaths := strings.Split(w.Volume, ":")
		if len(volumePaths) == 2 && volumePaths[0] != "" && volumePaths[1] != "" {
			// 进行挂载。
			if err := w.mountVolume(volumePaths); err != nil {
				xlog.Errorf("workspace mount volume err:[%+v] path:[%+v]", err, volumePaths)
				return err
			}
		} else {
			xlog.Warnf("Volume parameter input is not correct.")
		}
	}
	return nil
}

func (w *WorkSpace) createReadOnlyLayer() error {
	tarFolderPath := fmt.Sprintf("%s/%s/", RootURL, w.ImageName)
	exists, err := fileutils.PathExists(tarFolderPath)
	if err != nil {
		xlog.Errorf("create read only layer err:[%+v] path:[%+v]", err, tarFolderPath)
		return err
	}
	if exists {
		return nil
	}
	// 不存在则创建。
	if err := os.MkdirAll(tarFolderPath, 0622); err != nil {
		return err
	}

	imageURL := fmt.Sprintf("%s/%s.tar", RootURL, w.ImageName)
	// TODO: 疑问。
	// 执行tar命令，并且获取到对应的输出流。
	if _, err := exec.Command("tar", "-xvf", imageURL, "-C", tarFolderPath).CombinedOutput(); err != nil {
		return err
	}
	return nil
}

func (w *WorkSpace) createWriteLayer() error {
	writePath := fmt.Sprintf(WriteLayerUrl, w.ContainerName)
	// 目录权限：八进制，用户、用户组、其他用户，4：读权限 2：写权限 1：执行权限。
	// 7: 0111 表示不管是用户还是用户组或者是其他用户都有读写执行权限。
	if err := os.MkdirAll(writePath, 0777); err != nil {
		return err
	}
	return nil
}

// createMountPoint 创建挂载点。
// TODO: 这里创建的挂载点事什么意思：针对进程创建容器的文件挂载点。
func (w *WorkSpace) createMountPoint() error {
	mountPath := fmt.Sprintf(MountBasePath, w.ContainerName)
	if err := os.MkdirAll(mountPath, 0777); err != nil {
		return err
	}
	// TODO: 挂载的命令应该抽象出来。
	tempWriteLayer := fmt.Sprintf(WriteLayerUrl, w.ContainerName)
	tempImageLocation := fmt.Sprintf("%s/%s", RootURL, w.ImageName)
	dirs := fmt.Sprintf("dirs=%s:%s", tempWriteLayer, tempImageLocation)
	if _, err := exec.Command("mount", "-t", "aufs", "-o", dirs, "none", mountPath).CombinedOutput(); err != nil {
		return err
	}
	return nil
}

// mountVolume 挂载配置的挂载点。
func (w *WorkSpace) mountVolume(volumePaths []string) error {
	// TODO: 后续优化抽象，先进行实现。
	parentPath := volumePaths[0]
	if err := os.MkdirAll(parentPath, 0777); err != nil {
		return err
	}
	containerPath := volumePaths[1]
	mountPath := fmt.Sprintf(MountBasePath, containerPath)
	containerVolumePath := fmt.Sprintf("%s/%s", mountPath, containerPath)
	if err := os.MkdirAll(containerVolumePath, 0777); err != nil {
		return err
	}
	dirs := fmt.Sprintf("dirs=%s", parentPath)
	if _, err := exec.Command("mount", "-t", "aufs", "-o", dirs, "none", containerVolumePath).CombinedOutput(); err != nil {
		return err
	}
	return nil
}
