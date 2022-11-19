package limiters

import (
	"context"
	"errors"
)

var (
	ErrResourceExhausted = errors.New("resource exhausted")
)

type TokenLimiter interface {
	Acquire()
	AcquireNoWait() error
	Release()
}

type CancellableTokenLimiter interface {
	TokenLimiter
	AcquireCtx(ctx context.Context) error
}
