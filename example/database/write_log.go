package database

import "os"

// Append-Only Logs

// 1. 一个数据库需要附加索引“indexes”, 来高效的查询数据
// 2. 如何删除数据.

func LogCreate(path string) (*os.File, error) {
	return os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0664)
}

// 仅追加数据, 不会修改原有的数据. 不需要重名操作
func LogAppend(fp *os.File, line string) error {
	if _, err := fp.WriteString(line + "\n"); err != nil {
		return err
	}
	return fp.Sync()
}
