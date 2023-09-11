package multiple_flag

import (
	"sync/atomic"
)

// 使用了左移运算符 1 << iota，使得常量的值呈现 1、2、4、8... 这样的递增效果
// 为了位运算方便，通过对属性进行位运算，来决定输出内容，其本质上跟基于位运算的权限管理是一样的
const (
	Ldate = 1 << iota
	Ltime
	Lshortfile
	LstdFlags = Ldate | Ltime
)

type Logger struct {
	flag atomic.Int32 // properties
	// ...
}

func New(flag int) *Logger {
	l := new(Logger)
	l.SetFlags(flag)
	return l
}

func (l *Logger) SetFlags(flag int) {
	l.flag.Store(int32(flag))
}

func (l *Logger) Flags() int {
	return int(l.flag.Load())
}

func (l *Logger) Println(msg string) {
	if l.Flags()&(Ldate|Ltime) != 0 { // 包含 Ldate 或者 Ltime
		//...
	}
	if l.Flags()&(Lshortfile|LstdFlags) != 0 { // 包含 Lshortfile 或者 LstdFlags
		// ...
	}
}