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

func (c *tokenLimiterChannelFull) Acquire() {
	<-c.ch
}

func (c *tokenLimiterChannelFull) AcquireNoWait() error {
	select {
	case <-c.ch:
		return nil
	default:
		return ErrResourceExhausted
	}
}

func (c *tokenLimiterChannelFull) AcquireCtx(ctx context.Context) error {
	if ctx.Err() == nil {
		select {
		case <-c.ch:
			return nil
		case <-ctx.Done():
			break
		}
	}
	return ctx.Err()
}

func (c *tokenLimiterChannelFull) Release() {
	c.ch <- struct{}{}
}
