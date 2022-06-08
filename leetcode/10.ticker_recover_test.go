package leetcode

import (
	"log"
	"testing"
	"time"
)

func proc() {
	panic("error")
}

func Test10(t *testing.T) {
	go func() {
		// 1 要求每秒钟调用一次proc函数
		// 2 要求程序不能退出
		ticker := time.NewTicker(time.Second * 1)
		for {
			<-ticker.C
			go func() {
				defer func() {
					if a := recover(); a != nil {
						log.Println(a)
					}
				}()
				proc()
			}()
		}
	}()
	time.Sleep(6 * time.Second)
	//select {}
}
