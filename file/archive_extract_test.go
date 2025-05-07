package file_test

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"testing"
)

// 归档
func TestArchive(t *testing.T) {
	outFile, err := os.Create("test.zip") // 打包到什么文件
	if err != nil {
		panic(err)
	}
	defer outFile.Close()

	zipWr := zip.NewWriter(outFile)
	defer zipWr.Close()
	
	var filesToArchive = []struct{ // mock file
		Name, Body string
	}{
		{"test1.txt", "String contents of file"},
		{"test2.txt", "fleet"},
	}
	// 将要打包的内容写入到打包文件中
	for _, file := range filesToArchive {
		fw, err := zipWr.Create(file.Name)
		if err != nil {
			panic(err)
		}
		if _, err = fw.Write([]byte(file.Body)); err != nil {
			panic(err)
		}
	}
}

// 提取
func TestExtract(t *testing.T) {
	zipRd, err := zip.OpenReader("test.zip")
	if err != nil {
		panic(err)
	}
	defer zipRd.Close()
	for _, file := range zipRd.Reader.File {
		zippedFile, err := file.Open()
		if err != nil {
			panic(err)
		}
		defer zippedFile.Close()
		targetDir := "./"
		extractedFilePath := filepath.Join(targetDir, file.Name)
		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(extractedFilePath, file.Mode()); err != nil { // 创建文件夹并设置同样的权限
				panic(err)			
			}
		} else {
			outputFile, err := os.OpenFile(extractedFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
			if err != nil {
				panic(err)
			}
			defer outputFile.Close()
			if _, err := io.Copy(outputFile, zippedFile); err != nil {
				panic(err)
			}
		}
		zippedFile.Close()
	}
}