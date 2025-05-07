package database

import "os"

func LogCreate(path string) (*os.File, error) {
	return os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0664)
}

// LogAppend appends a line to the file.
// append-only is that is does not modify the existing data.
/*
	缺点：
		不好搜索查询
*/
func LogAppend(fp *os.File, line string) error {
	buf := []byte(line)
	buf = append(buf, '\n')
	_, err := fp.Write(buf)
	if err != nil {
		return err
	}
	return fp.Sync()
}