package file_test

import (
	"fmt"
	"io"
	"os"
	"testing"
)

// 拷贝文件
func TestCopy(t *testing.T) {
	f, err := os.Open("test.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	fBk, err := os.Create("test_bk.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if _, err := io.Copy(fBk, f); err != nil {
		panic(err)
	}
	if err = fBk.Sync(); err != nil {
		panic(err)
	}
}

// 跳转到文件指定位置
// f.Seek(0, 1) // 当前位置
// f.Seek(0, 0) // 开始位置
func TestSeek(t *testing.T) {
	f, err := os.Open("test.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	
	var offset int64 = 5 // 偏移多少byte
	var whence = 0 // 0 = 文件开始位置, 1 = 当前位置 2 = 文件结尾
	
	newPosition, err := f.Seek(offset, whence)
	if err != nil {
		panic(err)
	}
	fmt.Println("Just moved to 5:", newPosition)
	
	// 退回2byte
	ret, err := f.Seek(-2, 1)
	if err != nil {
		panic(err)
	}
	fmt.Println("Just moved back two:", ret)
	
	// 获取当前位置
	correntPosition, err := f.Seek(0, 1)
	if err != nil {
		panic(err)
	}
	fmt.Println("corrent position:", correntPosition)

	// 调转到文件开始处
	startPosition, err := f.Seek(0, 0)
	if err != nil {
		panic(err)
	}
	fmt.Println("position ofter seeking 0,0", startPosition)
}