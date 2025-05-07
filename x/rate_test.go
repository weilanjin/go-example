package x_test

import (
	"context"
	"log"
	"os"
	"sync"
	"testing"

	"golang.org/x/time/rate"
)

type ApiConnection struct {
	rate *rate.Limiter
}

func Open() *ApiConnection {
	return &ApiConnection{
		rate: rate.NewLimiter(rate.Limit(1), 1),
	}
}

func (api *ApiConnection) ReadFile(ctx context.Context) error {
	if err := api.rate.Wait(ctx); err != nil {
		return err
	}

	// pretend we do work here
	return nil
}

func (a *ApiConnection) ResolveAddress(ctx context.Context) error {
	if err := a.rate.Wait(ctx); err != nil {
		return err
	}

	// pretend we do work here
	return nil
}

func TestRate(t *testing.T) {
	defer log.Printf("Done.")
	log.SetOutput(os.Stdout)

	apiConn := Open()
	var wg sync.WaitGroup
	wg.Add(20)
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			if err := apiConn.ReadFile(context.Background()); err != nil {
				log.Printf("cannot ReadFIle: %v\n", err)
			}
			log.Println("ReadFile")
		}()
	}
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			if err := apiConn.ResolveAddress(context.Background()); err != nil {
				log.Printf("cannot ResolveAddress: %v\n", err)
			}
			log.Println("ResolveAddress")
		}()
	}
	wg.Wait()
}
