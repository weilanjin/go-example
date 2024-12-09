package main

import (
	"bufio"
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
	recipe "go.etcd.io/etcd/client/v3/experimental/recipes"
	"log"
	"math/rand"
	"os"
	"time"
)

// 分布式锁
// 分布式读写锁

func Locker(clt *clientv3.Client) {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	// useLock(clt)
	useMutex(clt)
}

func useLock(clt *clientv3.Client) {
	session, err := concurrency.NewSession(clt, concurrency.WithTTL(10)) // 默认60s 持有锁的节点会自动释放锁
	if err != nil {
		panic(err)
	}
	defer session.Close()
	mux := concurrency.NewLocker(session, *lockName)

	log.Println("acquiring lock")
	mux.Lock()
	log.Println("acquired lock")

	time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
	mux.Unlock()

	log.Println("released lock")
}

func useMutex(clt *clientv3.Client) {
	session, err := concurrency.NewSession(clt, concurrency.WithTTL(10)) // 默认60s 持有锁的节点会自动释放锁
	if err != nil {
		panic(err)
	}
	defer session.Close()
	mux := concurrency.NewMutex(session, *lockName)

	// 在请求锁之前查询key
	log.Printf("before acquiring. key: %s\n", mux.Key())

	log.Println("acquiring lock")
	if err := mux.Lock(context.TODO()); err != nil {
		panic(err)
	}
	log.Println("acquired lock")

	time.Sleep(time.Duration(rand.Intn(10)) * time.Second)

	if err := mux.Unlock(context.TODO()); err != nil {
		panic(err)
	}
	log.Println("released lock")
}

// 读写锁
// 读写锁的等待顺序
// - 当写锁被持有时, 对读锁和写锁的请求会等待写锁的释放
// - 当读锁被持有时, 对写锁的请求会等待读锁的释放,对读锁的请求可以直接获得锁
// - 当读锁被持有时, 这时候如果有一个节点请求写锁, 则会等待前面的读锁释放;如果此时再有对读锁的请求,则会被阻塞,直到前面的写锁释放.(和标准库一样)
func useRWMutex(clt *clientv3.Client) {
	session, err := concurrency.NewSession(clt, concurrency.WithTTL(10)) // 默认60s 持有锁的节点会自动释放锁
	if err != nil {
		panic(err)
	}
	defer session.Close()
	rwMutex := recipe.NewRWMutex(session, "rw")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		action := scanner.Text()
		switch action {
		case "w":
			writeLocker(rwMutex)
		case "r":
			readLocker(rwMutex)
		default:
			fmt.Println("unknown action")
		}
	}
}

func readLocker(rwMutex *recipe.RWMutex) {
	log.Println("acquiring read lock")
	if err := rwMutex.RLock(); err != nil {
		panic(err)
	}
	log.Println("acquired read lock")

	time.Sleep(time.Duration(rand.Intn(10)) * time.Second)

	if err := rwMutex.RUnlock(); err != nil {
		panic(err)
	}

	log.Println("released read lock")
}

func writeLocker(rwMutex *recipe.RWMutex) {
	log.Println("acquiring write lock")
	if err := rwMutex.Lock(); err != nil {
		panic(err)
	}
	log.Println("acquired write lock")

	time.Sleep(time.Duration(rand.Intn(10)) * time.Second)

	if err := rwMutex.Unlock(); err != nil {
		panic(err)
	}
	log.Println("released write lock")
}