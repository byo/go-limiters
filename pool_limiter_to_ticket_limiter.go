package limiters

import "context"

type poolLimiterToTokenLimiter struct {
	l PoolLimiter[struct{}]
}

func NewTokenLimiterFromPoolLimiter(maxTokens int) CancellableTokenLimiter {
	return &poolLimiterToTokenLimiter{
		l: NewPoolLimiter(maxTokens, func() struct{} { return struct{}{} }),
	}
}

func (l *poolLimiterToTokenLimiter) Acquire() {
	l.l.Acquire()
}

func (l *poolLimiterToTokenLimiter) AcquireNoWait() error {
	_, err := l.l.AcquireNoWait()
	return err
}

func (l *poolLimiterToTokenLimiter) AcquireCtx(ctx context.Context) error {
	_, err := l.l.AcquireCtx(ctx)
	return err
}

func (l *poolLimiterToTokenLimiter) Release() {
	l.l.Release(struct{}{})
}
