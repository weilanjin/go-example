package pool

import (
	"log/slog"
	"net/http"
	"sync"
	"testing"
	"time"
)

func TestHttpPool(t *testing.T) {
	var p sync.Pool
	p.New = func() any {
		return &http.Client{
			Timeout: 5 * time.Second,
		}
	}
	var wg sync.WaitGroup
	wg.Add(10)
	go func() {
		for i := 0; i < 10; i++ {
			defer wg.Done()
			client := p.Get().(*http.Client)
			defer p.Put(client)
			resp, err := client.Get("http://www.baidu.com")
			if err != nil {
				t.Error(err)
			}
			slog.Info(resp.Status)
			resp.Body.Close()
		}
	}()
	wg.Wait()
}
