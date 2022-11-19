package limiters

import "context"

type tokenLimiterChannelEmpty struct {
	ch chan struct{}
}

func NewTokenLimiterChannelEmpty(maxTokens int) CancellableTokenLimiter {
	return &tokenLimiterChannelEmpty{
		ch: make(chan struct{}, maxTokens),
	}
}

func (c *tokenLimiterChannelEmpty) Acquire() {
	c.ch <- struct{}{}
}

func (c *tokenLimiterChannelEmpty) AcquireNoWait() error {
	select {
	case c.ch <- struct{}{}:
		return nil
	default:
		return ErrResourceExhausted
	}
}

func (c *tokenLimiterChannelEmpty) AcquireCtx(ctx context.Context) error {
	if ctx.Err() == nil {
		select {
		case c.ch <- struct{}{}:
			return nil
		case <-ctx.Done():
			break
		}
	}
	return ctx.Err()
}

func (c *tokenLimiterChannelEmpty) Release() {
	<-c.ch
}
