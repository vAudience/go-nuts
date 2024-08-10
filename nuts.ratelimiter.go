package gonuts

import (
	"context"
	"sync"
	"time"
)

// RateLimiter implements a token bucket rate limiter
type RateLimiter struct {
	rate       float64
	bucketSize float64
	tokens     float64
	lastRefill time.Time
	mu         sync.Mutex
}

// NewRateLimiter creates a new RateLimiter
//
// Parameters:
//   - rate: the rate at which tokens are added to the bucket (per second)
//   - bucketSize: the maximum number of tokens the bucket can hold
//
// Returns:
//   - *RateLimiter: a new instance of RateLimiter
//
// Example usage:
//
//	limiter := gonuts.NewRateLimiter(10, 100)  // 10 tokens per second, max 100 tokens
//
//	for i := 0; i < 1000; i++ {
//	    if limiter.Allow() {
//	        // Perform rate-limited operation
//	        fmt.Println("Operation allowed:", i)
//	    } else {
//	        fmt.Println("Operation throttled:", i)
//	    }
//	    time.Sleep(time.Millisecond * 50)  // Simulate some work
//	}
func NewRateLimiter(rate, bucketSize float64) *RateLimiter {
	return &RateLimiter{
		rate:       rate,
		bucketSize: bucketSize,
		tokens:     bucketSize,
		lastRefill: time.Now(),
	}
}

// Allow checks if a request is allowed under the rate limit
//
// Returns:
//   - bool: true if the request is allowed, false otherwise
func (rl *RateLimiter) Allow() bool {
	return rl.AllowN(1)
}

// AllowN checks if n requests are allowed under the rate limit
//
// Parameters:
//   - n: the number of tokens to request
//
// Returns:
//   - bool: true if the requests are allowed, false otherwise
func (rl *RateLimiter) AllowN(n float64) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	rl.refill(now)

	if rl.tokens >= n {
		rl.tokens -= n
		return true
	}
	return false
}

// Wait blocks until a request is allowed or the context is cancelled
//
// Parameters:
//   - ctx: a context for cancellation
//
// Returns:
//   - error: nil if a token was acquired, or an error if the context was cancelled
//
// Example usage:
//
//	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
//	defer cancel()
//
//	err := limiter.Wait(ctx)
//	if err != nil {
//	    fmt.Println("Failed to acquire token:", err)
//	    return
//	}
//	// Perform rate-limited operation
func (rl *RateLimiter) Wait(ctx context.Context) error {
	return rl.WaitN(ctx, 1)
}

// WaitN blocks until n requests are allowed or the context is cancelled
//
// Parameters:
//   - ctx: a context for cancellation
//   - n: the number of tokens to request
//
// Returns:
//   - error: nil if tokens were acquired, or an error if the context was cancelled
func (rl *RateLimiter) WaitN(ctx context.Context, n float64) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if rl.AllowN(n) {
				return nil
			}
			time.Sleep(time.Millisecond * 10) // Small sleep to prevent tight loop
		}
	}
}

func (rl *RateLimiter) refill(now time.Time) {
	elapsed := now.Sub(rl.lastRefill).Seconds()
	rl.tokens = min(rl.bucketSize, rl.tokens+elapsed*rl.rate)
	rl.lastRefill = now
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
