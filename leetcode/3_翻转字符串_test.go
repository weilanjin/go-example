package leetcode

import (
	"log"
	"testing"
)

func Test3(t *testing.T) {
	str := "I like Go programing language."
	r := []rune(str)
	sLen := len(str)
	for i := 0; i < sLen/2; i++ {
		r[i], r[sLen-i-1] = r[sLen-i-1], r[i] // 翻转字符串
	}
	log.Println(string(r))
}
