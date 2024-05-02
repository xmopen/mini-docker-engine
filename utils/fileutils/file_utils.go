package fileutils

import "os"

// GetCgroupPath 获取CGroup路径
func GetCgroupPath(subSystemName, cgoupPath string, b bool) (string, error) {
	return "", nil
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
