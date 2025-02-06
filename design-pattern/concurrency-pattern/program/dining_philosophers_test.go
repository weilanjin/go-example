package program

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/fatih/color"
)

// 哲学家就餐问题
// 1971年计算机科学家Edsger Dijkstra 提出了一个同步原语问题
// [假设有五台计算机都试图访问无份共享的磁带驱动器]
//
// 就餐问题:
// 冥想--饿了--吃饭--冥想
// 假设有五位哲学家围坐在一张圆形餐桌旁,餐桌有无尽的可口的饭菜🍚,但是只有根筷子🥢,每根筷子
// 都位于两位哲学家之间.哲学家吃饭时,必须拿起自己左右两边的两根筷子,吃完饭后再放回筷子, 这
// 样其他哲学家也可以拿起筷子吃饭了.
/*
	形成死锁的四个条件
	- 禁止占用(No Preemption): 系统资源不能被强制地从一个线程中退出
	- 持有和等待(Hold and Wait): 一个线程在等待时持有并发资源 (持有并发资源的线程还在等待其他资源)
	- 互斥(Mutual Exclusion): 资源在同一时刻只能被分配给一个线程 (资源具有排他性)
	- 循环等待(Circular Waiting): 一系列线程相互持有其他线程所需要的资源.(线程之间必须有一个循环依赖的关系)
*/

// Chopstick 代表筷子
type Chopstick struct {
	sync.Mutex
}

// Philosopher 代表哲学家
type Philosopher struct {
	name           string // 哲学家名字
	leftChopstick  *Chopstick
	rightChopstick *Chopstick
	status         string // 冥想、饿了、吃饭、持有一根筷子并请求另一个筷子
}

func (p *Philosopher) dine() {
	for {
		mark(p, "冥想")
		randomPause(10)
		mark(p, "饿了")
		p.leftChopstick.Lock() // 先尝试拿起左手边的筷子
		randomPause(100)
		p.rightChopstick.Lock() // 再尝试拿起右手边的筷子
		mark(p, "吃饭")
		randomPause(10)
		p.rightChopstick.Unlock() // 先尝试放下右手边筷子
		p.leftChopstick.Unlock()  // 再尝试放下左手边筷子
	}
}

// 随机暂停一段时间
func randomPause(max int) {
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(max)))
}

// 显示此哲学家的状态
func mark(p *Philosopher, action string) {
	fmt.Printf("%s 开始 %s\n", p.name, action)
	p.status = fmt.Sprintf("%s 开始 %s\n", p.name, action)
}

func TestDpV1(t *testing.T) {
	go http.ListenAndServe(":8080", nil)
	// 哲学家的数量
	count := 5
	chopsticks := make([]*Chopstick, count)
	for i := 0; i < count; i++ {
		chopsticks[i] = &Chopstick{}
	}
	names := []string{
		color.RedString("哲学家1"),
		color.MagentaString("哲学家2"),
		color.CyanString("哲学家3"),
		color.GreenString("哲学家4"),
		color.WhiteString("哲学家5"),
	}
	philosophers := make([]*Philosopher, count) // 创建哲学家,给他们分配左右两边的筷子
	for i := 0; i < count; i++ {
		philosophers[i] = &Philosopher{
			name:           names[i],
			leftChopstick:  chopsticks[i],
			rightChopstick: chopsticks[(i+1)%count],
		}
		go philosophers[i].dine()

		// 解法1 最后一位哲学家不参与比赛, 避免死锁
		// if i < count-1 {
		// 	go philosophers[i].dine()
		// }
	}
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	fmt.Println("退出中... 每位哲学家的状态:")
	for _, p := range philosophers {
		fmt.Println(p.status)
	}
}
