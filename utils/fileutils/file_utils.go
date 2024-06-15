package fileutils

import (
	"fmt"
	"os"
	"path"
)

// GetCgroupPath 获取CGroup路径
func GetCgroupPath(subSystemName, cgoupPath string, autoCreate bool) (string, error) {
	cgroupRoot, err := FindCgroupMountPoint(subSystemName)
	if err != nil {
		return "", err
	}
	if _, err := os.Stat(path.Join(cgroupRoot, cgoupPath)); err != nil {
		if autoCreate && os.IsNotExist(err) {
			// 文件不存在则创建。
			if err := os.Mkdir(path.Join(cgroupRoot, cgoupPath), 0755); err != nil {
				return "", fmt.Errorf("create cgroup err:[%+v]", err)
			}
		}
		return "", err
	}

	return path.Join(cgroupRoot, cgoupPath), nil
}

// IOWriteFile 重写ioutil.WriteFile接口，如果文件不存在则创建
func IOWriteFile(path string, content []byte, perm os.FileMode) error {
	_, err := os.Stat(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		if err := os.Mkdir(path, 0755); err != nil {
			return err
		}
	}
	// 文件确定存在.
	return os.WriteFile(path, content, perm)
}

// PathExists 校验指定Path路径是否存在对应的文件。
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
