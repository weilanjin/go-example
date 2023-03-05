package cache_test

import (
	"regexp"
	"testing"
)

// 解析出单位
func TestParseUnit(t *testing.T) {
	size := "512GB"
	re, _ := regexp.Compile("^[0-9]+")

	loc := re.FindStringIndex(size)
	// unit := string(re.ReplaceAll([]byte(size), []byte("")))
	t.Log(size[:loc[1]], size[loc[1]:])
}
