package main

import (
	"bufio"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	recipe "go.etcd.io/etcd/client/v3/experimental/recipes"
	"os"
	"strconv"
	"strings"
)

// 分布式优先队列
// priority 值越小优先级越高
func priorityQueue(clt *clientv3.Client) {
	q := recipe.NewPriorityQueue(clt, *queueName)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		items := strings.Split(scanner.Text(), ",")
		switch items[0] {
		case "push":
			if len(items) != 3 {
				fmt.Println("must set value to push")
				continue
			}
			priority, err := strconv.Atoi(items[2])
			if err != nil {
				fmt.Println("priority must be int")
				continue
			}
			q.Enqueue(items[1], uint16(priority)) // 入队
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