package utils

import (
	"fmt"
	"os"
)

func HasWritePermission(path string) (bool, error) {
	// 获取目录的文件信息
	info, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	// 检查文件是否为目录，并判断是否有写权限
	if info.IsDir() {
		mode := info.Mode()
		// 检查当前用户是否具有写权限
		return mode&0200 != 0, nil
	}
	return false, fmt.Errorf("%s is not a directory", path)
}
