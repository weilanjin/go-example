package singleflight

import (
	"context"
	"golang.org/x/sync/singleflight"
)

// v2支持范型

type GroupV2[T any] struct {
	sf singleflight.Group
}

type ResultV2[T any] struct {
	Val *T
	Err error
}

func (g *GroupV2[T]) Do(ctx context.Context, key string, fn func() (*T, error)) (*T, error) {
	v, err, _ := g.sf.Do(key, func() (any, error) {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			return fn()
		}
	})
	if err != nil {
		return new(T), err
	}
	return v.(*T), err
}

func (g *GroupV2[T]) DoChan(key string, fn func() (*T, error)) <-chan ResultV2[T] {
	ch := make(chan ResultV2[T], 1)
	go func() {
		defer close(ch)
		res := <-g.sf.DoChan(key, func() (any, error) {
			val, err := fn()
			return val, err
		})
		ch <- ResultV2[T]{
			Val: res.Val.(*T),
			Err: res.Err,
		}
	}()
	return ch
}

func (g *GroupV2[T]) Forget(key string) {
	g.sf.Forget(key)
}