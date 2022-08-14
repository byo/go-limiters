package limiters

import "sync"

type tokenLimiterCounterMutex struct {
	cnt int
	m   sync.Mutex
	c   sync.Cond
}

func NewTokenLimiterCounterMutex(maxTokens int) TokenLimiter {
	ret := &tokenLimiterCounterMutex{
		cnt: maxTokens,
	}
	ret.c.L = &ret.m
	return ret
}

func (c *tokenLimiterCounterMutex) Acquire() ReleaseFunc {
	c.m.Lock()
	defer c.m.Unlock()

	for c.cnt <= 0 {
		c.c.Wait()
	}
	c.cnt--

	return func() {
		c.m.Lock()
		defer c.m.Unlock()

		c.cnt++
		c.c.Signal()
	}
}

func (c *tokenLimiterCounterMutex) AcquireNoWait() (ReleaseFunc, error) {
	c.m.Lock()
	defer c.m.Unlock()

	for c.cnt <= 0 {
		return nil, ErrResourceExhausted
	}

	c.cnt--

	return func() {
		c.m.Lock()
		defer c.m.Unlock()

		c.cnt++
		c.c.Signal()
	}, nil
}
