package main

import (
	"bufio"
	"flag"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	recipe "go.etcd.io/etcd/client/v3/experimental/recipes"
	"os"
	"strings"
)

// 分布式队列

var queueName = flag.String("queue-name", "test-queue", "queue name")

func queue(clt *clientv3.Client) {
	q := recipe.NewQueue(clt, *queueName)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		items := strings.Split(scanner.Text(), ",")
		switch items[0] {
		case "push":
			if len(items) != 2 {
				fmt.Println("must set value to push")
				continue
			}
			q.Enqueue(items[1]) // 入队
		case "pop":
			v, err := q.Dequeue()
			if err != nil {
				panic(err)
			}
			fmt.Println(v) // 输出出队元素
		case "quit", "exit":
			return
		default:
			fmt.Println("unknown action")
		}
	}
}