package limiters

import "context"

type tokenLimiterChannelFull struct {
	ch chan struct{}
}

func NewTokenLimiterChanelFull(maxTokens int) CancellableTokenLimiter {
	ret := &tokenLimiterChannelFull{
		ch: make(chan struct{}, maxTokens),
	}

	for i := 0; i < maxTokens; i++ {
		ret.ch <- struct{}{}
	}

	return ret
}

func (c *tokenLimiterChannelFull) Acquire() ReleaseFunc {
	<-c.ch
	return func() { c.ch <- struct{}{} }
}

func (c *tokenLimiterChannelFull) AcquireNoWait() (ReleaseFunc, error) {
	select {
	case <-c.ch:
		return func() { c.ch <- struct{}{} }, nil
	default:
		return nil, ErrResourceExhausted
	}
}

func (c *tokenLimiterChannelFull) AcquireCtx(ctx context.Context) (ReleaseFunc, error) {
	if ctx.Err() == nil {
		select {
		case <-c.ch:
			return func() { c.ch <- struct{}{} }, nil
		case <-ctx.Done():
			break
		}
	}
	return nil, ctx.Err()
}
