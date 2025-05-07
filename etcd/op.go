package main

import (
	"context"
	"log/slog"

	"go.etcd.io/etcd/client/v3"
)

var kvClient clientv3.KV

func init() {
	kvClient = clientv3.NewKV(client)
}

func OpPut(key, value string) {
	op := clientv3.OpPut(key, value)
	do, err := kvClient.Do(context.Background(), op)
	if err != nil {
		slog.Error("op put", slog.Any("err", err))
	} else {
		slog.Info("op put", slog.Int64("revision", do.Put().Header.Revision))
	}
}

func OpGet(key string) {
	op := clientv3.OpGet(key)
	do, err := kvClient.Do(context.Background(), op)
	if err != nil {
		slog.Error("op get", slog.Any("err", err))
	} else {
		slog.Info("op get", slog.Any("kvs", do.Get().Kvs), slog.Int64("revision", do.Put().Header.Revision))
	}
}

func OpDelete(key string) {
	op := clientv3.OpDelete(key)
	do, err := kvClient.Do(context.Background(), op)
	if err != nil {
		slog.Error("op delete", slog.Any("err", err))
	} else {
		slog.Info("op delete", slog.Int64("revision", do.Put().Header.Revision))
	}
}