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

func (c *tokenLimiterChannelEmpty) Acquire() ReleaseFunc {
	c.ch <- struct{}{}
	return func() { <-c.ch }
}

func (c *tokenLimiterChannelEmpty) AcquireNoWait() (ReleaseFunc, error) {
	select {
	case c.ch <- struct{}{}:
		return func() { <-c.ch }, nil
	default:
		return nil, ErrResourceExhausted
	}
}

func (c *tokenLimiterChannelEmpty) AcquireCtx(ctx context.Context) (ReleaseFunc, error) {
	if ctx.Err() == nil {
		select {
		case c.ch <- struct{}{}:
			return func() { <-c.ch }, nil
		case <-ctx.Done():
			break
		}
	}
	return nil, ctx.Err()
}
