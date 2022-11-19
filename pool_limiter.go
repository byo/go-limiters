//go:build go1.18
// +build go1.18

package limiters

import (
	"context"
)

type poolLimiter[T any] struct {
	ch chan T
}

func NewPoolLimiter[T any](poolSize int, gen func() T) PoolLimiter[T] {

	ret := &poolLimiter[T]{
		ch: make(chan T, poolSize),
	}

	// Preallocate the pool of objects
	for i := 0; i < poolSize; i++ {
		ret.ch <- gen()
	}

	return ret
}

func (p *poolLimiter[T]) Acquire() T {
	return <-p.ch
}

func (p *poolLimiter[T]) Release(item T) {
	p.ch <- item
}

func (p *poolLimiter[T]) AcquireNoWait() (T, error) {
	var item T
	select {
	case item = <-p.ch:
		return item, nil
	default:
		return item, ErrResourceExhausted
	}
}

func (p *poolLimiter[T]) AcquireCtx(ctx context.Context) (T, error) {
	var item T

	if ctx.Err() == nil {
		select {
		case item = <-p.ch:
			return item, nil
		case <-ctx.Done():
			break
		}
	}

	return item, ctx.Err()
}
