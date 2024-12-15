package limit

import (
	"context"
	"golang.org/x/time/rate"
	"log"
	"testing"
	"time"
)

/*
	令牌桶的处理方式:
	- 假设用户配置的处理速率为r, 则每隔1/r每秒就会将一个令牌加入令牌桶中.
	- 假设令牌桶最多可以存放n个令牌, 如果新令牌到达时令牌桶已经满,那么这个
      新令牌会被丢弃.
	- 当处理一个请求时,就从令牌桶中删除一个令牌.
*/
// golang.org/x/time/rate
// Allow、Reserve、Wait 都会消耗一个令牌
// func (lim *Limiter) Allow() bool - 如果没有令牌可用, 返回false, 不会被阻塞
// func (lim *Limiter) Reserve() *Reservation - 如果没有令牌可用,这个方法将返回未来可用令牌,和调用者等待的时间
// func (lim *Limiter) Wait(ctx context.Context) (err error) - 如果没有令牌可用, 这个方法将会被阻塞,直到有令牌可用, 或者上下文被取消

func TestLimit(t *testing.T) {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)

	var limit rate.Limit = 5                   // 第一种: 每秒产生5个令牌
	limit = rate.Every(200 * time.Millisecond) // 第二种: 每200毫秒产生一个令牌

	limiter := rate.NewLimiter(limit, 3) // 令牌桶容量3

	for i := range 15 {
		log.Printf("got #%d, err:%v", i, limiter.Wait(context.Background()))
	}
}