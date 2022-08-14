package limiters

import (
	"context"
	"errors"
)

var (
	ErrResourceExhausted = errors.New("resource exhausted")
)

type ReleaseFunc func()

type TokenLimiter interface {
	Acquire() ReleaseFunc
	AcquireNoWait() (ReleaseFunc, error)
}

type CancellableTokenLimiter interface {
	TokenLimiter
	AcquireCtx(ctx context.Context) (ReleaseFunc, error)
}
