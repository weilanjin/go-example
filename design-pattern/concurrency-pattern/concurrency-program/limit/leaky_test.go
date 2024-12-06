package limit

import (
	"go.uber.org/ratelimit"
	"log"
	"testing"
	"time"
)

// 漏桶
// 主要目的是控制将数据注入网络的速率, 平衡网络上的突发流量, 突发流量可用被整形成稳定的流量
// 因为它可以保持一个常量的输出速率,所以可以用来进行限流,并且因为使用了buffer缓存,所以可以
// 平滑处理突发请求.
//
// 桶有一点的容量,底部有孔,并且以固定的速率处理请求(水以一定的速度流出).
// 调用者以随机的速率向漏桶中放入请求.

/*
	漏桶
	- 如果漏桶已满,那么新的请求会被丢弃.(漏桶溢出)
	- 如果流入速度总是 < 流出速度, 漏桶总是处于不满的状态, 则不会有请求被丢弃
	- 如果流入速度总是 > 流出速度, 漏桶总是处于满的状态, 则不会有请求被丢弃
	- 如果有突发请求, 漏桶有一定的缓存作用,那么缓存满了才会丢弃请求.所以,在一定情况下可以削峰填谷,平滑请求的处理
*/
// 漏桶算法
// type Limiter
//	  func New(rate int, opts ...Option) Limiter
// type Option
// 	  func Per(per time.Duration) Option 设置时间窗口,默认窗口是1s
// 	  func WithClock(clock Clock) Option 设置时钟, 方便测试
// 	  func WithSlack(slack int) Option 设置一个宽松的值, 允许限流器积累一定的令牌,允许一定大小的突发流量

func TestLeakyBucket(t *testing.T) {
	// ratelimit.New(100) // 100/s 每秒产生100个令牌
	// ratelimit.New(100, ratelimit.Per(time.Minute)) // 每分钟产生100个令牌
	rl := ratelimit.New(1, ratelimit.WithSlack(3)) // 每秒产生一个令牌, 允许3个令牌积累, 允许一定大小的突发流量
	for i := 0; i < 10; i++ {
		rl.Take()
		log.Printf("got #%d", i)
		if i == 3 {
			time.Sleep(5 * time.Second)
		}
	}
}