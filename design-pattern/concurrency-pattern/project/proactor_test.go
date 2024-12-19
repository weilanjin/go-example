package project

import (
	"github.com/xtaci/gaio"
	"log"
	"net"
	"testing"
)

// Proactor 模式是一种用于处理异步IO操作的设计模式, 允许事件驱动应用程序且复用并异步分发请求.
// 使用一组异步操作(如异步读取、写入等), 当操作完成时会触发一个事件通知, 应用程序可以在事件通知处理完成的I/O 操作结果.
//
// 几个组件构成
// - 异步操作: 用于处理异步IO操作(异步读取、写入等),异步操作会向操作系统发送IO请求,并立即返回,不会阻塞应用程序的执行
// - 事件处理器: 用于处理异步操作完成后的事件通知.
// - 事件驱动器: Event Demulitplexer 用于监听多个异步操作的完成事件, 并将事件通知传递给相应的事件处理器
//
// 可以高效的处理大量的并发IO操作, 并避免了IO操作时应用程序的阻塞
//
// 反应器模式: 使用一组同步IO操作来处理并发的IO请求
// Proactor 模式: 使用一组异步IO操作来处理并发的IO请求

func TestEchoServer(t *testing.T) {
	watcher, err := gaio.NewWatcher()
	if err != nil {
		t.Fatal(err)
	}
	defer watcher.Close()
	go echoServer(watcher)
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		t.Fatal(err)
	}
	log.Println("echo server listening on", ln.Addr())
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		log.Println("new client", conn.RemoteAddr())

		// submit the first async read IO request
		if err = watcher.Read(nil, conn, make([]byte, 1024)); err != nil {
			log.Println(err)
			continue
		}
	}
}

func echoServer(watcher *gaio.Watcher) {
	for {
		rs, err := watcher.WaitIO()
		if err != nil {
			log.Println(err)
			return
		}
		for _, r := range rs {
			switch r.Operation {
			case gaio.OpRead:
				if r.Error == nil { // 读完成事件
					watcher.Write(nil, r.Conn, r.Buffer[:r.Size])
				}
			case gaio.OpWrite: // 写完成事件
				if r.Error == nil {
					watcher.Read(nil, r.Conn, r.Buffer[:r.Size])
				}
			}
		}
	}
}