package gonuts

import (
	"errors"
	"sync"
	"time"
)

// CircuitBreakerState represents the current state of the circuit breaker
type CircuitBreakerState int

const (
	StateClosed CircuitBreakerState = iota
	StateOpen
	StateHalfOpen
)

// CircuitBreaker implements the Circuit Breaker pattern
type CircuitBreaker struct {
	mu sync.Mutex

	failureThreshold uint
	resetTimeout     time.Duration
	halfOpenSuccess  uint

	failures  uint
	successes uint
	state     CircuitBreakerState
	lastError error
	expiry    time.Time
}

// ErrCircuitOpen is returned when the circuit breaker is in the open state
var ErrCircuitOpen = errors.New("circuit breaker is open")

// NewCircuitBreaker creates a new CircuitBreaker
//
// Parameters:
//   - failureThreshold: number of failures before opening the circuit
//   - resetTimeout: duration to wait before attempting to close the circuit
//   - halfOpenSuccess: number of successes in half-open state to close the circuit
//
// Returns:
//   - *CircuitBreaker: a new instance of CircuitBreaker
//
// Example usage:
//
//	cb := gonuts.NewCircuitBreaker(5, 10*time.Second, 2)
//
//	err := cb.Execute(func() error {
//	    return someRiskyOperation()
//	})
//
//	if err != nil {
//	    if errors.Is(err, gonuts.ErrCircuitOpen) {
//	        log.Println("Circuit is open, not attempting operation")
//	    } else {
//	        log.Printf("Operation failed: %v", err)
//	    }
//	}
func NewCircuitBreaker(failureThreshold uint, resetTimeout time.Duration, halfOpenSuccess uint) *CircuitBreaker {
	return &CircuitBreaker{
		failureThreshold: failureThreshold,
		resetTimeout:     resetTimeout,
		halfOpenSuccess:  halfOpenSuccess,
		state:            StateClosed,
	}
}

// Execute runs the given function if the circuit is closed or half-open
//
// Parameters:
//   - f: the function to execute
//
// Returns:
//   - error: nil if the function succeeds, ErrCircuitOpen if the circuit is open,
//     or the error returned by the function
func (cb *CircuitBreaker) Execute(f func() error) error {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	now := time.Now()

	if cb.state == StateOpen {
		if now.After(cb.expiry) {
			cb.state = StateHalfOpen
			cb.failures = 0
			cb.successes = 0
		} else {
			return ErrCircuitOpen
		}
	}

	err := f()

	if err != nil {
		cb.failures++
		cb.lastError = err

		if cb.failures >= cb.failureThreshold {
			cb.state = StateOpen
			cb.expiry = now.Add(cb.resetTimeout)
		}

		return err
	}

	if cb.state == StateHalfOpen {
		cb.successes++

		if cb.successes >= cb.halfOpenSuccess {
			cb.state = StateClosed
		}
	} else {
		// Reset failures on success in closed state
		cb.failures = 0
	}

	return nil
}

// State returns the current state of the circuit breaker
func (cb *CircuitBreaker) State() CircuitBreakerState {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	return cb.state
}

// LastError returns the last error that occurred
func (cb *CircuitBreaker) LastError() error {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	return cb.lastError
}
