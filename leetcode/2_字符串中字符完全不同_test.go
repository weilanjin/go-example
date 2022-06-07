package leetcode

import (
	"log"
	"strings"
	"testing"
)

// 判断字符串中字符是否完全不同
func Test2(t *testing.T) {
	str := "xnzf"
	// strings.Count 来判断在⼀个字符串中包含 的另外⼀个字符串的数量
	unique := isUniqueString1(str)
	log.Println(unique)

	// strings.Index
	unique2 := isUniqueString2(str)
	log.Println(unique2)
}

func isUniqueString1(s string) bool {
	if len(s) > 3000 {
		return false
	}
	for _, b := range s {
		if b > 127 {
			return false
		}
		if strings.Count(s, string(b)) > 1 {
			return false
		}
	}
	return true
}

func isUniqueString2(s string) bool {
	if len(s) > 3000 {
		return false
	}
	for k, b := range s {
		if b > 127 {
			return false
		}
		if strings.Index(s, string(b)) != k {
			return false
		}
	}
	return true
}
