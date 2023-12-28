package file_test

import (
	"fmt"
	"os"
	"testing"
)

func TestTempDir(t *testing.T) {
	// 在系统临时文件夹中创建一个临时文件夹 lance
	tempDirPath, err := os.MkdirTemp("", "lance")
	if err != nil {
		panic(err)
	}
	fmt.Println("Temp dir created: ",tempDirPath)
	// 在临时文件夹中创建临时文件
	tempFile, err := os.CreateTemp(tempDirPath, "test.txt")
	if err != nil {
		panic(err)
	}
	defer tempFile.Close()
	fmt.Println("Temp file created: ", tempFile.Name())
	
	// 操作完成手动移除
	if err = os.RemoveAll(tempDirPath); err != nil {
		panic(err)
	}
}

// output:
// Temp dir created:  /var/folders/ch/9dxh1q013dj93fgd2fdqtl300000gn/T/lance2382240141
// Temp file created:  /var/folders/ch/9dxh1q013dj93fgd2fdqtl300000gn/T/lance2382240141/test.txt207760967