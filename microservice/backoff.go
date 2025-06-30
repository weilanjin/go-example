package microservice

import (
	"context"
	"fmt"
	"math/rand/v2"
	"time"
)

// 带抖动的指数回退

func Do(ctx context.Context, fn func(context.Context) error) error {
	const (
		maxAttempts = 10
		baseDelay   = 1 * time.Second
		maxDelay    = 60 * time.Second
	)
	delay := baseDelay

	var timer *time.Timer
	defer func() {
		if timer != nil {
			timer.Stop()
		}
	}()

	for range maxAttempts {
		if err := fn(ctx); err == nil { // request
			return nil
		}
		delay *= 2 // 两倍两倍的递增
		delay = min(delay, maxDelay)

		/*
			rand.Float64 [0.0, 1.0]
			*0.5 [0.0, 0.5] 压缩一半
			-0.25 【-0.25， 0.25】 整个范围向左 <- 平移 0.25
		*/
		jitter := mulitplyDuration(delay, rand.Float64()*0.5-0.25) // ±25%
		sleepTime := delay + jitter

		timer = time.NewTimer(sleepTime)
		select {
		case <-ctx.Done():
			timer.Stop()
			return ctx.Err()
		case <-timer.C:
		}
	}
	return fmt.Errorf("failed after %d attempts", maxAttempts)
}

func mulitplyDuration(d time.Duration, mul float64) time.Duration {
	return time.Duration(float64(d) * mul)
}
