// https://redis.io/docs/manual/patterns/distributed-locks/
package redsync

import (
	"sync"
	"testing"
	"time"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/redigo"
	"github.com/weilanjin/go-example/microservice/redis/initialize"
)

func TestDistributedLock(t *testing.T) {

	// Create a pool with go-redis (or redigo) which is the pool redisync will
	// use while communicating with Redis. This can also be any pool that
	// implements the `redis.Pool` interface.
	rdb := initialize.Redis()
	pool := redigo.NewPool(rdb) // or, pool := redigo.NewPool(...)

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
