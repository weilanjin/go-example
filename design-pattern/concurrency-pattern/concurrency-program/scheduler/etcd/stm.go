package main

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
	"math/rand"
	"sync"
)

// Software Transactional Memory 软件事物内存

func STM(clt *clientv3.Client) {
	// 设置5个账户,每个账户100元, 共500元
	totalAccounts := 5
	for i := range totalAccounts {
		k := fmt.Sprintf("accts/%d", i)
		if _, err := clt.Put(context.TODO(), k, "100"); err != nil {
			panic(err)
		}
	}
	// 主要的事务逻辑
	exchange := func(stm concurrency.STM) error {
		from, to := rand.Intn(totalAccounts), rand.Intn(totalAccounts)
		if from == to {
			return nil // 自己不能转账给自己
		}

		// 读取帐号值
		fromK, toK := fmt.Sprintf("accts/%d", from), fmt.Sprintf("accts/%d", to)
		fromBalance, toBalance := stm.Get(fromK), stm.Get(toK)
		fromInt, toInt := 0, 0
		fmt.Sscanf(fromBalance, "%d", &fromInt)
		fmt.Sscanf(toBalance, "%d", &toInt)

		// 模拟转账
		after := fromInt / 2
		fromInt, toInt = fromInt-after, toInt+after

		// 写入帐号值
		stm.Put(fromK, fmt.Sprintf("%d", fromInt))
		stm.Put(toK, fmt.Sprintf("%d", toInt))
		return nil
	}
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := range 100 {
				if _, err := concurrency.NewSTM(clt, exchange); err != nil {
					fmt.Println(j, "error:", err)
				}
			}
		}()
	}
	wg.Wait()

	// 检查帐号最后的值
	sum := 0
	accts, err := clt.Get(context.TODO(), "accts/", clientv3.WithPrefix())
	if err != nil {
		panic(err)
	}
	for _, kv := range accts.Kvs { // 遍历所有帐号
		fmt.Println(string(kv.Key), string(kv.Value))
		i := 0
		fmt.Sscanf(string(kv.Value), "%d", &i)
		sum += i
	}
	fmt.Println("sum:", sum)
}