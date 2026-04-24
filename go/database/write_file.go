package database

import (
	"fmt"
	"math/rand"
	"os"
)

/*
	文件和数据库的区别
	文件 不好维护元数据
*/

// WriteFile writes data to a file
/*
	缺点：
		1.It truncates the file before updating it, what if the file needs to be read concurrently?
		2.Writing data to files may not be atomic, depending(取决于) on the size of the write。
		  Concurrent readers might(可能) get incomplete(完整的) data.
		3.数据什么时候落盘。写入系统调用，数据可能在操作系统的页面缓冲中，What's the state of the file when
 		  the system crashes and reboots?
*/
func WriteFile(path string, data []byte) error {
	fp, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0664)
	if err != nil {
		return err
	}
	defer fp.Close()
	_, err = fp.Write(data)
	return err
}

// WriteFileAtomic writes data to a file atomically
// 现将数据转存储到一个临时文件中
// 然后将临时文件重命名为最终文件
// the rename operation is atomic
func WriteFileAtomic(path string, data []byte) error {
	tmp := fmt.Sprintf("%s.tmp.%d", path, rand.Intn(100))
	fp, err := os.OpenFile(tmp, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0664)
	if err != nil {
		return err
	}
	defer fp.Close()

	if _, err = fp.Write(data); err != nil {
		os.Remove(tmp)
		return err
	}
	// 在重命名前将数据落盘，防止系统崩溃损坏文件。
	if err = fp.Sync(); err != nil {
		os.Remove(tmp)
		return err
	}
	// os.Rename 是原子操作(Unix)
	if err := os.Rename(tmp, path); err != nil {
		os.Remove(tmp)
		return err
	}
	return nil
}