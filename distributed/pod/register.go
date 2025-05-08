package pod

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

var ErrRegister = errors.New("register fail")

type Register struct {
	rdb               redis.UniversalClient
	heartbeatInterval time.Duration
}

func NewRegister(rdb redis.UniversalClient, heartbeatInterval time.Duration) *Register {
	return &Register{
		rdb:               rdb,
		heartbeatInterval: heartbeatInterval,
	}
}

func (r *Register) RegisterInstance(podInfo *Instance) error {
	ctx := context.Background()
	if err := r.register(ctx, podInfo.ID, podInfo); err != nil {
		return err
	}
	go r.heartbeat(podInfo)
	return nil
}

func (r *Register) heartbeat(podInfo *Instance) {
	// 定期更新服务心跳，防止被服务发现机制认为已下线
	ticker := time.NewTicker(r.heartbeatInterval)
	for range ticker.C {
		podInfo.LashBeat = time.Now()
		if err := r.register(context.Background(), podInfo.ID, podInfo); err != nil {
			slog.Error("register fail", "podInfo", podInfo, "err", err)
		}
		slog.Info("Service heartbeat updated.", "pod", podInfo.Name)
	}
}

func (r *Register) Sweep() {
	pods, err := r.ServiceList()
	if err != nil {
		slog.Error("sweep fail", "err", err)
		return
	}
	for _, pod := range pods {
		if time.Since(pod.LashBeat) > r.heartbeatInterval*2 {
			if err := r.rdb.HDel(context.Background(), r.Key(), pod.ID).Err(); err != nil {
				slog.Error("sweep fail", "err", err)
			}
		}
	}
}

func (r *Register) ServiceList() ([]*Instance, error) {
	ctx := context.Background()
	serviceMap, err := rdb.HGetAll(ctx, r.Key()).Result()
	if err != nil {
		return nil, err
	}
	pods := make([]*Instance, 0, len(serviceMap))
	for _, v := range serviceMap {
		pod := &Instance{}
		if err := pod.UnmarshalBinary([]byte(v)); err != nil {
			return nil, err
		}
		pods = append(pods, pod)
	}
	return pods, nil
}
func (r *Register) register(ctx context.Context, serviceName string, podInfo *Instance) error {
	if err := r.rdb.HSet(ctx, r.Key(), serviceName, podInfo).Err(); err != nil {
		return errors.Join(ErrRegister, err)
	}
	return nil
}

func (r *Register) Key() string {
	return registerSetKey
}
