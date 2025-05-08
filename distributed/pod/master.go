package pod

import (
	"context"
	"time"

	"github.com/weilanjin/go-example/pkg/uid"

	"github.com/redis/go-redis/v9"
)

type Master struct {
	register *Register
	IsMaster bool
}

func NewMaster(rdb redis.UniversalClient, heartbeatInterval time.Duration) *Master {
	return &Master{
		register: NewRegister(rdb, heartbeatInterval),
	}
}

func (m *Master) Start(ctx context.Context) error {
	podInfo := &Instance{
		ID:       uid.UUID(),
		Address:  "127.0.0.1",
		Name:     "push-service",
		Port:     8080,
		Version:  "v1",
		LashBeat: time.Now(),
	}
	if err := m.register.RegisterInstance(podInfo); err != nil {
		return err
	}
	if m.IsMaster {
		m.register.Sweep()
	}
	return nil
}

func (m *Master) Stop(ctx context.Context) error {
	return nil
}
