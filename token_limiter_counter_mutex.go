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

func (c *tokenLimiterCounterMutex) Acquire() {
	c.m.Lock()
	defer c.m.Unlock()

	for c.cnt <= 0 {
		c.c.Wait()
	}
	c.cnt--
}

func (c *tokenLimiterCounterMutex) AcquireNoWait() error {
	c.m.Lock()
	defer c.m.Unlock()

	for c.cnt <= 0 {
		return ErrResourceExhausted
	}

	c.cnt--

	return nil
}

func (c *tokenLimiterCounterMutex) Release() {
	c.m.Lock()
	defer c.m.Unlock()

	c.cnt++
	c.c.Signal()
}
