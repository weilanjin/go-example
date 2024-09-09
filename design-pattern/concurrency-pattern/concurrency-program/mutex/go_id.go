package mutex

import (
	"log/slog"
	"runtime"
	"strconv"
	"strings"
)

// 获取当前goroutine的id
// 获取 goroutine id 的库
// petermattis/goid // 对各个版本兼容比较好
// kortschak/goroutine // 更简洁

func GoID() int {
	var buf [64]byte
	n := runtime.Stack(buf[:], false) // 读取堆栈信息
	// 从堆栈信息中找到goroutine哪一行,把id解析出来
	// log.Printf("buf: %s", buf) goroutine 6 [running]:
	idField := strings.Fields(strings.TrimLeft(string(buf[:n]), "goroutine "))[0]
	id, err := strconv.Atoi(idField)
	if err != nil {
		slog.Error("get goroutine id failed", "err", err)
	}
	return id
}
