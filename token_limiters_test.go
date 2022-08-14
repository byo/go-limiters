package limiters_test

import (
	"context"
	"reflect"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/byo/go-limiters"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getTypeName(i interface{}) string {
	t := reflect.TypeOf(i)
	if t.Kind() == reflect.Pointer {
		return t.Elem().Name()
	}
	return t.Name()
}

func TestTokenLimiters(t *testing.T) {

	const maxTokens = 5
	const goroutines = 100
	const delay = time.Millisecond

	for _, l := range []limiters.TokenLimiter{
		limiters.NewTokenLimiterChannelEmpty(maxTokens),
		limiters.NewTokenLimiterChanelFull(maxTokens),
		limiters.NewTokenLimiterCounterMutex(maxTokens),
	} {
		t.Run(getTypeName(l), func(t *testing.T) {
			t.Run("Acquire", func(t *testing.T) {

				var currentCount, maxEverCount int64

				st := time.Now()
				wg := sync.WaitGroup{}

				for i := 0; i < goroutines; i++ {
					wg.Add(1)
					go func() {
						defer wg.Done()

						defer l.Acquire()()

						newVal := atomic.AddInt64(&currentCount, 1)
						oldMax := atomic.LoadInt64(&maxEverCount)

						if newVal > oldMax {
							atomic.CompareAndSwapInt64(&maxEverCount, oldMax, newVal)
						}

						time.Sleep(delay)

						defer atomic.AddInt64(&currentCount, -1)
					}()
				}

				wg.Wait()

				assert.LessOrEqual(t, maxEverCount, int64(10))
				assert.GreaterOrEqual(t, time.Since(st), goroutines*delay/maxTokens)
				assert.Less(t, time.Since(st), 2*goroutines*delay/maxTokens)
			})

			t.Run("AcquireNoWait", func(t *testing.T) {

				// Fill up all slots
				for i := 0; i < maxTokens; i++ {
					c, err := l.AcquireNoWait()
					assert.NoError(t, err)
					defer c()
				}

				// One more slot must instantly error out
				releaseFunc, err := l.AcquireNoWait()
				require.ErrorIs(t, err, limiters.ErrResourceExhausted)
				require.Nil(t, releaseFunc)
			})

			t.Run("AcquireCtx", func(t *testing.T) {
				if l, ok := l.(limiters.CancellableTokenLimiter); ok {
					t.Run("WithTimeout", func(t *testing.T) {

						ctx := context.Background()

						// Fill up all slots
						for i := 0; i < maxTokens; i++ {
							c, err := l.AcquireCtx(ctx)
							assert.NoError(t, err)
							defer c()
						}

						// One more slot must fail once the context times out
						ctxTimeout, cancelFunc := context.WithTimeout(ctx, time.Millisecond*20)
						defer cancelFunc()

						c, err := l.AcquireCtx(ctxTimeout)
						require.ErrorIs(t, err, context.DeadlineExceeded)
						require.Nil(t, c)
					})

					t.Run("WithCancel", func(t *testing.T) {

						ctx, cancelFunc := context.WithCancel(context.Background())
						cancelFunc()

						for i := 0; i < maxTokens; i++ {
							c, err := l.AcquireCtx(ctx)
							// We must always fail, even if there are some slots left
							require.ErrorIs(t, err, context.Canceled)
							require.Nil(t, c)
						}
					})
				}
			})
		})
	}
}

func BenchmarkTokenLimiters(b *testing.B) {

	const maxTokens = 5

	for _, l := range []limiters.TokenLimiter{
		limiters.NewTokenLimiterChannelEmpty(maxTokens),
		limiters.NewTokenLimiterChanelFull(maxTokens),
		limiters.NewTokenLimiterCounterMutex(maxTokens),
	} {
		b.Run(getTypeName(l), func(b *testing.B) {
			b.Run("acquire one", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					l.Acquire()()
				}
			})
			b.Run("acquire all", func(b *testing.B) {
				rf := make([]limiters.ReleaseFunc, maxTokens)
				for i := 0; i < b.N; i++ {
					for i := 0; i < maxTokens; i++ {
						rf[i] = l.Acquire()
					}
					for i := maxTokens - 1; i >= 0; i-- {
						rf[i]()
					}
				}
			})
		})

	}
}
