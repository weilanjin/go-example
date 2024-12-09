package main

import (
	"flag"
	clientv3 "go.etcd.io/etcd/client/v3"
	"strings"
)

var (
	nodeID    = flag.Int("id", 0, "node id")
	addr      = flag.String("addr", "http://127.0.0.1:2379,http://127.0.0.1:2389,http://127.0.0.1:2399", "listen address")
	electName = flag.String("name", "test-elect", "election name")
	lockName  = flag.String("lock", "test-lock", "lock name")
)

func main() {
	flag.Parse()
	endpoints := strings.Split(*addr, ",")
	cli, err := clientv3.New(clientv3.Config{
		Endpoints: endpoints,
	})
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	// Leader(cli)
	// Locker(cli)
	// queue(cli)
	// priorityQueue(cli)
	// Barrier(cli)
	// DoubleBarrier(cli)
	STM(cli)
}