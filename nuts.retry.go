package gonuts

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

// Retry attempts to execute the given function with exponential backoff.
//
// Parameters:
//   - ctx: A context.Context for cancellation.
//   - attempts: The maximum number of retry attempts.
//   - initialDelay: The initial delay between retries.
//   - maxDelay: The maximum delay between retries.
//   - f: The function to be executed.
//
// Returns:
//   - error: nil if the function succeeds, otherwise the last error encountered.
//
// The function uses exponential backoff with jitter to space out retry attempts.
// It will stop retrying if the context is cancelled or the maximum number of attempts is reached.
//
// Example usage:
//
//	err := gonuts.Retry(ctx, 5, time.Second, time.Minute, func() error {
//	    return someUnreliableOperation()
//	})
//	if err != nil {
//	    log.Printf("Operation failed after retries: %v", err)
//	}
func Retry(ctx context.Context, attempts int, initialDelay, maxDelay time.Duration, f func() error) error {
	var err error
	for i := 0; i < attempts; i++ {
		err = f()
		if err == nil {
			return nil
		}

		if i == attempts-1 {
			break
		}

		select {
		case <-ctx.Done():
			return fmt.Errorf("operation cancelled: %w", ctx.Err())
		case <-time.After(backoffDuration(i, initialDelay, maxDelay)):
		}
	}
	return fmt.Errorf("operation failed after %d attempts: %w", attempts, err)
}

func backoffDuration(attempt int, initialDelay, maxDelay time.Duration) time.Duration {
	delay := initialDelay * (1 << uint(attempt))
	if delay > maxDelay {
		delay = maxDelay
	}
	// Add jitter
	jitter := time.Duration(float64(delay) * (0.5 + rand.Float64()/2))
	return jitter
}
