package file_test

import (
	"compress/gzip"
	"io"
	"os"
	"testing"
)

// 压缩
func TestCompress(t *testing.T) {
	outputFile, err := os.Create("test.txt.gz")
	if err != nil {
		panic(err)
	}
	gzipWr := gzip.NewWriter(outputFile)
	defer gzipWr.Close()

	_, err = gzipWr.Write([]byte("Gophers rule!\n"))
	if err != nil {
		panic(err)
	}
}

// 解压
// 标准库支持：gzip、zlib、bz2、flate、lzw
func TestUncompresss(t *testing.T) {
	gzipFile, err := os.Open("test.txt.gz")
	if err != nil {
		panic(err)
	}
	gzipRd, err := gzip.NewReader(gzipFile)
	if err != nil {
		panic(err)
	}
	defer gzipRd.Close()
	outfileWr, err := os.Create("unzipped.txt")
	if err != nil {
		panic(err)
	}
	defer outfileWr.Close()
	if _, err := io.Copy(outfileWr, gzipRd); err != nil {
		panic(err)
	}
}
