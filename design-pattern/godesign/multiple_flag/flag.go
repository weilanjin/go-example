package multiple_flag

import (
	"sync/atomic"
)

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