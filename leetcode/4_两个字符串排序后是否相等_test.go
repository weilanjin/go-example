package leetcode

import (
	"log"
	"strings"
	"testing"
)

func Test4(t *testing.T) {
	s1 := "wei lan jin"
	//s2 := "che yon lin"
	s2 := "ian jew lin"
	regroup := isRegroup(s1, s2)
	log.Println(regroup)
}

// 重新排序后判断两个字符串是否相等
func isRegroup(s1 string, s2 string) bool {
	sl1 := len(s1)
	sl2 := len(s2)
	if sl1 != sl2 {
		return false
	}
	for _, b := range s1 {
		if strings.Count(s1, string(b)) != strings.Count(s2, string(b)) {
			return false
		}
	}
	return true
}
