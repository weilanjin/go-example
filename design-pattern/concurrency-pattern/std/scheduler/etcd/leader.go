package main

import (
	"bufio"
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
	"log"
	"os"
)

// Leader 选主

func Leader(cli *clientv3.Client) {
	// 创建一个并发的session
	session, err := concurrency.NewSession(cli)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	el := concurrency.NewElection(session, *electName)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		switch line {
		case "elect": // 启动选取
			go elect(el)
		case "proclaim": // 宣告,只是设置主节点的值
			proclaim(el)
		case "resign": // 放弃主
			resign(el)
		case "watch": // 监听主从变化的事件
			go watch(el)
		case "query": // 主动查询
			query(el)
		case "rev": // 查看版本号
			revision(el)
		case "exit", "quit":
			return
		default:
			fmt.Println("unknown action")
		}
	}
}

var count int

// elect 选主
// Campaign 这个是一个阻塞方法
// 取消阻塞条件
// - 成功选到主节点
// - 方法返回错误
// - ctx 被撤销
func elect(el *concurrency.Election) {
	log.Println("acampaigning for ID:", *nodeID)
	// 调用Campaign方法选主, 主节点的值为value-<主节点ID>-<count>
	if err := el.Campaign(context.Background(), fmt.Sprintf("value-%d-%d", *nodeID, count)); err != nil {
		log.Printf("failed to elect: %v", err)
		return
	}
	count++
	log.Println("campaigned for ID:", *nodeID)
}

// proclaim 为主节点设置新值, 但不会重新选主
func proclaim(el *concurrency.Election) {
	log.Println("proclaiming for ID:", *nodeID)
	// 调用Proclaim方法宣告,主节点的值为value-<主节点ID>-<count>
	if err := el.Proclaim(context.Background(), fmt.Sprintf("value-%d-%d", *nodeID, count)); err != nil {
		log.Printf("failed to proclaim: %v", err)
		return
	}
	count++
	log.Println("proclaimed for ID:", *nodeID)
}

// resign 重新选主, 有可能另一个节点被选为主节点
// 当前的主节点辞去作为主节点,开始重新一轮选举
func resign(el *concurrency.Election) {
	log.Println("resigning for ID:", *nodeID)
	// 调用Resign方法放弃主
	if err := el.Resign(context.Background()); err != nil {
		log.Printf("failed to resign: %v", err)
		return
	}
	log.Println("resigned for ID:", *nodeID)
}

// watch 监听主从变化的事件
func watch(el *concurrency.Election) {
	// 监听主从变化的事件
	// 显示主节点的变化信息, 不会返回主节点的全部历史变化信息,只会返回最近的一条变化信息以及之后·
	ch := el.Observe(context.Background())
	log.Println("start to watch for ID:", *nodeID)

	for range 10 {
		resp := <-ch
		log.Println("observed event:", string(resp.Kvs[0].Key), string(resp.Kvs[0].Value))
	}
}

// query 查询主节点的信息
func query(el *concurrency.Election) {
	resp, err := el.Leader(context.Background())
	if err != nil {
		log.Printf("failed to query leader: %v", err)
		return
	}
	log.Println("current leader:", string(resp.Kvs[0].Key), string(resp.Kvs[0].Value))
}

// revision 可以查询主节点的版本号
func revision(el *concurrency.Election) {
	rev := el.Rev()
	log.Println("current revision:", rev)
}