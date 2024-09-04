package database

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

// 文件和数据库的区别
// 文件 不好维护元数据

// 问题
// 1.文件更新前截断, 如果并发读怎么办?
// 2.将数据写到文件可能不是原子的,取决于写入的大小,并发读可能会得不到完整的数据.
// 3.数据实际如何持久化到磁盘? 数据可能仍在写系统调用返回操作系统的页缓存.当系统崩溃并重启,文件的状态是什么?
func SaveData1(path string, data []byte) error {
	fp, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0664)
	if err != nil {
		return err
	}
	defer fp.Close()
	_, err = fp.Write(data)
	return err
}

// atomic renaming
// 1.将数据写到临时文件.
// 2.将临时文件重命名为目标文件.(原子操作)
//
// 如果系统崩溃前重命名, 不要影响源文件
// 同时读取文件没有问题
//
// 问题:
// 1.无法控制何时同步到磁盘, 元数据可能已经同步打磁盘上了, 系统正好崩溃. 造成文件损坏.
func SaveDate2(path string, data []byte) error {
	tmp := fmt.Sprintf("%s.tmp.%d", path, randomInt())
	fp, err := os.OpenFile(tmp, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0664)
	if err != nil {
		return err
	}
	defer fp.Close()
	_, err = fp.Write(data)
	if err != nil {
		os.Remove(tmp)
		return err
	}
	// os.Rename 是原子操作(Unix)
	return os.Rename(tmp, path)
}

func randomInt() int {
	rand.NewSource(time.Now().UnixNano())
	return rand.Intn(100)
}

// fsync
func SaveData3(path string, data []byte) error {
	tmp := fmt.Sprintf("%s.tmp.%d", path, randomInt())
	fp, err := os.OpenFile(tmp, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0664)
	if err != nil {
		return err
	}
	defer fp.Close()
	if _, err = fp.Write(data); err != nil {
		os.Remove(tmp)
		return err
	}
	if err := fp.Sync(); err != nil { // 确保数据已经写入磁盘
		os.Remove(tmp)
		return err
	}
	return os.Rename(tmp, path)
}
