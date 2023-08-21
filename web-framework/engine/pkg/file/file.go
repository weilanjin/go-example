package file

import (
	"fmt"
	"os"
)

// 文件是否存在
func Exists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func NotExist(src string) bool {
	_, err := os.Stat(src)
	return os.IsNotExist(err)
}

func IsNotExistMkDir(src string) error {
	if NotExist(src) {
		if err := MkDir(src); err != nil {
			return err
		}
	}
	return nil
}

func MkDir(src string) error {
	err := os.MkdirAll(src, os.ModeDir)
	return err
}

func Permission(src string) bool {
	_, err := os.Stat(src)
	return os.IsPermission(err)
}

func MustOpen(fileName, dir string) (*os.File, error) {
	if Permission(dir) {
		return nil, fmt.Errorf("permission denied dir: %s", dir)
	}

	if err := IsNotExistMkDir(dir); err != nil {
		return nil, fmt.Errorf("error during make dir %s, err:%s", dir, err)
	}

	f, err := os.OpenFile(dir+string(os.PathSeparator)+fileName, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, fmt.Errorf("fail to open file, err: %s", err)
	}
	return f, nil
}