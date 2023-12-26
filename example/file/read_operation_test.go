package file_test

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"testing"
)

var f *os.File
func init() {
	var err error
	// 打开只读
	f, err = os.Open("basic_operation_test.go") // OpenFile(name, O_RDONLY, 0)
	if err != nil {
		panic(err)
	}
}

func TestRead(t *testing.T) {
	defer f.Close()
	buf := make([]byte, 2048)
	byteSize, err := f.Read(buf)
	if err != nil {
		panic(err)
	}
	
	buf = buf[:byteSize:byteSize] // Clip 缩容
	fmt.Printf("%v, %s, %d\n", buf, buf, byteSize)
}

// 通过文件对象读取
func TestReadAll(t *testing.T) {
	defer f.Close()
	data, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", data)
}

// 通过文件名读取
func TestReadFile(t *testing.T) {
	data, err := os.ReadFile("test.txt")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", data)
}

// 读取大文件（按行读取）
func TestBufioRead(t *testing.T) {
	defer f.Close()
	bufReader := bufio.NewReader(f) // defaultBufSize = 4096
	for {
		line, err := bufReader.ReadString('\n')
		fmt.Printf(line)
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
	}
	fmt.Println()
}

// 更专业 读取大文件
func TestScanerRead(t *testing.T) {
	defer f.Close()
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines) // bufio.ScanWords, 可以自定义分隔符 SplitFunc
	for {
		if !scanner.Scan() { // 循环获取下一个 token
			if err := scanner.Err(); err != nil && err != io.EOF {
				panic(err)
			}
			break
		}
		fmt.Println(scanner.Text())	// scanner.Bytes()
	}
}

func TestReadFull(t *testing.T) {
	defer f.Close()

	// io.ReadFull()在文件的字节数小于byte slice字节数的时候会返回错误
	buf := make([]byte, 3)
	n, err := io.ReadFull(f, buf)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s, %d", buf, n)
}

