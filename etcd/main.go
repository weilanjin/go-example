package main

import (
	"context"
	"log/slog"
	"time"

	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var client *clientv3.Client

func init() {
	config := clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2380"}, // 集群列表
		DialTimeout: 5 * time.Second,
	}
	var err error
	if client, err = clientv3.New(config); err != nil {
		panic(err)
	}
}

const (
	Username = "username"
	k        = "greeting"
)

func main() {
	ctx := context.Background()

	// 申请一个lease(租约) 10s
	lease := clientv3.NewLease(client)
	lgr, err := lease.Grant(ctx, 5)
	if err != nil {
		slog.Error("lease grant", slog.Any("ttl", "10s"), slog.Any("err", err))
		return
	}

	leaseID := lgr.ID
	go autoLease(lease, leaseID)

	kvClient := clientv3.NewKV(client)

	// ResponseHeader 描述集群信息
	// Put 一个KV, 让它与租约关联起来, 从而实现10s后自动过期
	_, err = kvClient.Put(ctx, Username, "lanjin.wei", clientv3.WithLease(leaseID))
	if err != nil {
		slog.Error("etcdctl put kv", slog.String(Username, "lanjin.wei"), slog.Any("err", err))
	}

	time.Sleep(6 * time.Second)

	gr, err := kvClient.Get(ctx, Username)
	if err != nil {
		slog.Error("etcdctl get k", slog.Any("key", Username), slog.Any("err", err))
	} else if gr.Count == 0 { // count = 0 时 kv 过期
		slog.Info("kv已过期", slog.String("key", Username))
	}
	slog.Info("get kvs", slog.Any("kvs", gr.Kvs))
}

// 自动续租
func autoLease(lease clientv3.Lease, leaseID clientv3.LeaseID) {
	keepRespCh, err := lease.KeepAlive(context.Background(), leaseID)
	if err != nil {
		slog.Error("lease keep alive", slog.Any("err", err))
		return
	}

	isEnd := false
	for {
		select {
		case keepResp := <-keepRespCh:
			if keepResp != nil { // k 已经过期. etcd 掉线过久
				slog.Info("auto lease fail", slog.Int64("leaseID", int64(leaseID)))
				isEnd = true
			} else { // 每秒会续租一次
				slog.Info("auto lease success", slog.Int64("leaseID", int64(leaseID)))
			}
		}
		if isEnd {
			break
		}
	}
}

// watchKv 监听kv的修改或删除
func watchKv() {
	ctx := context.Background()
	gr, err := clientv3.NewKV(client).Get(ctx, Username)
	if err != nil {
		slog.Error("etcdctl get k", slog.Any("key", Username), slog.Any("err", err))
	} else if gr.Count == 0 { // count = 0 时 kv 过期
		slog.Info("kv已过期", slog.String("key", Username))
	}

	// 当前etcd集群事务ID, 单调递增
	watchStartRevs := gr.Header.Revision + 1
	// 创建一个 watcher
	watcher := clientv3.NewWatcher(client)
	// 启动监听
	slog.Info("从该版本开始监听", slog.Int64("revision", watchStartRevs))
	wcCh := watcher.Watch(ctx, Username, clientv3.WithRev(watchStartRevs))

	// 处理kv变化的事件
	for watchResp := range wcCh {
		for _, e := range watchResp.Events {
			switch e.Type {
			case mvccpb.PUT:
				slog.Info("update:", slog.Any("new kv", e.Kv), slog.Any("old kv", e.PrevKv))
			case mvccpb.DELETE:
				slog.Info("delete:", slog.Int64("Revision", e.Kv.ModRevision))
			}
		}
	}
}
