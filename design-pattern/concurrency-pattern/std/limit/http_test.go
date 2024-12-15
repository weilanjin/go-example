package limit

import (
	"github.com/go-redis/redis_rate/v10"
	"github.com/redis/go-redis/v9"
	"net/http"
	"strconv"
	"time"
)

var limiter *redis_rate.Limiter

func init() {
	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	limiter = redis_rate.NewLimiter(rdb)
}

func rateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		res, err := limiter.Allow(r.Context(), "token:123", redis_rate.PerSecond(100))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		h := w.Header()
		h.Set("RateLimit-Remaining", strconv.Itoa(res.Remaining))
		if res.Allowed == 0 {
			seconds := int(res.ResetAfter / time.Second)
			h.Set("RateLimit-Reset", strconv.Itoa(seconds))
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		// 获得到令牌
		next.ServeHTTP(w, r)
	})
}