package file_test

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"io"
	"os"
	"testing"
)

func TestHash(t *testing.T) {
	data, err := os.ReadFile("go_spec.html")
	if err != nil {
		panic(err)
	}
	// 计算Hash
	fmt.Printf("Md5: %x\n\n", md5.Sum(data))
	fmt.Printf("Sha1: %x\n\n", sha1.Sum(data))
	fmt.Printf("Sha256: %x\n\n", sha256.Sum256(data))
	fmt.Printf("Sha512: %x\n\n", sha512.Sum512(data))
}

func TestHash1(t *testing.T) {
	file, err := os.Open("go_spec.html")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	hasher := md5.New()
	_, err = io.Copy(hasher, file)
	if err != nil {
		panic(err)
	}
	// 传递 nil 作为参数，因为我们不通参数传递数据，而是通过writer接口。
	sum := hasher.Sum(nil)
	fmt.Printf("Md5 checksum: %x\n", sum)
}