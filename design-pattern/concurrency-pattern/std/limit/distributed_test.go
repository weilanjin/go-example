package limit

import (
	"context"
	"github.com/go-redis/redis_rate/v10"
	"github.com/redis/go-redis/v9"
	"log"
	"sync"
	"testing"
	"time"
)

// func NewLimiter(rdb rediser) *Limiter // 并不是创建这个限流器时设定流的速率和容量, 而是在请求时传入限流的参数
// func (l Limiter) Allow(ctx context.Context, key string, limit Limit) (*Result, error)
// func (l Limiter) AllowAtMost(ctx context.Context, key string, limit Limit, n int) (*Result, error)
// func (l Limiter) AllowN(ctx context.Context, key string, limit Limit, n int) (*Result, error)
// func (l *Limiter) Reset(ctx context.Context, key string) error // 获取一个令牌,并重置所有的限制以及之前的请求统计

func TestRedis(t *testing.T) {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)

	var wg sync.WaitGroup
	wg.Add(2)

	for i := range 2 {
		go func(i int) {
			defer wg.Done()
			ctx := context.Background()
			rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
			limiter := redis_rate.NewLimiter(rdb)
			for range 10 {
				res, err := limiter.Allow(ctx, "token:123", redis_rate.PerSecond(5))
				if err != nil {
					log.Println(err)
				}
				log.Println(i, "allowed", res.Allowed, "remaining", res.Remaining, "retry after", res.RetryAfter)
				if res.Allowed == 0 {
					time.Sleep(res.RetryAfter)
				}
			}
		}(i)
	}
	wg.Wait()
}