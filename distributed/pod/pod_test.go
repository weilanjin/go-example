package pod

import (
	"testing"
	"time"

	"golang.org/x/exp/slog"
)

func TestRegister(t *testing.T) {
	// podInfo := &Instance{
	// 	ID:       uid.UUID(),
	// 	Address:  "127.0.0.1",
	// 	Name:     "push-service",
	// 	Port:     8080,
	// 	Version:  "v1",
	// 	LashBeat: time.Now(),
	// }

	r := NewRegister(rdb, time.Second*3)
	// r.Register(podInfo)

	for {
		list, err := r.ServiceList()
		if err != nil {
			slog.Error("get pod set", "err", err)
			continue
		}
		for _, v := range list {
			slog.Info("get pod set", "data", v)
		}
		time.Sleep(3 * time.Second)
	}
}
