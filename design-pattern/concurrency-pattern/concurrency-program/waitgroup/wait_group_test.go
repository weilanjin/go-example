package waitgroup

import (
	"context"
	"sync"
	"testing"

	"golang.org/x/sync/errgroup"
)

func TestWG(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Done()
	wg.Wait()
}

func TestErrGroup(t *testing.T) {
	eg, _ := errgroup.WithContext(context.Background())
	eg.Go(func() error {
		return nil
	})
	eg.Go(func() error {
		return nil
	})
	eg.Wait()
}
