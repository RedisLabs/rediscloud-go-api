package internal

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFixedWindowCountRateLimiter_Wait(t *testing.T) {
	windowSize := 2 * time.Second
	windowLimit := 10

	limiter := newFixedWindowCountRateLimiter(windowLimit, windowSize)

	ctx := context.Background()
	start := time.Now()
	runs := 3
	count := 0
	for range windowLimit * runs {
		err := limiter.Wait(ctx)
		require.NoError(t, err)
		count++
	}
	end := time.Now()
	assert.Equal(t, runs*windowLimit, count)
	assert.Greater(t, end.Sub(start), windowSize.Nanoseconds()*int64(runs-1))
}

func TestFixedWindowCountRateLimiter_Update(t *testing.T) {
	windowSize := 2 * time.Second
	windowLimit := 10

	limiter := newFixedWindowCountRateLimiter(windowLimit, windowSize)

	ctx := context.Background()
	start := time.Now()
	runs := 2
	count := 0
	assert.NoError(t, limiter.Update(0))
	for range windowLimit * runs {
		t.Logf("%s\n", time.Now().String())
		err := limiter.Wait(ctx)
		require.NoError(t, err)
		count++
	}
	end := time.Now()
	assert.Equal(t, runs*windowLimit, count)
	assert.Greater(t, end.Sub(start), windowSize.Nanoseconds()*int64(runs))

}
