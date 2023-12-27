package file_test

import (
	"fmt"
	"io/fs"
	"os"
	"testing"
	"time"
)

// 打开文件 Open、OpenFile
func TestOpen(t *testing.T) {
	f, err := os.Open("doc.go") // 只读方式打开 OpenFile(name, O_RDONLY, 0)
	if err != nil {
		panic(err)
	}
	_ = f.Close()
	
	// 最后一个参数是权限模式permission mode
	// 第二个是打开时的属性
	// os.O_RDONLY // 只读
	// os.O_WRONLY // 只写
	// os.O_RDWR // 读写
	// os.O_APPEND // 往文件中添建（Append）
	// os.O_CREATE // 如果文件不存在则先创建
	// os.O_TRUNC // 文件打开时裁剪文件
	// os.O_EXCL // 和O_CREATE一起使用，文件不能存在
	// os.O_SYNC // 以同步I/O的方式打开
	f, err = os.OpenFile("test.txt", os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	_ = f.Close()
}

// Create Empty File
func TestTouch(t *testing.T) {
	f, err := os.Create("test.txt") // OpenFile(name, O_RDWR|O_CREATE|O_TRUNC, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	fmt.Println(f)
}

// 裁剪一个文件到100个字节。
// 如果文件本来就少于100个字节，则文件中原始内容得以保留，剩余的字节以null字节填充。
// 如果文件本来超过100个字节，则超过的字节会被抛弃。
// 这样我们总是得到精确的100个字节的文件。
// 传入0则会清空文件。
func TestTruncate(t *testing.T) {
	if err := os.Truncate("test.txt", 100); err != nil {
		panic(err)
	}
}

// 获取文件信息
func TestGetFileInfo(t *testing.T) {
	fileInfo, err := os.Stat("test.txt")
	if err != nil {
		panic(err)
	}
	
	fmt.Println("File name:", fileInfo.Name())
	fmt.Println("Size in bytes:", fileInfo.Size())
	fmt.Println("Permissions:", fileInfo.Mode())
	fmt.Println("Last modified:", fileInfo.ModTime())
	fmt.Println("Is Directory: ", fileInfo.IsDir())
	fmt.Printf("System interface type: %T\n", fileInfo.Sys())
	fmt.Printf("System info: %+v\n\n", fileInfo.Sys())
	
	info := fs.FormatFileInfo(fileInfo)
	fmt.Println(info)
}

// 重命名和移动
func TestMv(t *testing.T) {
	oldpath := "test.txt"
	newpath := "../exec/test.txt"
	if err := os.Rename(oldpath, newpath); err != nil {
		panic(err)
	}
}

// 删除文件
func TestRm(t *testing.T) {
	if err := os.Remove("../exec/test.txt"); err != nil {
		panic(err)
	}
}

// 检查文件是否存在 IsNotExist、IsExist
func TestIsExist(t *testing.T) {
	fileInfo, err := os.Stat("test.txt")
	if err != nil {
		if os.IsNotExist(err) {
			panic(err)
		}
		fmt.Println(err)
	}
	fmt.Println(fileInfo)
}

// 检查文件的读写权限
func TestRDWR(t *testing.T) {
	// 测试文件的写权限
	f, err := os.OpenFile("test.txt", os.O_WRONLY, 0666)
	if err != nil {
		if os.IsPermission(err) {
			fmt.Println("Error: Write permission denied.")
		}
	}
	f.Close()
	
	// 测试文件的读权限 os.O_RDONLY
	f, err = os.OpenFile("test.txt", os.O_RDONLY, 0666)
	if err != nil {
		if os.IsPermission(err) {
			fmt.Println("Error: Read permission denied.")
		}
	}
	f.Close()
}

// 修改文件的读写执行权
func TestChmod(t *testing.T) {
	// r - 4
	// w - 2
	// x - 1
	// rwx -> 7 可读写执行、rw- -> 6 读写、r-x -> 5 读执行
	
	// user u group g other o
	// 777 -- u rwx | g rwx | o rwx
	// 755 -- u rwx | g r-x | o r-x
	// 711 -- u rwx | g --x | o --x
	// 700 -- u rwx | g --- | o ---
	// 644 -- u rw- | g r-- | o r--
	// 600 -- u rw- | g --- | o ---
	if err := os.Chmod("test.txt", 0000); err != nil {
		panic(err)
	}
}

// 改变文件所有者
func TestChown(t *testing.T) {
	if 	err := os.Chown("test.txt", os.Getuid(), os.Getgid()); err != nil {
		panic(err)
	}
}

// 改变时间戳
func TestChtimes(t *testing.T) {
	twoDaysFromNow := time.Now().Add(48 * time.Hour)
	lastAccessTime := twoDaysFromNow
	lastModifyTime := twoDaysFromNow
	if err := os.Chtimes("test.txt", lastAccessTime, lastModifyTime); err != nil {
		panic(err)
	}
}

func TestLink(t *testing.T) {
	// 创建一个硬链接
	// 创建后同一个文件内容会有两个文件名,改变一个文件的内容会影响另一个.
	// 删除和重命名不会影响另一个
	if err := os.Link("test.txt", "test_also.txt"); err != nil {
		panic(err)
	}
	// 创建一个软链接
	if err := os.Symlink("test.txt", "test_sym.txt"); err != nil {
		panic(err)
	}
	fileInfo, err := os.Lstat("test_sym.txt")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Link info: %+v\n", fileInfo)
	// 改变软连接的拥有者不会影响原文件
	if err = os.Lchown("test_sym.txt", os.Getuid(), os.Getgid()); err != nil {
		panic(err)
	}
}