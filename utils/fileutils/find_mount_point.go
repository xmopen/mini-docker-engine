package fileutils

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/xmopen/godocker/utils/processutils"
)

// FindCgroupMountPoint 获取进程内的指定子系统的挂载点。
func FindCgroupMountPoint(subsystem string) (string, error) {
	processFileInfo, err := os.Open(processutils.ProcessSelfMountInfo)
	if err != nil {
		return "", err
	}
	defer processFileInfo.Close()

	// 按照指定格式内容读取文件。
	scanner := bufio.NewScanner(processFileInfo)
	for scanner.Scan() {
		txt := scanner.Text()
		fields := strings.Split(txt, " ")
		for _, opt := range strings.Split(fields[len(fields)-1], ",") {
			if opt == subsystem {
				return fields[4], nil
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}
	return "", fmt.Errorf("not found subsystem mount point:[%+v]", subsystem)
}
