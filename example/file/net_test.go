package file_test

import (
	"io"
	"net/http"
	"os"
	"testing"
)

// 下载网络资源
func TestDownload(t *testing.T) {
	newFile, err := os.Create("go_spec.html")
	if err != nil {
		panic(err)
	}
	defer newFile.Close()
	resp, err := http.Get("https://go.dev/ref/spec")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	_, err = io.Copy(newFile, resp.Body)
	if err != nil {
		panic(err)
	}
}

