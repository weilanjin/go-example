package channel

import "unsafe"

type g struct {
	schedlink guintptr
}

type guintptr uintptr

func (p *guintptr) set(gp *g) {
	*p = guintptr(unsafe.Pointer(gp))
}

func (p *guintptr) ptr() *g {
	return (*g)(unsafe.Pointer(*p))
}

type gList struct {
	head guintptr
}

func (l *gList) empty() bool {
	return l.head == 0
}

func (l *gList) push(gp *g) {
	gp.schedlink = l.head
	l.head.set(gp)
}

func (l *gList) pop() *g {
	gp := l.head.ptr()
	if gp != nil {
		l.head = gp.schedlink
	}
	return gp
}