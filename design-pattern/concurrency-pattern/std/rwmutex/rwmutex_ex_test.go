package rwmutex

import (
	"log"
	"testing"
	"time"
)

// 1.创建 100 个 reader
// 2.创建 50 个 writer
// 3.再创建 50 个 reader

// 2024/09/10 23:06:27 readers:  150           // 一共启动了 150 个 reader
// 2024/09/10 23:06:27 departing readers:  100 // 前 100 个 reader 是 departing 状态
// 2024/09/10 23:06:27 writer:  50
func TestRWMutexEx(t *testing.T) {
	var rwmutex RWMutexEx

	for i := 0; i < 100; i++ {
		go func() {
			rwmutex.RLock()
			time.Sleep(time.Hour)
			rwmutex.RUnlock()
		}()
	}

	time.Sleep(time.Second)

	for i := 0; i < 50; i++ {
		go func() {
			rwmutex.Lock()
			time.Sleep(time.Hour)
			rwmutex.Unlock()
		}()
	}

	time.Sleep(time.Second)

	for i := 0; i < 50; i++ {
		go func() {
			rwmutex.RLock()
			time.Sleep(time.Hour)
			rwmutex.RUnlock()
		}()
	}
	time.Sleep(time.Second)

	log.Println("readers: ", rwmutex.ReaderCount())
	log.Println("departing readers: ", rwmutex.ReaderWait())
	log.Println("writer: ", rwmutex.WriterCount())
}
