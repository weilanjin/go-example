package main

import (
	"bufio"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
	recipe "go.etcd.io/etcd/client/v3/experimental/recipes"
	"os"
)

// 分布式计数屏障
// 两段屏障 初始化需要 count 值
// Enter 阻塞,直到有count个节点调用了Enter方法,这些阻塞才能继续放行
// Leave 当一个节点调用Leave时,会被阻塞,直到有count个节点都调用Leave方法,这些阻塞才能继续放行

func DoubleBarrier(cli *clientv3.Client) {
	session, err := concurrency.NewSession(cli)
	if err != nil {
		panic(err)
	}
	b := recipe.NewDoubleBarrier(session, "double_barrier", 5)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		switch scanner.Text() {
		case "hold":
			b.Enter()
			fmt.Println("barrier enter")
		case "release":
			b.Leave()
			fmt.Println("barrier leave")
		case "quit", "exit":
			return
		default:
			fmt.Println("unknown action")
		}
	}
}