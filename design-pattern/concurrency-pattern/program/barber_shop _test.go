package program

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"testing"
)

// 理发师问题
// 计算机科学家 Edsger Dijkstra 在 1965 提出的, 在Silberschatz、Galvin 和 Gagne 的
// Operating System Concepts 有此问题的变种
//
// 一个理发店有1个理发师💇和几个座位
// - 如果没有顾客, 这个理发师就躺在理发椅上睡觉
// - 顾客必须唤醒理发师,让他开始理发
// - 如果有一位顾客到来, 理发师正在理发:
// 		* 如果还有空闲的座位, 则此顾客坐下
//		* 如果座位都做满了,则此顾客离开
// - 理发师理完发后,需要检查是否有等待的顾客
// 		* 如果有, 则请一位顾客起来开始理发
//		* 如果没有,理发师则去睡觉

// 一个队列, 有多个并发写(Multiple Writer, 顾客) 和 一个并发读(Single Reader, 理发师)
// 有一位顾客到来,则座位数+1,
// 理发师叫起一位等待的顾客开始理发,则座位数-1

type Semaphore chan struct{}

func (s Semaphore) Acquire() {
	s <- struct{}{}
}

func (s Semaphore) TryAcquire() bool {
	select {
	case s <- struct{}{}: // 还有空闲的座位
		return true
	default: // 没有空闲的座位了, 顾客离开
		return false
	}
}

func (s Semaphore) Release() {
	<-s
}

// Go 官方扩展库 semaphore.Weighted 同步原语,在Acquire()之前调用Release()会 panic

var seats = make(Semaphore, 3)

// 有多个理发师

func barber(name string) { // 理发师
	for {
		log.Println(name + " 老师尝试请求一位顾客")
		seats.Release()
		log.Println(name + " 老师叫起一位顾客, 开始理发")
		randomPause(2000)
	}
}

func customers() { // 模拟顾客陆续到来
	for {
		randomPause(1000)
		go customer()
	}
}

func customer() { // 顾客
	if seats.TryAcquire() {
		log.Println("customer come in")
	} else {
		log.Println("customer leave")
	}
}

func TestBarber(t *testing.T) {
	go barber("Tony")
	go barber("Jerry")
	go barber("Lucy")
	go customers()
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
}