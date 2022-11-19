//go:build go1.18
// +build go1.18

package limiters

import "context"

type PoolLimiter[T any] interface {
	// Acquire acquires a resource from the pool and blocks if the pool is empty until another
	// goroutine returns the resource back to the pool
	Acquire() T

	// AcquireNoWait acquires a resource from the pool, immediately returns an error if the pool is empty
	AcquireNoWait() (T, error)

	// AcquireCtx acquires resource from the pool and blocks if the pool is empty until either another
	// goroutine returns the resource back to the pool or the context is done
	AcquireCtx(ctx context.Context) (T, error)

	Release(t T)
}
