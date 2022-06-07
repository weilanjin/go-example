package leetcode

import (
	"log"
	"strings"
	"testing"
)

func Test5(t *testing.T) {
	s := "How's the weather today."
	s1 := replaceSpace(s)
	log.Println(s1)
}

// " " -> 20%
// strings.Replace(s, "20%", -1) // -1 替换所有
func replaceSpace(s string) string {
	var sb strings.Builder
	for _, v := range s {
		if string(v) == " " {
			sb.WriteString("20%")
		} else {
			sb.WriteString(string(v))
		}
	}
	return sb.String()
}
