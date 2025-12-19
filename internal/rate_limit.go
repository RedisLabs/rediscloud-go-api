package internal

import (
	"context"
	"sync"
	"time"
)

type RateLimiter interface {
	// Wait will verify one request can be sent or wait if it can't.
	Wait(ctx context.Context) error
	// Update the rate limiter when the server returns more information about the current limits.
	Update(remaining int) error
}

// A fixedWindowCountRateLimiter is a rate limiter that will count the number of requests within a period (or window)
// and block the caller for the expected remaining period in the window.
//
// The window will start again after the last one closes and the count will be reset.
// Since other requests can happen outside the SDK, callers can calls the Update() function to update the remaining
// event in the window.
//
// This rate limiter tries to model the server-side behaviour as best it can, however, it doesn't know exactly when
// the server-side window starts or ends, so it can be misaligned. Therefore, the callers still need to retry requests
// if a status code 429 (Too Many Requests) is received.
type fixedWindowCountRateLimiter struct {
	limit       int
	period      time.Duration
	windowStart *time.Time
	count       int
	mu          *sync.Mutex
}

func newFixedWindowCountRateLimiter(limit int, period time.Duration) *fixedWindowCountRateLimiter {
	return &fixedWindowCountRateLimiter{
		limit:  limit,
		period: period,
		mu:     &sync.Mutex{},
	}
}

// Wait will block the caller when the number of requests has exceeded the limit in the current window.
// This function allows bursting so it will only block when the limit is reached.
func (rl *fixedWindowCountRateLimiter) Wait(ctx context.Context) error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	// Start window on first requests
	if rl.windowStart == nil {
		now := time.Now()
		rl.windowStart = &now
	}

	windowEnd := rl.windowStart.Add(rl.period)
	if time.Now().After(windowEnd) {
		rl.count = 0
		rl.windowStart = &windowEnd
		windowEnd = rl.windowStart.Add(rl.period)
	}

	if rl.count >= rl.limit {
		delay := windowEnd.Sub(time.Now())
		rl.mu.Unlock()
		err := sleepWithContext(ctx, delay)
		rl.mu.Lock()
		if err != nil {
			return err
		}
		// After sleeping, the window may have reset - recheck and reset count if needed
		now := time.Now()
		windowEnd = rl.windowStart.Add(rl.period)
		if now.After(windowEnd) {
			rl.count = 0
			rl.windowStart = &windowEnd
		}
	}
	rl.count++
	return nil
}

func (rl *fixedWindowCountRateLimiter) Update(remaining int) error {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	rl.count = rl.limit - remaining
	return nil
}

func sleepWithContext(ctx context.Context, d time.Duration) error {
	timer := time.NewTimer(d)
	select {
	case <-ctx.Done():
		if !timer.Stop() {
			<-timer.C // Drain the timer channel to prevent leaks
		}
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

var _ RateLimiter = &fixedWindowCountRateLimiter{}
