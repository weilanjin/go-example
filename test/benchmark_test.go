package test

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
)

const num = 10000

func BenchmarkStringSprintf(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var str string
		for i := 0; i < num; i++ {
			str = fmt.Sprintf("%s%d", str, i)
		}
	}
	b.StopTimer()
}

func BenchmarkStringJoin(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var str string
		for i := 0; i < num; i++ {
			str = str + strconv.Itoa(i)
		}
	}
	b.StopTimer()
}

func BenchmarkStringStringBuilder(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var sb strings.Builder
		for i := 0; i < num; i++ {
			sb.WriteString(strconv.Itoa(i))
		}
		_ = sb.String()
	}
	b.StopTimer()
}

// go test -bench=".*"
// goos: darwin
// goarch: arm64
// pkg: lovec.wlj/example/test
// BenchmarkStringSprintf-10                     69          17848097 ns/op
// BenchmarkStringJoin-10                        82          14657208 ns/op
// BenchmarkStringStringBuilder-10             6164            194327 ns/op
