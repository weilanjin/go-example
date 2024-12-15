package cond

import (
	"log"
	"math/rand"
	"sync"
	"testing"
	"time"
)

/*
短跑比赛
1.每个运动员(单独的一个子goroutine)就位后,都把变量 ready + 1

	并调用Broadcast方法(或者Signal方法,因为只有一个裁判员).

2.裁判员(主goroutine)检查条件, 如果条件不满足(ready != 10),

	则调用wait方法阻塞自己,直到条件满足(ready == 10), 才能继续执行,
	击响发令枪,宣布比赛开始.
*/
func TestMatch(t *testing.T) {
	c := sync.NewCond(&sync.Mutex{})
	var ready int
	for i := 0; i < 10; i++ {
		go func(i int) {
			time.Sleep(time.Duration(rand.Int63n(10)) * time.Second)
			c.L.Lock()
			ready++
			c.L.Unlock()

			log.Printf("运动员#%d已就绪\n", i)

			c.Broadcast() // 运动员 i 准备就绪, 通知裁判裁判员
		}(i)
	}
	c.L.Lock()
	for ready != 10 { // 检查条件是否满足
		c.Wait()
		log.Println("裁判员被唤醒一次")
	}
	c.L.Unlock()

	log.Println("比赛开始")
}
