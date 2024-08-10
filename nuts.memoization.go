package gonuts

import (
	"fmt"
	"reflect"
	"sync"
	"time"
)

// MemoizedFunc is a wrapper for a memoized function
type MemoizedFunc struct {
	mu    sync.RWMutex
	cache map[string]cacheEntry
	f     interface{}
	ttl   time.Duration
}

type cacheEntry struct {
	result interface{}
	expiry time.Time
}

// Memoize creates a memoized version of the given function
//
// Parameters:
//   - f: the function to memoize (must be a function type)
//   - ttl: time-to-live for cached results (use 0 for no expiration)
//
// Returns:
//   - *MemoizedFunc: a memoized version of the input function
//
// The memoized function will cache results based on input parameters.
// Subsequent calls with the same parameters will return the cached result.
// Cached results expire after the specified TTL (if non-zero).
//
// Example usage:
//
//	expensiveFunc := func(x int) int {
//	    time.Sleep(time.Second) // Simulate expensive operation
//	    return x * 2
//	}
//
//	memoized := gonuts.Memoize(expensiveFunc, 5*time.Minute)
//
//	start := time.Now()
//	result, err := memoized.Call(42)
//	fmt.Printf("First call took %v: %v\n", time.Since(start), result)
//
//	start = time.Now()
//	result, err = memoized.Call(42)
//	fmt.Printf("Second call took %v: %v\n", time.Since(start), result)
func Memoize(f interface{}, ttl time.Duration) *MemoizedFunc {
	return &MemoizedFunc{
		cache: make(map[string]cacheEntry),
		f:     f,
		ttl:   ttl,
	}
}

// Call invokes the memoized function with the given arguments
//
// Parameters:
//   - args: the arguments to pass to the memoized function
//
// Returns:
//   - interface{}: the result of the function call
//   - error: any error that occurred during the function call or type checking
func (m *MemoizedFunc) Call(args ...interface{}) (interface{}, error) {
	key := fmt.Sprintf("%v", args)

	m.mu.RLock()
	entry, found := m.cache[key]
	m.mu.RUnlock()

	if found && (m.ttl == 0 || time.Now().Before(entry.expiry)) {
		return entry.result, nil
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// Check again in case another goroutine updated the cache
	entry, found = m.cache[key]
	if found && (m.ttl == 0 || time.Now().Before(entry.expiry)) {
		return entry.result, nil
	}

	v := reflect.ValueOf(m.f)
	t := v.Type()

	if t.Kind() != reflect.Func {
		return nil, fmt.Errorf("memoize: not a function")
	}

	if len(args) != t.NumIn() {
		return nil, fmt.Errorf("memoize: wrong number of arguments: got %d, want %d", len(args), t.NumIn())
	}

	var in []reflect.Value
	for i, arg := range args {
		if reflect.TypeOf(arg) != t.In(i) {
			return nil, fmt.Errorf("memoize: argument %d has wrong type: got %T, want %v", i, arg, t.In(i))
		}
		in = append(in, reflect.ValueOf(arg))
	}

	result := v.Call(in)

	var out interface{}
	if len(result) > 0 {
		out = result[0].Interface()
	}

	m.cache[key] = cacheEntry{
		result: out,
		expiry: time.Now().Add(m.ttl),
	}

	return out, nil
}
