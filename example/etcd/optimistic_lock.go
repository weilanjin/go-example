package main

import (
	"context"
	"log/slog"
	"time"

	"go.etcd.io/etcd/clientv3"
)

// 乐观锁 分布式锁

// lease实现锁自动过期: 防止 etcd 宕机后重启不能自动释放.
// op 操作
// txn事务: if else then

func tx() {
	// 1.上锁(创建租约) 自动续租 拿着租约去抢占一个key
	lease := clientv3.NewLease(client)
	// 申请一个5s的租约
	leaseGrantResp, err := lease.Grant(context.Background(), 5)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	leaseID := leaseGrantResp.ID

	// 准备一个用于取消自动续租的context
	ctx, cancel := context.WithCancel(context.Background())
	// 确保函数退出,自动续租会停止
	defer cancel()
	// Revoke租约告诉leaseID直接释放
	defer lease.Revoke(context.TODO(), leaseID)

	// 5s 之后自动续约
	keepRespChan, err := lease.KeepAlive(ctx, leaseID)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	go func() {
		for {
			select {
			case keepResp := <-keepRespChan:
				if keepResp == nil {
					slog.Error("租约已经失效")
					goto END
				} else {
					slog.Info("自动续约:", slog.Int64("revision", keepResp.Revision))
				}
			}
		}
	END:
	}()
	// if 不存在key, then 设置它, else 抢锁失败
	kv := clientv3.NewKV(client)

	jobName := "/cron/lock/job1"
	// 创建事务
	txn := kv.Txn(context.TODO())
	txn.If(clientv3.Compare(clientv3.CreateRevision(jobName), "=", 0)).
		Then(clientv3.OpPut(jobName, "work", clientv3.WithLease(leaseID))).
		Else(clientv3.OpGet(jobName))
	if tr, err := txn.Commit(); err != nil {
		slog.Error("tx commit", slog.Any("err", err))
		return
	} else if !tr.Succeeded { // 是否抢到锁
		slog.Info("锁被占用", slog.Any("kvs", tr.Responses[0].GetResponseRange().Kvs[0].Value))
	}
	// 2.处理业务
	slog.Info("处理业务")
	time.Sleep(5 * time.Second)
	// 3.释放锁(取消自动续租) 释放租约
	// defer
}
