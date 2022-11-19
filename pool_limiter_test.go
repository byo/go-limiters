//go:build go1.18
// +build go1.18

package limiters_test

import (
	"sync/atomic"
	"testing"

	"github.com/byo/go-limiters"
	"github.com/stretchr/testify/assert"
)

func TestPoolLimiter(t *testing.T) {

	const poolSize = 10
	const itemSize = 1024 * 1024 * 16

	var itemsCreated int64

	pool := limiters.NewPoolLimiter(
		poolSize,
		func() []byte {
			atomic.AddInt64(&itemsCreated, 1)
			return make([]byte, itemSize)
		},
	)

	assert.EqualValues(t, poolSize, itemsCreated)

	item := pool.Acquire()
	defer pool.Release(item)

	assert.Len(t, item, itemSize)
}

func BenchmarkPoolLimiters(b *testing.B) {
	pool := limiters.NewPoolLimiter(10, func() int { return 0 })

	for i := 0; i < b.N; i++ {
		val := pool.Acquire()
		pool.Release(val)
	}
}
