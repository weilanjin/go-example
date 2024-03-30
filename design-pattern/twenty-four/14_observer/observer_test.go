package observer

import "testing"

func TestObserver(t *testing.T) {
	subject := NewSubject()
	reader := NewReader("reader1")
	reader2 := NewReader("reader2")
	reader3 := NewReader("reader3")
	subject.Attach(reader)
	subject.Attach(reader2)
	subject.Attach(reader3)

	subject.UpdateContext("observer mode")
}
