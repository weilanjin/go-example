package redsync

import (
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	goredislib "github.com/redis/go-redis/v9"
	"sync"
	"testing"
	"time"
)

func TestDistributedLock(t *testing.T) {

	// Create a pool with go-redis (or redigo) which is the pool redisync will
	// use while communicating with Redis. This can also be any pool that
	// implements the `redis.Pool` interface.
	client := goredislib.NewClient(&goredislib.Options{
		Addr: "localhost:6379",
	})
	pool := goredis.NewPool(client) // or, pool := redigo.NewPool(...)

	// Create an instance of redisync to be used to obtain a mutual exclusion
	// lock.
	rs := redsync.New(pool)

	// Obtain a new mutex by using the same name for all instances wanting the
	// same lock.
	mutexname := "my-global-mutex"
	mutex := rs.NewMutex(mutexname)

	var wg sync.WaitGroup
	gNum := 2
	wg.Add(gNum)
	for i := 0; i < gNum; i++ {
		go func() {
			defer wg.Done()
			t.Log("start get lock")
			if err := mutex.Lock(); err != nil {
				panic(err)
			}

			t.Log("lock success")
			time.Sleep(time.Second * 3)
			t.Log("start unlock")

			if ok, err := mutex.Unlock(); !ok || err != nil {
				t.Log(err)
				panic("unlock failed")
			}
			t.Log("unlock success")
		}()
	}
	wg.Wait()
}