package project

import (
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"testing"
)

// 多线程问题
// - 线程可能在不同的CPU间调度,不能充分利用CPU的亲和性(Affinity)
// - 避免数据竞争,线程需要使用同步原语, 降低了程序的性能
//
// 优化
// 在单CPU上运行反而不是低级设计, 程序没有数据竞争,减少了同步的性能损坏,
// 通过每个进程绑定一个CPU核,又能利用CPU的亲和性 (Redis)

// 一个程序启动多个进程,后再启动指定数量的子进程,子进程共享所监听的Socket(文件)

func TestProcess(t *testing.T) {
	// 单进程模式, 简单启动一个TCP Server 即可
	ln, err := net.Listen("tcp", ":8972")
	if err != nil {
		panic(err)
	}
	start(ln) // 处理 net.Listener
}

func TestMultiProcess(t *testing.T) {
	// 主进程
	// 先启动一个 TCP Server
	addr, err := net.ResolveTCPAddr("tcp", ":8972")
	if err != nil {
		panic(err)
	}
	ln, err := net.ListenTCP("tcp", addr)
	if err != nil {
		panic(err)
	}
	f, err := ln.File() // 等到句柄
	if err != nil {
		panic(err)
	}

	// 启动指定数量的子进程
	children := make([]*exec.Cmd, 10)
	for i := range children {
		children[i] = exec.Command(os.Args[0], "-prefork", "-child")
		children[i].Stdout = os.Stdout
		children[i].Stderr = os.Stderr
		children[i].ExtraFiles = []*os.File{f} // 把主进程监听的Socket传给子进程, 它们共享一个socket
		if err := children[i].Start(); err != nil {
			panic(err)
		}
	}
	for _, child := range children {
		if err := child.Wait(); err != nil {
			log.Printf("failed to start chid's starting: %v", err)
		}
	}
	os.Exit(0)
}

func start(ln net.Listener) {
	log.Println("started")
	for {
		conn, e := ln.Accept()
		if e != nil {
			log.Println("accept failed, err:", e)
			continue
		}
		go io.Copy(conn, conn) // 实现echo协议,将收到的数据原样返回
	}
}