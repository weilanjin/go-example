package file_test

import (
	"bufio"
	"fmt"
	"os"
	"testing"
)

func TestOpenFile(t *testing.T) {
	// 覆盖写
	f, err := os.OpenFile("test.txt", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	f.WriteString("OK\n")
	f.Write([]byte("Bytes!\n"))
}

func TestOsWiterFile(t *testing.T) {
	// os.OpenFile("test.txt", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if err := os.WriteFile("test.txt", []byte("Hi\n"), 0666); err != nil {
		panic(err)
	}
}

func TestBufioWiter(t *testing.T) {
	// 打开只写文件
	f, err := os.OpenFile("test.txt", os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	bufWr := bufio.NewWriter(f) // 默认 4096 bufio.NewWriterSize(w io.Writer, size int) 可以指定buf大小
	bufWr.WriteString("你好")
	// 还有多少字节可用
	fmt.Printf("Bytes buffered: %d\n", bufWr.Buffered()) // 6
	fmt.Printf("Available buffer: %d\n", bufWr.Available()) // 4090
	bufWr.Flush()
}

