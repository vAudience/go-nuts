# gonuts_code.md

## nuts.sanitize.go

```go
package gonuts

import "strings"

// THIS IS NOT REALLY SAFE! JUST A ROUGH WAY TO MAKE IT A LITTLE HARDER

var SANITIZE_SQLSAFER = []string{"`", "´", "'", " or ", " OR ", "=", ";", ":", "(", ")"}

func SanitizeString(badStringsList []string, stringToClean string) (cleanString string) {
	cleanString = stringToClean
	for _, badThing := range badStringsList {
		cleanString = strings.ReplaceAll(cleanString, badThing, "")
	}
	return cleanString
}
```

## nuts.processstats.go

```go
package gonuts

import (
	"fmt"
	"os"
	"runtime"
)

func PrintMemoryUsage() bool {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	content := fmt.Sprintf("[MemoryUsage] PID=(%d) | Alloc = %v MiB | TotalAlloc = %v MiB | Sys = %v MiB | NumGC = %v", os.Getpid(), bToMb(int64(m.Alloc)), bToMb(int64(m.TotalAlloc)), bToMb(int64(m.Sys)), m.NumGC)
	L.Debugf(content)
	return true
}

func bToMb(b int64) int64 {
	return b / 1024 / 1024
}
```

## nuts.pagination.go

```go
package gonuts

import (
	"fmt"
	"math"
)

// PaginationInfo contains information about the current pagination state
type PaginationInfo struct {
	CurrentPage  int   `json:"current_page"`
	PerPage      int   `json:"per_page"`
	TotalItems   int64 `json:"total_items"`
	TotalPages   int   `json:"total_pages"`
	FirstItem    int   `json:"first_item"`
	LastItem     int   `json:"last_item"`
	FirstPage    int   `json:"first_page"`
	LastPage     int   `json:"last_page"`
	NextPage     *int  `json:"next_page"`
	PreviousPage *int  `json:"previous_page"`
}

// NewPaginationInfo creates a new PaginationInfo instance
//
// Parameters:
//   - currentPage: the current page number
//   - perPage: the number of items per page
//   - totalItems: the total number of items in the dataset
//
// Returns:
//   - *PaginationInfo: a new instance of PaginationInfo
//
// Example usage:
//
//	pagination := gonuts.NewPaginationInfo(2, 10, 95)
//	fmt.Printf("Current Page: %d\n", pagination.CurrentPage)
//	fmt.Printf("Total Pages: %d\n", pagination.TotalPages)
//	fmt.Printf("Next Page: %v\n", *pagination.NextPage)
//
//	// Output:
//	// Current Page: 2
//	// Total Pages: 10
//	// Next Page: 3
func NewPaginationInfo(currentPage, perPage int, totalItems int64) *PaginationInfo {
	if currentPage < 1 {
		currentPage = 1
	}
	if perPage < 1 {
		perPage = 10 // Default to 10 items per page
	}

	totalPages := int(math.Ceil(float64(totalItems) / float64(perPage)))
	if currentPage > totalPages {
		currentPage = totalPages
	}

	firstItem := (currentPage-1)*perPage + 1
	lastItem := firstItem + perPage - 1
	if int64(lastItem) > totalItems {
		lastItem = int(totalItems)
	}

	var nextPage, previousPage *int
	if currentPage < totalPages {
		next := currentPage + 1
		nextPage = &next
	}
	if currentPage > 1 {
		prev := currentPage - 1
		previousPage = &prev
	}

	return &PaginationInfo{
		CurrentPage:  currentPage,
		PerPage:      perPage,
		TotalItems:   totalItems,
		TotalPages:   totalPages,
		FirstItem:    firstItem,
		LastItem:     lastItem,
		FirstPage:    1,
		LastPage:     totalPages,
		NextPage:     nextPage,
		PreviousPage: previousPage,
	}
}

// Offset calculates the offset for database queries
//
// Returns:
//   - int: the offset to use in database queries
func (p *PaginationInfo) Offset() int {
	return (p.CurrentPage - 1) * p.PerPage
}

// Limit returns the number of items per page
//
// Returns:
//   - int: the number of items per page
func (p *PaginationInfo) Limit() int {
	return p.PerPage
}

// HasNextPage checks if there is a next page
//
// Returns:
//   - bool: true if there is a next page, false otherwise
func (p *PaginationInfo) HasNextPage() bool {
	return p.NextPage != nil
}

// HasPreviousPage checks if there is a previous page
//
// Returns:
//   - bool: true if there is a previous page, false otherwise
func (p *PaginationInfo) HasPreviousPage() bool {
	return p.PreviousPage != nil
}

// PageNumbers returns a slice of page numbers to display
//
// Parameters:
//   - max: the maximum number of page numbers to return
//
// Returns:
//   - []int: a slice of page numbers
//
// This method is useful for generating pagination controls in user interfaces.
// It aims to provide a balanced range of page numbers around the current page.
func (p *PaginationInfo) PageNumbers(max int) []int {
	if max >= p.TotalPages {
		pages := make([]int, p.TotalPages)
		for i := range pages {
			pages[i] = i + 1
		}
		return pages
	}

	half := max / 2
	start := p.CurrentPage - half
	end := p.CurrentPage + half

	if start < 1 {
		start = 1
		end = max
	}

	if end > p.TotalPages {
		end = p.TotalPages
		start = p.TotalPages - max + 1
		if start < 1 {
			start = 1
		}
	}

	pages := make([]int, end-start+1)
	for i := range pages {
		pages[i] = start + i
	}
	return pages
}

// String returns a string representation of the pagination info
//
// Returns:
//   - string: a string representation of the pagination info
func (p *PaginationInfo) String() string {
	return fmt.Sprintf("Page %d of %d (Total items: %d, Per page: %d)",
		p.CurrentPage, p.TotalPages, p.TotalItems, p.PerPage)
}
```

## nuts.conditionalexec.go

```go
package gonuts

// Condition represents a function that returns a boolean
type Condition func() bool

// Action represents a function that performs some action
type Action func()

// ConditionalExecution executes actions based on conditions
type ConditionalExecution struct {
	conditions []Condition
	actions    []Action
	elseAction Action
}

// NewConditionalExecution creates a new ConditionalExecution
//
// Returns:
//   - *ConditionalExecution: a new instance of ConditionalExecution
//
// Example usage:
//
//	ce := gonuts.NewConditionalExecution().
//	    If(func() bool { return x > 10 }).
//	    Then(func() { fmt.Println("x is greater than 10") }).
//	    ElseIf(func() bool { return x > 5 }).
//	    Then(func() { fmt.Println("x is greater than 5 but not greater than 10") }).
//	    Else(func() { fmt.Println("x is 5 or less") }).
//	    Execute()
func NewConditionalExecution() *ConditionalExecution {
	return &ConditionalExecution{}
}

// If adds a condition to the execution chain
//
// Parameters:
//   - condition: a function that returns a boolean
//
// Returns:
//   - *ConditionalExecution: the ConditionalExecution instance for method chaining
func (ce *ConditionalExecution) If(condition Condition) *ConditionalExecution {
	ce.conditions = append(ce.conditions, condition)
	return ce
}

// Then adds an action to be executed if the previous condition is true
//
// Parameters:
//   - action: a function to be executed
//
// Returns:
//   - *ConditionalExecution: the ConditionalExecution instance for method chaining
func (ce *ConditionalExecution) Then(action Action) *ConditionalExecution {
	ce.actions = append(ce.actions, action)
	return ce
}

// ElseIf is an alias for If to improve readability
//
// Parameters:
//   - condition: a function that returns a boolean
//
// Returns:
//   - *ConditionalExecution: the ConditionalExecution instance for method chaining
func (ce *ConditionalExecution) ElseIf(condition Condition) *ConditionalExecution {
	return ce.If(condition)
}

// Else adds an action to be executed if all conditions are false
//
// Parameters:
//   - action: a function to be executed
//
// Returns:
//   - *ConditionalExecution: the ConditionalExecution instance for method chaining
func (ce *ConditionalExecution) Else(action Action) *ConditionalExecution {
	ce.elseAction = action
	return ce
}

// Execute runs the conditional execution chain
//
// This method evaluates each condition in order and executes the corresponding
// action for the first true condition. If no conditions are true and an Else
// action is defined, it executes the Else action.
func (ce *ConditionalExecution) Execute() {
	for i, condition := range ce.conditions {
		if condition() {
			if i < len(ce.actions) {
				ce.actions[i]()
			}
			return
		}
	}
	if ce.elseAction != nil {
		ce.elseAction()
	}
}

// ExecuteWithFallthrough runs the conditional execution chain with fallthrough behavior
//
// This method is similar to Execute, but it continues to evaluate conditions and
// execute actions even after a true condition is found, until it encounters a condition
// that returns false or reaches the end of the chain.
func (ce *ConditionalExecution) ExecuteWithFallthrough() {
	for i, condition := range ce.conditions {
		if condition() {
			if i < len(ce.actions) {
				ce.actions[i]()
			}
		} else {
			return
		}
	}
	if ce.elseAction != nil {
		ce.elseAction()
	}
}

// IfThen is a convenience function for simple if-then execution
//
// Parameters:
//   - condition: a function that returns a boolean
//   - action: a function to be executed if the condition is true
//
// Example usage:
//
//	gonuts.IfThen(
//	    func() bool { return x > 10 },
//	    func() { fmt.Println("x is greater than 10") },
//	)
func IfThen(condition Condition, action Action) {
	if condition() {
		action()
	}
}

// IfThenElse is a convenience function for simple if-then-else execution
//
// Parameters:
//   - condition: a function that returns a boolean
//   - thenAction: a function to be executed if the condition is true
//   - elseAction: a function to be executed if the condition is false
//
// Example usage:
//
//	gonuts.IfThenElse(
//	    func() bool { return x > 10 },
//	    func() { fmt.Println("x is greater than 10") },
//	    func() { fmt.Println("x is 10 or less") },
//	)
func IfThenElse(condition Condition, thenAction, elseAction Action) {
	if condition() {
		thenAction()
	} else {
		elseAction()
	}
}
```

## nuts.memoization.go

```go
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
```

## nuts.eventemitter.go

```go
package gonuts

import (
	"fmt"
	"reflect"
	"sync"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

// EventEmitter is a flexible publish-subscribe event system with named listeners
type EventEmitter struct {
	listeners map[string]map[string]reflect.Value
	mu        sync.RWMutex
}

// NewEventEmitter creates a new EventEmitter
//
// Returns:
//   - *EventEmitter: a new instance of EventEmitter
//
// Example usage:
//
//	emitter := gonuts.NewEventEmitter()
func NewEventEmitter() *EventEmitter {
	return &EventEmitter{
		listeners: make(map[string]map[string]reflect.Value),
	}
}

// On subscribes a named function to an event
//
// Parameters:
//   - event: the name of the event to subscribe to
//   - name: a unique name for this listener (if empty, a unique ID will be generated)
//   - fn: the function to be called when the event is emitted
//
// Returns:
//   - string: the name or generated ID of the listener
//   - error: any error that occurred during subscription
//
// Example usage:
//
//	id, err := emitter.On("userLoggedIn", "logLoginTime", func(username string) {
//	    fmt.Printf("User logged in: %s at %v\n", username, time.Now())
//	})
//	if err != nil {
//	    log.Printf("Error subscribing to event: %v", err)
//	}
func (ee *EventEmitter) On(event, name string, fn interface{}) (string, error) {
	if reflect.TypeOf(fn).Kind() != reflect.Func {
		return "", fmt.Errorf("third argument to On must be a function")
	}

	ee.mu.Lock()
	defer ee.mu.Unlock()

	if ee.listeners[event] == nil {
		ee.listeners[event] = make(map[string]reflect.Value)
	}

	if name == "" {
		var err error
		name, err = gonanoid.New()
		if err != nil {
			return "", fmt.Errorf("failed to generate unique ID: %w", err)
		}
	}

	ee.listeners[event][name] = reflect.ValueOf(fn)
	return name, nil
}

// Off unsubscribes a named function from an event
//
// Parameters:
//   - event: the name of the event to unsubscribe from
//   - name: the name or ID of the listener to unsubscribe
//
// Returns:
//   - error: any error that occurred during unsubscription
func (ee *EventEmitter) Off(event, name string) error {
	ee.mu.Lock()
	defer ee.mu.Unlock()

	if listeners, ok := ee.listeners[event]; ok {
		if _, exists := listeners[name]; exists {
			delete(listeners, name)
			return nil
		}
	}
	return fmt.Errorf("listener not found for event: %s, name: %s", event, name)
}

// Emit triggers an event with the given arguments
//
// Parameters:
//   - event: the name of the event to emit
//   - args: the arguments to pass to the event listeners
//
// Returns:
//   - error: any error that occurred during emission
func (ee *EventEmitter) Emit(event string, args ...interface{}) error {
	ee.mu.RLock()
	defer ee.mu.RUnlock()

	if listeners, ok := ee.listeners[event]; ok {
		for _, listener := range listeners {
			if err := ee.callListener(listener, args); err != nil {
				return err
			}
		}
	}
	return nil
}

// EmitConcurrent triggers an event with the given arguments, calling listeners concurrently
//
// Parameters:
//   - event: the name of the event to emit
//   - args: the arguments to pass to the event listeners
//
// Returns:
//   - error: any error that occurred during emission
func (ee *EventEmitter) EmitConcurrent(event string, args ...interface{}) error {
	ee.mu.RLock()
	listeners := ee.listeners[event]
	ee.mu.RUnlock()

	if len(listeners) == 0 {
		return nil
	}

	var wg sync.WaitGroup
	errChan := make(chan error, len(listeners))

	for _, listener := range listeners {
		wg.Add(1)
		go func(l reflect.Value) {
			defer wg.Done()
			if err := ee.callListener(l, args); err != nil {
				errChan <- err
			}
		}(listener)
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

// Once subscribes a one-time named function to an event
//
// Parameters:
//   - event: the name of the event to subscribe to
//   - name: a unique name for this listener (if empty, a unique ID will be generated)
//   - fn: the function to be called when the event is emitted
//
// Returns:
//   - string: the name or generated ID of the listener
//   - error: any error that occurred during subscription
//
// The function will be automatically unsubscribed after it is called once.
func (ee *EventEmitter) Once(event, name string, fn interface{}) (string, error) {
	wrapper := reflect.MakeFunc(reflect.TypeOf(fn), func(args []reflect.Value) []reflect.Value {
		reflect.ValueOf(fn).Call(args)
		ee.Off(event, name)
		return nil
	})
	return ee.On(event, name, wrapper.Interface())
}

// ListenerCount returns the number of listeners for a given event
//
// Parameters:
//   - event: the name of the event
//
// Returns:
//   - int: the number of listeners for the event
func (ee *EventEmitter) ListenerCount(event string) int {
	ee.mu.RLock()
	defer ee.mu.RUnlock()

	return len(ee.listeners[event])
}

// ListenerNames returns a list of all listener names for a given event
//
// Parameters:
//   - event: the name of the event
//
// Returns:
//   - []string: a slice containing all listener names for the event
func (ee *EventEmitter) ListenerNames(event string) []string {
	ee.mu.RLock()
	defer ee.mu.RUnlock()

	names := make([]string, 0, len(ee.listeners[event]))
	for name := range ee.listeners[event] {
		names = append(names, name)
	}
	return names
}

// Events returns a list of all events that have listeners
//
// Returns:
//   - []string: a slice containing all events with listeners
func (ee *EventEmitter) Events() []string {
	ee.mu.RLock()
	defer ee.mu.RUnlock()

	events := make([]string, 0, len(ee.listeners))
	for event := range ee.listeners {
		events = append(events, event)
	}
	return events
}

func (ee *EventEmitter) callListener(listener reflect.Value, args []interface{}) error {
	listenerType := listener.Type()
	if listenerType.NumIn() != len(args) {
		return fmt.Errorf("event handler expects %d arguments, but got %d", listenerType.NumIn(), len(args))
	}

	callArgs := make([]reflect.Value, len(args))
	for i, arg := range args {
		expectedType := listenerType.In(i)
		argValue := reflect.ValueOf(arg)
		if !argValue.Type().AssignableTo(expectedType) {
			return fmt.Errorf("argument %d has wrong type: got %v, want %v", i, argValue.Type(), expectedType)
		}
		callArgs[i] = argValue
	}

	listener.Call(callArgs)
	return nil
}
```

## nuts.logger.go

```go
package gonuts

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// CHECK https://stackoverflow.com/questions/68472667/how-to-log-to-stdout-or-stderr-based-on-log-level-using-uber-go-zap
func Init_Logger(targetLevel zapcore.Level, instanceId string, log2file bool, logfilePath string) *zap.SugaredLogger {
	// fmt.Printf("[nuts.logger] adding logfile: (%s)(%t)(%s)", instanceId, log2file, logfilePath)
	var LogConfig = zap.NewDevelopmentConfig()
	LogConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	LogConfig.EncoderConfig.EncodeTime = SyslogTimeEncoder
	LogConfig.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	LogConfig.Level = zap.NewAtomicLevelAt(targetLevel)
	if log2file {
		if logfilePath != "" {
			logfileName := logfilePath + "log_" + time.Now().Format("2006-01-02T15:04:05Z07:00") + "_" + instanceId + ".txt"
			LogConfig.OutputPaths = append(LogConfig.OutputPaths, logfileName)
			fmt.Printf("[nuts.logger] adding logfile: (%s)", logfileName)
		}
	}
	// LogConfig.Level.SetLevel(zap.DebugLevel)
	logger, err := LogConfig.Build()
	if err != nil {
		fmt.Printf("[nuts.logger] ERROR! failed to create logger PANIC! \n%s", err)
		panic(err)
	}
	defer logger.Sync() // flushes buffer, if any
	return logger.Sugar()
}

func SyslogTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("15:04:05.000"))
}

func SetLoglevel(loglevel string, instanceId string, log2file bool, logfilePath string) {
	switch loglevel {
	case "DEBUG":
		L = Init_Logger(zap.DebugLevel, instanceId, log2file, logfilePath)
		fmt.Println("[SetLoglevel] LogLevel set to DEBUG.")
	case "INFO":
		L = Init_Logger(zap.InfoLevel, instanceId, log2file, logfilePath)
		fmt.Println("[SetLoglevel] LogLevel set to INFO.")
	case "WARN":
		L = Init_Logger(zap.WarnLevel, instanceId, log2file, logfilePath)
		fmt.Println("[SetLoglevel] LogLevel set to WARN.")
	case "ERROR":
		L = Init_Logger(zap.ErrorLevel, instanceId, log2file, logfilePath)
		fmt.Println("[SetLoglevel] LogLevel set to ERROR.")
	case "FATAL":
		L = Init_Logger(zap.FatalLevel, instanceId, log2file, logfilePath)
		fmt.Println("[SetLoglevel] LogLevel set to FATAL.")
	case "PANIC":
		L = Init_Logger(zap.PanicLevel, instanceId, log2file, logfilePath)
		fmt.Println("[SetLoglevel] LogLevel set to PANIC.")
	default:
		L = Init_Logger(zap.DebugLevel, instanceId, log2file, logfilePath)
		fmt.Println("[SetLoglevel] LogLevel set to DEFAULT (DEBUG).")
	}
}

func GetPrettyJson(object any) (pretty string) {
	pretty = ""
	jsonBytes, err := json.Marshal(object)
	if err != nil {
		return "failed to marshal to json :("
	}
	var prettyJSON bytes.Buffer
	json.Indent(&prettyJSON, jsonBytes, "", "\t")
	pretty = prettyJSON.String()
	return pretty
}
```

## nuts.concurrentmap.go

```go
package gonuts

import (
	"fmt"
	"sync"
)

// ConcurrentMap is a thread-safe map implementation
type ConcurrentMap[K comparable, V any] struct {
	shards    []*mapShard[K, V]
	numShards int
}

type mapShard[K comparable, V any] struct {
	items map[K]V
	mu    sync.RWMutex
}

// NewConcurrentMap creates a new ConcurrentMap
//
// Parameters:
//   - numShards: the number of shards to use (must be > 0, default is 32 if <= 0)
//
// Returns:
//   - *ConcurrentMap[K, V]: a new instance of ConcurrentMap
//
// Example usage:
//
//	cm := gonuts.NewConcurrentMap[string, int](16)
//	cm.Set("foo", 42)
//	value, exists := cm.Get("foo")
//	if exists {
//	    fmt.Printf("Value: %d\n", value) // Output: Value: 42
//	}
//	cm.Delete("foo")
func NewConcurrentMap[K comparable, V any](numShards int) *ConcurrentMap[K, V] {
	if numShards <= 0 {
		numShards = 32 // Default number of shards
	}
	cm := &ConcurrentMap[K, V]{
		shards:    make([]*mapShard[K, V], numShards),
		numShards: numShards,
	}
	for i := 0; i < numShards; i++ {
		cm.shards[i] = &mapShard[K, V]{
			items: make(map[K]V),
		}
	}
	return cm
}

func (cm *ConcurrentMap[K, V]) getShard(key K) *mapShard[K, V] {
	hash := fmt.Sprintf("%v", key)
	return cm.shards[fnv32(hash)%uint32(cm.numShards)]
}

// Set adds a key-value pair to the map
//
// Parameters:
//   - key: the key to set
//   - value: the value to set
func (cm *ConcurrentMap[K, V]) Set(key K, value V) {
	shard := cm.getShard(key)
	shard.mu.Lock()
	defer shard.mu.Unlock()
	shard.items[key] = value
}

// Get retrieves a value from the map
//
// Parameters:
//   - key: the key to retrieve
//
// Returns:
//   - V: the value associated with the key
//   - bool: true if the key exists, false otherwise
func (cm *ConcurrentMap[K, V]) Get(key K) (V, bool) {
	shard := cm.getShard(key)
	shard.mu.RLock()
	defer shard.mu.RUnlock()
	val, ok := shard.items[key]
	return val, ok
}

// Delete removes a key-value pair from the map
//
// Parameters:
//   - key: the key to delete
func (cm *ConcurrentMap[K, V]) Delete(key K) {
	shard := cm.getShard(key)
	shard.mu.Lock()
	defer shard.mu.Unlock()
	delete(shard.items, key)
}

// Len returns the total number of items in the map
//
// Returns:
//   - int: the number of items in the map
func (cm *ConcurrentMap[K, V]) Len() int {
	count := 0
	for _, shard := range cm.shards {
		shard.mu.RLock()
		count += len(shard.items)
		shard.mu.RUnlock()
	}
	return count
}

// Clear removes all items from the map
func (cm *ConcurrentMap[K, V]) Clear() {
	for _, shard := range cm.shards {
		shard.mu.Lock()
		shard.items = make(map[K]V)
		shard.mu.Unlock()
	}
}

// Keys returns a slice of all keys in the map
//
// Returns:
//   - []K: a slice containing all keys in the map
func (cm *ConcurrentMap[K, V]) Keys() []K {
	keys := make([]K, 0, cm.Len())
	for _, shard := range cm.shards {
		shard.mu.RLock()
		for key := range shard.items {
			keys = append(keys, key)
		}
		shard.mu.RUnlock()
	}
	return keys
}

// Range calls the given function for each key-value pair in the map
//
// Parameters:
//   - f: the function to call for each key-value pair
//
// The range function is called with two arguments: key and value. If f returns false, range stops the iteration.
//
// Example usage:
//
//	cm.Range(func(key string, value int) bool {
//	    fmt.Printf("Key: %s, Value: %d\n", key, value)
//	    return true // continue iteration
//	})
func (cm *ConcurrentMap[K, V]) Range(f func(K, V) bool) {
	for _, shard := range cm.shards {
		shard.mu.RLock()
		for k, v := range shard.items {
			if !f(k, v) {
				shard.mu.RUnlock()
				return
			}
		}
		shard.mu.RUnlock()
	}
}

// fnv32 is a simple hash function based on FNV-1a
func fnv32(key string) uint32 {
	hash := uint32(2166136261)
	const prime32 = uint32(16777619)
	for i := 0; i < len(key); i++ {
		hash *= prime32
		hash ^= uint32(key[i])
	}
	return hash
}
```

## nuts.slicehelpers.go

```go
package gonuts

import (
	"sort"
)

// Contains checks if a slice contains a specific element.
//
// Example:
//
//	numbers := []int{1, 2, 3, 4, 5}
//	fmt.Println(Contains(numbers, 3)) // Output: true
func Contains[T comparable](slice []T, element T) bool {
	for _, v := range slice {
		if v == element {
			return true
		}
	}
	return false
}

// IndexOf returns the index of the first occurrence of an element in a slice, or -1 if not found.
//
// Example:
//
//	fruits := []string{"apple", "banana", "cherry"}
//	fmt.Println(IndexOf(fruits, "banana")) // Output: 1
func IndexOf[T comparable](slice []T, element T) int {
	for i, v := range slice {
		if v == element {
			return i
		}
	}
	return -1
}

// Remove removes all occurrences of an element from a slice.
//
// Example:
//
//	numbers := []int{1, 2, 3, 2, 4, 2, 5}
//	fmt.Println(Remove(numbers, 2)) // Output: [1 3 4 5]
func Remove[T comparable](slice []T, element T) []T {
	result := make([]T, 0, len(slice))
	for _, v := range slice {
		if v != element {
			result = append(result, v)
		}
	}
	return result
}

// RemoveAt removes the element at a specific index from a slice.
//
// Example:
//
//	letters := []string{"a", "b", "c", "d"}
//	fmt.Println(RemoveAt(letters, 2)) // Output: [a b d]
func RemoveAt[T any](slice []T, index int) []T {
	if index < 0 || index >= len(slice) {
		return slice
	}
	return append(slice[:index], slice[index+1:]...)
}

// Map applies a function to each element of a slice and returns a new slice.
//
// Example:
//
//	numbers := []int{1, 2, 3, 4, 5}
//	squared := Map(numbers, func(n int) int { return n * n })
//	fmt.Println(squared) // Output: [1 4 9 16 25]
func Map[T, R any](slice []T, f func(T) R) []R {
	result := make([]R, len(slice))
	for i, v := range slice {
		result[i] = f(v)
	}
	return result
}

// Filter returns a new slice containing only the elements that satisfy the predicate.
//
// Example:
//
//	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
//	evens := Filter(numbers, func(n int) bool { return n%2 == 0 })
//	fmt.Println(evens) // Output: [2 4 6 8 10]
func Filter[T any](slice []T, predicate func(T) bool) []T {
	result := make([]T, 0, len(slice))
	for _, v := range slice {
		if predicate(v) {
			result = append(result, v)
		}
	}
	return result
}

// Reduce applies a reducer function to all elements in a slice, returning a single value.
//
// Example:
//
//	numbers := []int{1, 2, 3, 4, 5}
//	sum := Reduce(numbers, 0, func(acc, n int) int { return acc + n })
//	fmt.Println(sum) // Output: 15
func Reduce[T, R any](slice []T, initial R, reducer func(R, T) R) R {
	result := initial
	for _, v := range slice {
		result = reducer(result, v)
	}
	return result
}

// Unique returns a new slice with duplicate elements removed.
//
// Example:
//
//	numbers := []int{1, 2, 2, 3, 3, 3, 4, 5, 5}
//	fmt.Println(Unique(numbers)) // Output: [1 2 3 4 5]
func Unique[T comparable](slice []T) []T {
	seen := make(map[T]struct{}, len(slice))
	result := make([]T, 0, len(slice))
	for _, v := range slice {
		if _, ok := seen[v]; !ok {
			seen[v] = struct{}{}
			result = append(result, v)
		}
	}
	return result
}

// Reverse returns a new slice with elements in reverse order.
//
// Example:
//
//	numbers := []int{1, 2, 3, 4, 5}
//	fmt.Println(Reverse(numbers)) // Output: [5 4 3 2 1]
func Reverse[T any](slice []T) []T {
	result := make([]T, len(slice))
	for i, v := range slice {
		result[len(slice)-1-i] = v
	}
	return result
}

// Sort sorts a slice in ascending order (requires a less function).
//
// Example:
//
//	numbers := []int{3, 1, 4, 1, 5, 9, 2, 6}
//	Sort(numbers, func(a, b int) bool { return a < b })
//	fmt.Println(numbers) // Output: [1 1 2 3 4 5 6 9]
func Sort[T any](slice []T, less func(T, T) bool) {
	sort.Slice(slice, func(i, j int) bool {
		return less(slice[i], slice[j])
	})
}

// Chunk splits a slice into chunks of a specified size.
//
// Example:
//
//	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
//	fmt.Println(Chunk(numbers, 3)) // Output: [[1 2 3] [4 5 6] [7 8 9] [10]]
func Chunk[T any](slice []T, chunkSize int) [][]T {
	if chunkSize <= 0 {
		return [][]T{slice}
	}
	chunks := make([][]T, 0, (len(slice)+chunkSize-1)/chunkSize)
	for i := 0; i < len(slice); i += chunkSize {
		end := i + chunkSize
		if end > len(slice) {
			end = len(slice)
		}
		chunks = append(chunks, slice[i:end])
	}
	return chunks
}

// Join concatenates the elements of a string slice, separated by the specified separator.
//
// Example:
//
//	words := []string{"Hello", "World", "Golang"}
//	fmt.Println(Join(words, ", ")) // Output: Hello, World, Golang
func Join(slice []string, separator string) string {
	return Reduce(slice, "", func(acc string, s string) string {
		if acc == "" {
			return s
		}
		return acc + separator + s
	})
}
```

## nuts.interval.go

```go
package gonuts

import "time"

/*
	intervalChannel := Interval(time.Duration(time.Second*1), func() { nuts.L.Debugf("tick ", time.Now()) }, true)
	intervalChannel.Stop()
*/

// creates a new GoInteval struct that allows running a function on a regular interal.
// the call function can trigger a stop of the timer by returning false instead of true
func Interval(call func() bool, duration time.Duration, runImmediately bool) *GoInterval {
	var iv GoInterval = GoInterval{
		tickDuration:     duration,
		call:             call,
		cancelChan:       make(chan bool),
		callHasCancelled: false,
	}
	iv.Start(duration, runImmediately)
	return &iv
}

type GoInterval struct {
	active           bool
	callHasCancelled bool
	tickDuration     time.Duration
	call             func() bool
	ticker           *time.Ticker
	cancelChan       chan (bool)
}

func (iv *GoInterval) Start(duration time.Duration, runImmediately bool) *GoInterval {
	if iv.ticker != nil {
		iv.ticker.Stop()
	}
	iv.tickDuration = duration
	iv.ticker = time.NewTicker(iv.tickDuration)
	iv.active = true
	go func() {
		for {
			select {
			case <-iv.cancelChan:
				return
			case <-iv.ticker.C:
				if iv.call != nil {
					if !iv.call() {
						if iv.active {
							iv.callHasCancelled = true
							iv.Stop()
						}
						return
					}
				}
			}
		}
	}()
	if runImmediately {
		if !iv.call() {
			iv.Stop()
		}
	}
	return iv
}

func (iv *GoInterval) Stop() *GoInterval {
	if !iv.active {
		return iv
	}
	iv.ticker.Stop()
	if !iv.callHasCancelled {
		iv.cancelChan <- true
	}
	iv.ticker.Reset(iv.tickDuration)
	iv.active = false
	return iv
}

func (iv *GoInterval) State() bool {
	return iv.active
}
```

## main.go

```go
package gonuts

import "go.uber.org/zap/zapcore"

var L = Init_Logger(zapcore.DebugLevel, "unknown", false, "logs/")

// @title gonuts package
// @version 0.2.0
// @description a weird collection of little helpers. take anything you want from them.

// @contact.name Toni
// @contact.email i@itsatony.com

// @license.name Unlicense
// @license.url http://unlicense.org/
```

## nuts.ratelimiter.go

```go
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
```

## nuts.pointerstovalues.go

```go
package gonuts

// Generic function to create a pointer to any type
func Ptr[T any](v T) *T {
	return &v
}

// Type-specific functions for common types

// StrPtr returns a pointer to the given string
func StrPtr(s string) *string {
	return &s
}

// IntPtr returns a pointer to the given int
func IntPtr(i int) *int {
	return &i
}

// BoolPtr returns a pointer to the given bool
func BoolPtr(b bool) *bool {
	return &b
}

// Float64Ptr returns a pointer to the given float64
func Float64Ptr(f float64) *float64 {
	return &f
}

// Float32Ptr returns a pointer to the given float32
func Float32Ptr(f float32) *float32 {
	return &f
}

// Int64Ptr returns a pointer to the given int64
func Int64Ptr(i int64) *int64 {
	return &i
}

// Int32Ptr returns a pointer to the given int32
func Int32Ptr(i int32) *int32 {
	return &i
}

// Uint64Ptr returns a pointer to the given uint64
func Uint64Ptr(u uint64) *uint64 {
	return &u
}

// Uint32Ptr returns a pointer to the given uint32
func Uint32Ptr(u uint32) *uint32 {
	return &u
}
```

## nuts.commandlineparameters.go

```go
package gonuts

type CommandLineParameters struct {
	SeedDb     *bool
	EnvFile    *string
	ConfigFile *string
}
```

## nuts.timestamps.go

```go
package gonuts

import (
	"time"
)

func TimeFromUnixTimestamp(timestamp int64) time.Time {
	return time.Unix(timestamp, 0)
}

// this is just to remember that javascript Date.now() converts like this
func TimeFromJSTimestamp(timestamp int64) time.Time {
	return time.Unix(timestamp/1000, 0)
}

// this is just to remember that javascript Date.now() converts like this
func TimeToJSTimestamp(t time.Time) int64 {
	return t.UnixMilli()
}
```

## nuts.urlbuilder.go

```go
package gonuts

import (
	"net/url"
	"path"
	"strings"
)

// URLBuilder provides a fluent interface for constructing URLs
type URLBuilder struct {
	scheme   string
	host     string
	path     string
	query    url.Values
	fragment string
}

// NewURLBuilder creates a new URLBuilder
//
// Parameters:
//   - baseURL: the base URL to start with (optional, can be empty)
//
// Returns:
//   - *URLBuilder: a new instance of URLBuilder
//
// Example usage:
//
//	builder := gonuts.NewURLBuilder("https://api.example.com")
//	url := builder.
//	    AddPath("v1").
//	    AddPath("users").
//	    AddQuery("page", "1").
//	    AddQuery("limit", "10").
//	    SetFragment("top").
//	    Build()
//
//	fmt.Println(url) // Output: https://api.example.com/v1/users?limit=10&page=1#top
func NewURLBuilder(baseURL string) *URLBuilder {
	builder := &URLBuilder{
		query: make(url.Values),
	}

	if baseURL != "" {
		parsed, err := url.Parse(baseURL)
		if err == nil {
			builder.scheme = parsed.Scheme
			builder.host = parsed.Host
			builder.path = parsed.Path
			builder.query = parsed.Query()
			builder.fragment = parsed.Fragment
		}
	}

	return builder
}

// SetScheme sets the scheme (protocol) of the URL
//
// Parameters:
//   - scheme: the scheme to set (e.g., "http", "https")
//
// Returns:
//   - *URLBuilder: the URLBuilder instance for method chaining
func (b *URLBuilder) SetScheme(scheme string) *URLBuilder {
	b.scheme = scheme
	return b
}

// SetHost sets the host of the URL
//
// Parameters:
//   - host: the host to set (e.g., "example.com", "api.example.com:8080")
//
// Returns:
//   - *URLBuilder: the URLBuilder instance for method chaining
func (b *URLBuilder) SetHost(host string) *URLBuilder {
	b.host = host
	return b
}

// AddPath adds a path segment to the URL
//
// Parameters:
//   - segment: the path segment to add
//
// Returns:
//   - *URLBuilder: the URLBuilder instance for method chaining
func (b *URLBuilder) AddPath(segment string) *URLBuilder {
	b.path = path.Join(b.path, segment)
	return b
}

// AddQuery adds a query parameter to the URL
//
// Parameters:
//   - key: the query parameter key
//   - value: the query parameter value
//
// Returns:
//   - *URLBuilder: the URLBuilder instance for method chaining
func (b *URLBuilder) AddQuery(key, value string) *URLBuilder {
	b.query.Add(key, value)
	return b
}

// SetFragment sets the fragment (hash) of the URL
//
// Parameters:
//   - fragment: the fragment to set
//
// Returns:
//   - *URLBuilder: the URLBuilder instance for method chaining
func (b *URLBuilder) SetFragment(fragment string) *URLBuilder {
	b.fragment = fragment
	return b
}

// Build constructs and returns the final URL as a string
//
// Returns:
//   - string: the constructed URL
func (b *URLBuilder) Build() string {
	var sb strings.Builder

	if b.scheme != "" {
		sb.WriteString(b.scheme)
		sb.WriteString("://")
	}

	sb.WriteString(b.host)

	if b.path != "" && b.path != "/" {
		if !strings.HasPrefix(b.path, "/") {
			sb.WriteString("/")
		}
		sb.WriteString(b.path)
	}

	if len(b.query) > 0 {
		sb.WriteString("?")
		sb.WriteString(b.query.Encode())
	}

	if b.fragment != "" {
		sb.WriteString("#")
		sb.WriteString(b.fragment)
	}

	return sb.String()
}

// BuildURL constructs and returns the final URL as a *url.URL
//
// Returns:
//   - *url.URL: the constructed URL
//   - error: any error that occurred during URL parsing
func (b *URLBuilder) BuildURL() (*url.URL, error) {
	return url.Parse(b.Build())
}
```

## nuts.mapreduce.go

```go
package gonuts

import (
	"context"
	"fmt"
	"runtime"
	"sync"
)

// MapFunc is a function that transforms an input value into an output value
type MapFunc[T, R any] func(T) R

// ReduceFunc is a function that combines two values into a single value
type ReduceFunc[R any] func(R, R) R

// ParallelSliceMap applies a function to each element of a slice concurrently
//
// This function processes a slice in parallel, applying the mapFunc to each element.
//
// Parameters:
//   - ctx: A context for cancellation
//   - input: A slice of input values
//   - mapFunc: A function to apply to each input value
//
// Returns:
//   - []R: A slice containing the mapped results
//   - error: An error if the operation was cancelled or failed
//
// Example usage:
//
//	input := []string{"hello", "world", "parallel", "map"}
//
//	result, err := ParallelSliceMap(
//		context.Background(),
//		input,
//		func(s string) int { return len(s) },  // Map: Get length of each string
//	)
//
//	if err != nil {
//		log.Fatalf("ParallelSliceMap failed: %v", err)
//	}
//	fmt.Printf("String lengths: %v\n", result)  // Output: String lengths: [5 5 8 3]
func ParallelSliceMap[T, R any](
	ctx context.Context,
	input []T,
	mapFunc MapFunc[T, R],
) ([]R, error) {
	numWorkers := runtime.GOMAXPROCS(0)
	chunkSize := (len(input) + numWorkers - 1) / numWorkers

	var wg sync.WaitGroup
	results := make([]R, len(input))
	errChan := make(chan error, numWorkers)

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			start := workerID * chunkSize
			end := start + chunkSize
			if end > len(input) {
				end = len(input)
			}

			for j := start; j < end; j++ {
				select {
				case <-ctx.Done():
					errChan <- ctx.Err()
					return
				default:
					results[j] = mapFunc(input[j])
				}
			}
		}(i)
	}

	wg.Wait()
	close(errChan)

	if err := <-errChan; err != nil {
		return nil, fmt.Errorf("map operation cancelled: %w", err)
	}

	return results, nil
}

// ConcurrentMapReduce performs a map-reduce operation concurrently on the input slice
//
// This function applies the mapFunc to each element of the input slice concurrently,
// and then reduces the results using the reduceFunc. It utilizes all available CPU cores
// for maximum performance and uses the existing ConcurrentMap for intermediate results.
//
// Parameters:
//   - ctx: A context for cancellation
//   - input: A slice of input values
//   - mapFunc: A function to apply to each input value
//   - reduceFunc: A function to combine mapped values
//   - initialValue: The initial value for the reduction
//
// Returns:
//   - R: The final reduced result
//   - error: An error if the operation was cancelled or failed
//
// Example usage:
//
//	input := []int{1, 2, 3, 4, 5}
//
//	result, err := ConcurrentMapReduce(
//		context.Background(),
//		input,
//		func(x int) int { return x * x },        // Map: Square each number
//		func(a, b int) int { return a + b },     // Reduce: Sum the squares
//		0,                                       // Initial value for sum
//	)
//
//	if err != nil {
//		log.Fatalf("MapReduce failed: %v", err)
//	}
//	fmt.Printf("Sum of squares: %d\n", result)  // Output: Sum of squares: 55
func ConcurrentMapReduce[T, R any](
	ctx context.Context,
	input []T,
	mapFunc MapFunc[T, R],
	reduceFunc ReduceFunc[R],
	initialValue R,
) (R, error) {
	// Perform parallel mapping
	mappedResults, err := ParallelSliceMap(ctx, input, mapFunc)
	if err != nil {
		return initialValue, fmt.Errorf("mapping operation failed: %w", err)
	}

	// Use ConcurrentMap for intermediate results
	intermediateResults := NewConcurrentMap[int, R](runtime.GOMAXPROCS(0))

	var wg sync.WaitGroup
	errChan := make(chan error, runtime.GOMAXPROCS(0))

	// Perform reduction in parallel
	for i := 0; i < runtime.GOMAXPROCS(0); i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			localResult := initialValue
			for j := workerID; j < len(mappedResults); j += runtime.GOMAXPROCS(0) {
				select {
				case <-ctx.Done():
					errChan <- ctx.Err()
					return
				default:
					localResult = reduceFunc(localResult, mappedResults[j])
				}
			}
			intermediateResults.Set(workerID, localResult)
		}(i)
	}

	wg.Wait()
	close(errChan)

	if err := <-errChan; err != nil {
		return initialValue, fmt.Errorf("reduction operation cancelled: %w", err)
	}

	// Final reduction
	finalResult := initialValue
	intermediateResults.Range(func(key int, value R) bool {
		finalResult = reduceFunc(finalResult, value)
		return true
	})

	return finalResult, nil
}
```

## nuts.retry.go

```go
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
```

## gonuts.enums.go

```go
package gonuts

import (
	"fmt"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"strings"
	"text/template"
)

// EnumDefinition represents the structure of an enum
type EnumDefinition struct {
	Name   string   // The name of the enum type
	Values []string // The values of the enum
}

// GenerateEnum generates the enum code based on the provided definition
//
// This function takes an EnumDefinition and returns a string containing
// the generated Go code for a type-safe enum. The generated code includes:
// - A new type based on int
// - Constant values for each enum value
// - A String() method for string representation
// - An IsValid() method to check if a value is valid
// - A Parse[EnumName]() function to convert strings to enum values
// - MarshalJSON() and UnmarshalJSON() methods for JSON encoding/decoding
//
// Parameters:
//   - def: An EnumDefinition struct containing the enum name and values
//
// Returns:
//   - string: The generated Go code for the enum
//   - error: An error if code generation fails
//
// Example usage:
//
//	def := EnumDefinition{
//	    Name:   "Color",
//	    Values: []string{"Red", "Green", "Blue"},
//	}
//	code, err := GenerateEnum(def)
//	if err != nil {
//	    log.Fatalf("Failed to generate enum: %v", err)
//	}
//	fmt.Println(code)
//
// The generated code can be used as follows:
//
//	var c Color = ColorRed
//	fmt.Println(c)                 // Output: Red
//	fmt.Println(c.IsValid())       // Output: true
//	c2, _ := ParseColor("Green")
//	fmt.Println(c2)                // Output: Green
//	jsonData, _ := json.Marshal(c)
//	fmt.Println(string(jsonData))  // Output: "Red"
func GenerateEnum(def EnumDefinition) (string, error) {
	const enumTemplate = `
// Code generated by gonuts. DO NOT EDIT.

package {{.PackageName}}

import (
	"fmt"
	"encoding/json"
)

// {{.Name}} represents an enumeration of {{.Name}} values
type {{.Name}} int

const (
	{{range $index, $value := .Values}}{{if $index}}_{{end}}{{$.Name}}{{$value}} {{$.Name}} = iota
	{{end}}
)

// {{.Name}}Values contains all valid string representations of {{.Name}}
var {{.Name}}Values = []string{
	{{range .Values}}"{{.}}",
	{{end}}
}

// String returns the string representation of the {{.Name}}
func (e {{.Name}}) String() string {
	if e < 0 || int(e) >= len({{.Name}}Values) {
		return fmt.Sprintf("Invalid{{.Name}}(%d)", int(e))
	}
	return {{.Name}}Values[e]
}

// IsValid checks if the {{.Name}} value is valid
func (e {{.Name}}) IsValid() bool {
	return e >= 0 && int(e) < len({{.Name}}Values)
}

// Parse{{.Name}} converts a string to a {{.Name}} value
func Parse{{.Name}}(s string) ({{.Name}}, error) {
	for i, v := range {{.Name}}Values {
		if v == s {
			return {{.Name}}(i), nil
		}
	}
	return {{.Name}}(-1), fmt.Errorf("invalid {{.Name}}: %s", s)
}

// MarshalJSON implements the json.Marshaler interface for {{.Name}}
func (e {{.Name}}) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface for {{.Name}}
func (e *{{.Name}}) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	v, err := Parse{{.Name}}(s)
	if err != nil {
		return err
	}
	*e = v
	return nil
}
`

	packageName, err := getPackageName()
	if err != nil {
		return "", fmt.Errorf("failed to get package name: %w", err)
	}

	data := struct {
		PackageName string
		EnumDefinition
	}{
		PackageName:    packageName,
		EnumDefinition: def,
	}

	var buf strings.Builder
	tmpl, err := template.New("enum").Parse(enumTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	err = tmpl.Execute(&buf, data)
	if err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

// getPackageName attempts to determine the package name of the current directory
func getPackageName() (string, error) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, ".", nil, parser.PackageClauseOnly)
	if err != nil {
		return "", err
	}

	for name := range pkgs {
		return name, nil
	}

	return "", fmt.Errorf("no package found in current directory")
}

// WriteEnumToFile generates the enum code and writes it to a file
//
// This function generates the enum code based on the provided EnumDefinition
// and writes it to a file. It also formats the generated code for readability.
//
// Parameters:
//   - def: An EnumDefinition struct containing the enum name and values
//   - filename: The name of the file to write the generated code to
//
// Returns:
//   - error: An error if code generation or file writing fails
//
// Example usage:
//
//	def := EnumDefinition{
//	    Name:   "Color",
//	    Values: []string{"Red", "Green", "Blue"},
//	}
//	err := WriteEnumToFile(def, "color_enum.go")
//	if err != nil {
//	    log.Fatalf("Failed to generate enum file: %v", err)
//	}
//
// This will create a file named "color_enum.go" in the current directory
// with the generated enum code. The generated file will contain a type-safe
// Color enum with values ColorRed, ColorGreen, and ColorBlue, along with
// helper methods for string conversion, validation, and JSON marshaling/unmarshaling.
func WriteEnumToFile(def EnumDefinition, filename string) error {
	code, err := GenerateEnum(def)
	if err != nil {
		return err
	}

	// Parse the generated code
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", code, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("failed to parse generated code: %w", err)
	}

	// Format the AST
	var buf strings.Builder
	err = format.Node(&buf, fset, file)
	if err != nil {
		return fmt.Errorf("failed to format code: %w", err)
	}

	// Write the formatted code to file
	err = os.WriteFile(filename, []byte(buf.String()), 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}
```

## nuts.ids.go

```go
package gonuts

import (
	"errors"
	"math/rand"
	"regexp"
	"strconv"
	"time"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

const idAlphabet string = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ-_"

var UUIDRegEx *regexp.Regexp = regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
var NotLegalIdCharacters *regexp.Regexp = regexp.MustCompile("[^A-Za-z0-9-_]")

var ErrBadUUID error = errors.New("uuid format error")
var ErrBadId error = errors.New("bad id format")
var ErrIllegalId error = errors.New("illegal id")
var ErrUnknownId error = errors.New("unknown id")
var ErrMalformedId error = errors.New("malformed id")

func NanoID(prefix string) (nid string) {
	nid, err := gonanoid.Generate(idAlphabet, 12)
	if err != nil {
		L.Error(err)
		nid = strconv.FormatInt(time.Now().UnixMicro(), 10)
	}
	nid = prefix + "_" + nid
	return nid
}

func NID(prefix string, length int) (nid string) {
	nid, err := gonanoid.Generate(idAlphabet, length)
	if err != nil {
		L.Error(err)
		nid = strconv.FormatInt(time.Now().UnixMicro(), 10)
	}
	if len(prefix) > 0 {
		nid = prefix + "_" + nid
	}
	return nid
}

func GenerateRandomString(letters []rune, length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
```

## nuts.circuitbreaker.go

```go
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
```

## nuts.sets.go

```go
package gonuts

import (
	"fmt"
	"strings"
	"sync"
)

// Set is a generic set data structure
type Set[T comparable] struct {
	items map[T]struct{}
	mu    sync.RWMutex
}

// NewSet creates a new Set
//
// Returns:
//   - *Set[T]: a new instance of Set
//
// Example usage:
//
//	intSet := gonuts.NewSet[int]()
//	intSet.Add(1, 2, 3)
//	fmt.Println(intSet.Contains(2)) // Output: true
//
//	stringSet := gonuts.NewSet[string]()
//	stringSet.Add("apple", "banana", "cherry")
//	fmt.Println(stringSet.Size()) // Output: 3
func NewSet[T comparable]() *Set[T] {
	return &Set[T]{
		items: make(map[T]struct{}),
	}
}

// Add adds items to the set
//
// Parameters:
//   - items: the items to add to the set
//
// Returns:
//   - *Set[T]: the Set instance for method chaining
func (s *Set[T]) Add(items ...T) *Set[T] {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, item := range items {
		s.items[item] = struct{}{}
	}
	return s
}

// Remove removes items from the set
//
// Parameters:
//   - items: the items to remove from the set
//
// Returns:
//   - *Set[T]: the Set instance for method chaining
func (s *Set[T]) Remove(items ...T) *Set[T] {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, item := range items {
		delete(s.items, item)
	}
	return s
}

// Contains checks if an item is in the set
//
// Parameters:
//   - item: the item to check for
//
// Returns:
//   - bool: true if the item is in the set, false otherwise
func (s *Set[T]) Contains(item T) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, exists := s.items[item]
	return exists
}

// Size returns the number of items in the set
//
// Returns:
//   - int: the number of items in the set
func (s *Set[T]) Size() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.items)
}

// Clear removes all items from the set
func (s *Set[T]) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items = make(map[T]struct{})
}

// ToSlice returns a slice of all items in the set
//
// Returns:
//   - []T: a slice containing all items in the set
func (s *Set[T]) ToSlice() []T {
	s.mu.RLock()
	defer s.mu.RUnlock()
	slice := make([]T, 0, len(s.items))
	for item := range s.items {
		slice = append(slice, item)
	}
	return slice
}

// Union returns a new set that is the union of this set and another set
//
// Parameters:
//   - other: the other set to union with
//
// Returns:
//   - *Set[T]: a new Set containing the union of both sets
func (s *Set[T]) Union(other *Set[T]) *Set[T] {
	result := NewSet[T]()
	s.mu.RLock()
	defer s.mu.RUnlock()
	other.mu.RLock()
	defer other.mu.RUnlock()
	for item := range s.items {
		result.items[item] = struct{}{}
	}
	for item := range other.items {
		result.items[item] = struct{}{}
	}
	return result
}

// Intersection returns a new set that is the intersection of this set and another set
//
// Parameters:
//   - other: the other set to intersect with
//
// Returns:
//   - *Set[T]: a new Set containing the intersection of both sets
func (s *Set[T]) Intersection(other *Set[T]) *Set[T] {
	result := NewSet[T]()
	s.mu.RLock()
	defer s.mu.RUnlock()
	other.mu.RLock()
	defer other.mu.RUnlock()
	for item := range s.items {
		if _, exists := other.items[item]; exists {
			result.items[item] = struct{}{}
		}
	}
	return result
}

// Difference returns a new set that is the difference of this set and another set
//
// Parameters:
//   - other: the other set to diff with
//
// Returns:
//   - *Set[T]: a new Set containing the difference of this set minus the other set
func (s *Set[T]) Difference(other *Set[T]) *Set[T] {
	result := NewSet[T]()
	s.mu.RLock()
	defer s.mu.RUnlock()
	other.mu.RLock()
	defer other.mu.RUnlock()
	for item := range s.items {
		if _, exists := other.items[item]; !exists {
			result.items[item] = struct{}{}
		}
	}
	return result
}

// String returns a string representation of the set
//
// Returns:
//   - string: a string representation of the set
func (s *Set[T]) String() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	items := make([]string, 0, len(s.items))
	for item := range s.items {
		items = append(items, fmt.Sprintf("%v", item))
	}
	return fmt.Sprintf("Set{%s}", strings.Join(items, ", "))
}
```

## nuts.passwords.go

```go
package gonuts

import "golang.org/x/crypto/bcrypt"

// NormalizePassword func for a returning the users input as a byte slice.
func NormalizePassword(p string) []byte {
	return []byte(p)
}

// GeneratePassword func for a making hash & salt with user password.
func GeneratePassword(p string) string {
	// Normalize password from string to []byte.
	bytePwd := NormalizePassword(p)

	// MinCost is just an integer constant provided by the bcrypt package
	// along with DefaultCost & MaxCost. The cost can be any value
	// you want provided it isn't lower than the MinCost (4).
	hash, err := bcrypt.GenerateFromPassword(bytePwd, bcrypt.MinCost)
	if err != nil {
		return err.Error()
	}

	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it.
	return string(hash)
}

// ComparePasswords func for a comparing password.
func ComparePasswords(hashedPwd string, inputPwd string) bool {
	// Since we'll be getting the hashed password from the DB it will be a string,
	// so we'll need to convert it to a byte slice.
	byteHash := NormalizePassword(hashedPwd)
	byteInput := NormalizePassword(inputPwd)

	// Return result.
	if err := bcrypt.CompareHashAndPassword(byteHash, byteInput); err != nil {
		return false
	}

	return true
}
```

## nuts.debounce.go

```go
package gonuts

import (
	"reflect"
	"sync"
	"time"
)

// Debounce that executes first call immediately and last call after delays in calls
func Debounce(fn any, duration time.Duration, callback func(int)) func(...any) {
	var timer *time.Timer
	var args []reflect.Value
	var callCount int
	var firstCall bool = true
	var mutex sync.Mutex
	fnVal := reflect.ValueOf(fn)

	return func(callArgs ...any) {
		mutex.Lock()
		defer mutex.Unlock()
		callCount++

		// Convert call arguments to reflect.Value
		args = make([]reflect.Value, len(callArgs))
		for i, arg := range callArgs {
			args[i] = reflect.ValueOf(arg)
		}

		if firstCall {
			// Execute the function immediately on the first call
			fnVal.Call(args)
			firstCall = false
		}

		// Reset the timer if it's already set
		if timer != nil {
			timer.Stop()
		}

		timer = time.AfterFunc(duration, func() {
			mutex.Lock()
			defer mutex.Unlock()

			// Call the function with the latest arguments after the duration elapses
			fnVal.Call(args)

			// If a callback is provided, call it with the number of accumulated calls
			if callback != nil {
				callback(callCount)
			}

			// Reset call count and first call state after executing
			callCount = 0
			firstCall = true
		})
	}
}

/*

func main() {
    // Example function to debounce
    printMessage := func(message string) {
        fmt.Println("Message:", message)
    }

    // Example callback to execute after debouncing
    countCalls := func(count int) {
        fmt.Println("Function was called", count, "times")
    }

    // Creating a debounced version of printMessage that executes after 2 seconds of inactivity
    debouncedPrint := Debounce(printMessage, 2*time.Second, countCalls)

    // Simulating rapid calls to the debounced function
    debouncedPrint("Hello")
    time.Sleep(1 * time.Second)
    debouncedPrint("World")
    time.Sleep(1 * time.Second)
    debouncedPrint("Again")
    time.Sleep(3 * time.Second) // Wait enough time to ensure the last call and the callback execute
}

*/
```

## nuts.projectmarkdown.go

```go
package gonuts

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/alecthomas/chroma/lexers"
	"github.com/bmatcuk/doublestar/v4"
	"gopkg.in/yaml.v3"
)

type MarkdownGeneratorConfig struct {
	BaseDir         *string  `yaml:"baseDir"`
	Includes        []string `yaml:"includes"`
	Excludes        []string `yaml:"excludes"`
	BaseHeaderLevel *int     `yaml:"baseHeaderLevel"`
	PrependText     *string  `yaml:"prependText"`
	Title           *string  `yaml:"title"`
	OutputFile      *string  `yaml:"outputFile"`
}

var (
	defaultExcludes = []string{
		"**/vendor/**",
		"**/node_modules/**",
		"**/.git/**",
		"**/.vault/**",
	}

	defaultIncludes = []string{
		"**/*.go", "**/*.js", "**/*.py", "**/*.java", "**/*.c", "**/*.cpp", "**/*.h",
		"**/*.cs", "**/*.rb", "**/*.php", "**/*.swift", "**/*.kt", "**/*.rs",
		"**/*.ts", "**/*.jsx", "**/*.tsx", "**/*.vue", "**/*.scala", "**/*.groovy",
		"**/*.sh", "**/*.bash", "**/*.zsh", "**/*.fish",
		"**/*.sql", "**/*.md", "**/*.yaml", "**/*.yml", "**/*.json", "**/*.xml",
		"**/*.html", "**/*.css", "**/*.scss", "**/*.sass", "**/*.less",
		"**/Dockerfile", "**/Makefile", "**/Jenkinsfile", "**/Gemfile",
		"**/.gitignore", "**/.dockerignore",
	}

	defaultBaseHeaderLevel = 2
	defaultPrependText     = "This is a generated markdown file containing code from the project."
	defaultTitle           = "Project Code Documentation"
)

// // EXAMPLE: Load config from YAML
// config, err := LoadConfigFromYAML("markdown_config.yaml")
// if err != nil {
//     log.Fatalf("Error loading config: %v", err)
// }

// // Generate markdown
// err = GenerateMarkdownFromFiles(config)
// if err != nil {
//     log.Fatalf("Error generating markdown: %v", err)
// }

// // markdown_config.yaml
// baseDir: /path/to/your/project
// includes:
//   - "**/*.go"
//   - "**/*.md"
// excludes:
//   - "**/test/**"

func LoadConfigFromYAML(filename string) (*MarkdownGeneratorConfig, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error reading YAML file: %w", err)
	}

	var config MarkdownGeneratorConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling YAML: %w", err)
	}

	return &config, nil
}

func ApplyDefaults(config *MarkdownGeneratorConfig) error {
	if config.BaseDir == nil {
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("error getting current working directory: %w", err)
		}
		config.BaseDir = &cwd
	}

	if config.BaseHeaderLevel == nil {
		config.BaseHeaderLevel = &defaultBaseHeaderLevel
	}

	if len(config.Excludes) == 0 {
		config.Excludes = defaultExcludes
	}

	if len(config.Includes) == 0 {
		config.Includes = defaultIncludes
	}

	if config.PrependText == nil {
		config.PrependText = &defaultPrependText
	}

	if config.Title == nil {
		config.Title = &defaultTitle
	}

	if config.OutputFile == nil {
		defaultOutputFile := filepath.Join(*config.BaseDir, "project_code.md")
		config.OutputFile = &defaultOutputFile
	}

	return nil
}

// GenerateMarkdownFromFiles generates a markdown file containing code snippets from files in a directory.
//
//	config := &gonuts.MarkdownGeneratorConfig{
//		BaseDir: gonuts.Ptr("/path/to/your/project"),
//		Includes: []string{
//				"**/*.go",
//				"**/*.md",
//				"**/*.yaml",
//		},
//		Excludes: []string{
//				"**/vendor/**",
//				"**/test/**",
//		},
//		Title: gonuts.Ptr("My Custom Project Documentation"),
//		BaseHeaderLevel: gonuts.Ptr(3),
//	}
func GenerateMarkdownFromFiles(config *MarkdownGeneratorConfig) error {
	err := ApplyDefaults(config)
	if err != nil {
		return fmt.Errorf("error applying defaults: %w", err)
	}

	// Validate config
	if *config.BaseDir == "" {
		return fmt.Errorf("base directory is required")
	}
	if *config.OutputFile == "" {
		return fmt.Errorf("output file is required")
	}

	// Add output file to excludes
	config.Excludes = append(config.Excludes, *config.OutputFile)

	// Get all files matching include patterns
	var files []string
	for _, pattern := range config.Includes {
		matches, err := doublestar.Glob(os.DirFS(*config.BaseDir), pattern)
		if err != nil {
			return fmt.Errorf("error matching include pattern %s: %w", pattern, err)
		}
		files = append(files, matches...)
	}

	// Remove excluded files
	files = filterExcludedFiles(files, config.Excludes)

	// Sort files based on the order of include patterns
	sort.Slice(files, func(i, j int) bool {
		return getPatternIndex(files[i], config.Includes) < getPatternIndex(files[j], config.Includes)
	})

	// Generate markdown content
	content := generateMarkdownContent(files, config)

	// Write markdown file
	err = os.WriteFile(*config.OutputFile, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("error writing markdown file: %w", err)
	}

	return nil
}

func filterExcludedFiles(files, excludes []string) []string {
	var result []string
	for _, file := range files {
		excluded := false
		for _, pattern := range excludes {
			match, err := doublestar.Match(pattern, file)
			if err == nil && match {
				excluded = true
				break
			}
		}
		if !excluded {
			result = append(result, file)
		}
	}
	return result
}

func getPatternIndex(file string, patterns []string) int {
	for i, pattern := range patterns {
		match, _ := doublestar.Match(pattern, file)
		if match {
			return i
		}
	}
	return len(patterns)
}

func generateMarkdownContent(files []string, config *MarkdownGeneratorConfig) string {
	var sb strings.Builder

	// Add title
	sb.WriteString(fmt.Sprintf("%s %s\n\n", strings.Repeat("#", *config.BaseHeaderLevel), *config.Title))

	// Add prepend text
	sb.WriteString(*config.PrependText + "\n\n")

	// Add file contents
	for _, file := range files {
		relPath, _ := filepath.Rel(*config.BaseDir, file)
		sb.WriteString(fmt.Sprintf("%s %s\n\n", strings.Repeat("#", *config.BaseHeaderLevel+1), relPath))

		content, err := ioutil.ReadFile(filepath.Join(*config.BaseDir, file))
		if err != nil {
			sb.WriteString(fmt.Sprintf("Error reading file: %s\n\n", err))
			continue
		}

		language := inferLanguage(file, string(content))
		sb.WriteString(fmt.Sprintf("```%s\n%s\n```\n\n", language, string(content)))
	}

	return sb.String()
}

func inferLanguage(filename, content string) string {
	lexer := lexers.Match(filename)
	if lexer == nil {
		lexer = lexers.Analyse(content)
	}
	if lexer == nil {
		return ""
	}
	return lexer.Config().Name
}
```

## nuts.stringSlice.go

```go
package gonuts

import "sort"

// @Summary StringSliceContains checks if a string slice contains a string.
func StringSliceContains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func StringSliceIndexOf(s []string, e string) int {
	for i, a := range s {
		if a == e {
			return i
		}
	}
	return -1
}

func StringSliceRemoveString(max int, sourceSlice []string, stringToRemove string) []string {
	if max == 0 {
		return sourceSlice
	}
	found := []int{}
	for i, a := range sourceSlice {
		if a == stringToRemove {
			found = append(found, i)
		}
	}
	sort.Sort(sort.Reverse(sort.IntSlice(found)))
	var resultSlice []string = sourceSlice
	for _, idx := range found {
		if max == -1 || max > 0 {
			resultSlice = StringSliceRemoveIndex(resultSlice, idx)
		}
	}
	return resultSlice
}

func StringSliceRemoveIndex(s []string, index int) []string {
	return append(s[:index], s[index+1:]...)
}

func StringSlicesHaveSameContent(source []string, compare []string) bool {
	if len(source) != len(compare) {
		return false
	}
	for _, item := range source {
		found := false
		for _, cmp := range compare {
			if item == cmp {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}
```

## nuts.filterjsonfields.go

```go
package gonuts

import "encoding/json"

// from here: https://stackoverflow.com/questions/17306358/removing-fields-from-struct-or-hiding-them-in-json-response by https://stackoverflow.com/users/7496198/chhaileng
func RemoveJsonFields(obj any, fieldsToRemove []string) (string, error) {
	toJson, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}
	if len(fieldsToRemove) == 0 {
		return string(toJson), nil
	}
	toMap := map[string]any{}
	json.Unmarshal([]byte(string(toJson)), &toMap)
	for _, field := range fieldsToRemove {
		delete(toMap, field)
	}
	toJson, err = json.Marshal(toMap)
	if err != nil {
		return "", err
	}
	return string(toJson), nil
}

func SelectJsonFields(obj any, fieldsToSelect []string) (string, error) {
	toJson, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}
	toMap := map[string]any{}
	json.Unmarshal([]byte(string(toJson)), &toMap)
	for key := range toMap {
		if !StringSliceContains(fieldsToSelect, key) {
			delete(toMap, key)
		}
	}
	toJson, err = json.Marshal(toMap)
	if err != nil {
		return "", err
	}
	return string(toJson), nil
}
```

## nuts.version.go

```go
package gonuts

import (
	"encoding/json"
	"os"
)

type VersionData struct {
	Version   string `json:"version"`
	GitCommit string `json:"gitCommit"`
	GitBranch string `json:"gitBranch"`
}

var vData VersionData = VersionData{Version: "0.0.0", GitCommit: "unknown", GitBranch: "unknown"}

func Init() {
	homedir, _ := os.Getwd()
	versionFile, err := os.ReadFile(homedir + "/" + "version.json")
	if err != nil {
		L.Errorf("[version.service] failed to load version file PANIC! \n%s", err)
		L.Panic(err)
	}
	err = json.Unmarshal([]byte(versionFile), &vData)
	if err != nil {
		L.Errorf("[version.service] failed to parse version file PANIC! \n%s", err)
		L.Panic(err)
	}
}

func GetVersion() string {
	return vData.Version
}

func GetGitCommit() string {
	return vData.GitCommit
}

func GetGitBranch() string {
	return vData.GitBranch
}

func GetVersionData() VersionData {
	return vData
}
```

## nuts.jsonpath.go

```go
package gonuts

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// JSONPathExtractor extracts values from JSON data using a path-like syntax
type JSONPathExtractor struct {
	data interface{}
}

// NewJSONPathExtractor creates a new JSONPathExtractor
//
// Parameters:
//   - jsonData: a string containing JSON data
//
// Returns:
//   - *JSONPathExtractor: a new instance of JSONPathExtractor
//   - error: any error that occurred during JSON parsing
//
// Example usage:
//
//	jsonData := `{
//	    "name": "John Doe",
//	    "age": 30,
//	    "address": {
//	        "street": "123 Main St",
//	        "city": "Anytown"
//	    },
//	    "phones": ["555-1234", "555-5678"]
//	}`
//
//	extractor, err := gonuts.NewJSONPathExtractor(jsonData)
//	if err != nil {
//	    log.Fatalf("Error creating extractor: %v", err)
//	}
//
//	name, err := extractor.Extract("name")
//	if err != nil {
//	    log.Fatalf("Error extracting name: %v", err)
//	}
//	fmt.Println("Name:", name)
//
//	city, err := extractor.Extract("address.city")
//	if err != nil {
//	    log.Fatalf("Error extracting city: %v", err)
//	}
//	fmt.Println("City:", city)
//
//	secondPhone, err := extractor.Extract("phones[1]")
//	if err != nil {
//	    log.Fatalf("Error extracting second phone: %v", err)
//	}
//	fmt.Println("Second Phone:", secondPhone)
func NewJSONPathExtractor(jsonData string) (*JSONPathExtractor, error) {
	var data interface{}
	err := json.Unmarshal([]byte(jsonData), &data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}
	return &JSONPathExtractor{data: data}, nil
}

// Extract retrieves a value from the JSON data using the given path
//
// Parameters:
//   - path: a string representing the path to the desired value
//
// Returns:
//   - interface{}: the extracted value
//   - error: any error that occurred during extraction
//
// The path syntax supports:
//   - Dot notation for object properties: "address.street"
//   - Bracket notation for array indices: "phones[0]"
//   - A combination of both: "users[0].name"
func (jpe *JSONPathExtractor) Extract(path string) (interface{}, error) {
	parts := strings.FieldsFunc(path, func(r rune) bool {
		return r == '.' || r == '[' || r == ']'
	})

	var current interface{} = jpe.data
	for _, part := range parts {
		switch v := current.(type) {
		case map[string]interface{}:
			var ok bool
			current, ok = v[part]
			if !ok {
				return nil, fmt.Errorf("key not found: %s", part)
			}
		case []interface{}:
			index, err := strconv.Atoi(part)
			if err != nil {
				return nil, fmt.Errorf("invalid array index: %s", part)
			}
			if index < 0 || index >= len(v) {
				return nil, fmt.Errorf("array index out of bounds: %d", index)
			}
			current = v[index]
		default:
			return nil, fmt.Errorf("cannot navigate further from %T", v)
		}
	}

	return current, nil
}

// ExtractString is a convenience method that extracts a string value
//
// Parameters:
//   - path: a string representing the path to the desired value
//
// Returns:
//   - string: the extracted string value
//   - error: any error that occurred during extraction or type assertion
func (jpe *JSONPathExtractor) ExtractString(path string) (string, error) {
	v, err := jpe.Extract(path)
	if err != nil {
		return "", err
	}
	s, ok := v.(string)
	if !ok {
		return "", fmt.Errorf("value at path %s is not a string", path)
	}
	return s, nil
}

// ExtractInt is a convenience method that extracts an int value
//
// Parameters:
//   - path: a string representing the path to the desired value
//
// Returns:
//   - int: the extracted int value
//   - error: any error that occurred during extraction or type assertion
func (jpe *JSONPathExtractor) ExtractInt(path string) (int, error) {
	v, err := jpe.Extract(path)
	if err != nil {
		return 0, err
	}
	switch n := v.(type) {
	case float64:
		return int(n), nil
	case int:
		return n, nil
	default:
		return 0, fmt.Errorf("value at path %s is not a number", path)
	}
}

// ExtractFloat is a convenience method that extracts a float64 value
//
// Parameters:
//   - path: a string representing the path to the desired value
//
// Returns:
//   - float64: the extracted float64 value
//   - error: any error that occurred during extraction or type assertion
func (jpe *JSONPathExtractor) ExtractFloat(path string) (float64, error) {
	v, err := jpe.Extract(path)
	if err != nil {
		return 0, err
	}
	f, ok := v.(float64)
	if !ok {
		return 0, fmt.Errorf("value at path %s is not a float", path)
	}
	return f, nil
}

// ExtractBool is a convenience method that extracts a bool value
//
// Parameters:
//   - path: a string representing the path to the desired value
//
// Returns:
//   - bool: the extracted bool value
//   - error: any error that occurred during extraction or type assertion
func (jpe *JSONPathExtractor) ExtractBool(path string) (bool, error) {
	v, err := jpe.Extract(path)
	if err != nil {
		return false, err
	}
	b, ok := v.(bool)
	if !ok {
		return false, fmt.Errorf("value at path %s is not a boolean", path)
	}
	return b, nil
}
```

## nuts.trie.go

```go
package gonuts

import (
	"strings"
)

// TrieNode represents a node in the Trie data structure.
type TrieNode struct {
	children map[rune]*TrieNode
	isEnd    bool
	value    interface{}
}

// Trie is a tree-like data structure for efficient string operations.
type Trie struct {
	root *TrieNode
}

// NewTrie creates and returns a new Trie.
//
// Example:
//
//	trie := NewTrie()
func NewTrie() *Trie {
	return &Trie{root: &TrieNode{children: make(map[rune]*TrieNode)}}
}

// Insert adds a word to the Trie.
//
// Example:
//
//	trie := NewTrie()
//	trie.Insert("apple")
func (t *Trie) Insert(word string) {
	node := t.root
	for _, ch := range word {
		if node.children[ch] == nil {
			node.children[ch] = &TrieNode{children: make(map[rune]*TrieNode)}
		}
		node = node.children[ch]
	}
	node.isEnd = true
}

// InsertWithValue adds a word to the Trie with an associated value.
//
// Example:
//
//	trie := NewTrie()
//	trie.InsertWithValue("apple", 42)
func (t *Trie) InsertWithValue(word string, value interface{}) {
	node := t.root
	for _, ch := range word {
		if node.children[ch] == nil {
			node.children[ch] = &TrieNode{children: make(map[rune]*TrieNode)}
		}
		node = node.children[ch]
	}
	node.isEnd = true
	node.value = value
}

// Search checks if a word exists in the Trie.
//
// Example:
//
//	trie := NewTrie()
//	trie.Insert("apple")
//	fmt.Println(trie.Search("apple"))  // Output: true
//	fmt.Println(trie.Search("app"))    // Output: false
func (t *Trie) Search(word string) bool {
	node := t.findNode(word)
	return node != nil && node.isEnd
}

// SearchWithValue checks if a word exists in the Trie and returns its associated value.
//
// Example:
//
//	trie := NewTrie()
//	trie.InsertWithValue("apple", 42)
//	value, found := trie.SearchWithValue("apple")
//	if found {
//	    fmt.Println(value)  // Output: 42
//	}
func (t *Trie) SearchWithValue(word string) (interface{}, bool) {
	node := t.findNode(word)
	if node != nil && node.isEnd {
		return node.value, true
	}
	return nil, false
}

// StartsWith checks if any word in the Trie starts with the given prefix.
//
// Example:
//
//	trie := NewTrie()
//	trie.Insert("apple")
//	fmt.Println(trie.StartsWith("app"))  // Output: true
//	fmt.Println(trie.StartsWith("ban"))  // Output: false
func (t *Trie) StartsWith(prefix string) bool {
	return t.findNode(prefix) != nil
}

// findNode is a helper function to find a node for a given word or prefix.
func (t *Trie) findNode(word string) *TrieNode {
	node := t.root
	for _, ch := range word {
		if node.children[ch] == nil {
			return nil
		}
		node = node.children[ch]
	}
	return node
}

// Delete removes a word from the Trie.
//
// Example:
//
//	trie := NewTrie()
//	trie.Insert("apple")
//	trie.Delete("apple")
//	fmt.Println(trie.Search("apple"))  // Output: false
func (t *Trie) Delete(word string) {
	t.delete(t.root, word, 0)
}

// delete is a recursive helper function for Delete.
func (t *Trie) delete(node *TrieNode, word string, index int) bool {
	if index == len(word) {
		if !node.isEnd {
			return false
		}
		node.isEnd = false
		return len(node.children) == 0
	}

	ch := rune(word[index])
	if node.children[ch] == nil {
		return false
	}

	shouldDeleteChild := t.delete(node.children[ch], word, index+1)

	if shouldDeleteChild {
		delete(node.children, ch)
		return len(node.children) == 0
	}

	return false
}

// FindWordsWithPrefix returns all words in the Trie that start with the given prefix.
//
// Example:
//
//	trie := NewTrie()
//	trie.Insert("apple")
//	trie.Insert("app")
//	trie.Insert("apricot")
//	fmt.Println(trie.FindWordsWithPrefix("app"))  // Output: [app apple]
func (t *Trie) FindWordsWithPrefix(prefix string) []string {
	node := t.findNode(prefix)
	if node == nil {
		return nil
	}

	var words []string
	t.findWordsWithPrefixDFS(node, prefix, &words)
	return words
}

// findWordsWithPrefixDFS is a helper function for FindWordsWithPrefix using depth-first search.
func (t *Trie) findWordsWithPrefixDFS(node *TrieNode, prefix string, words *[]string) {
	if node.isEnd {
		*words = append(*words, prefix)
	}

	for ch, child := range node.children {
		t.findWordsWithPrefixDFS(child, prefix+string(ch), words)
	}
}

// LongestCommonPrefix finds the longest common prefix of all words in the Trie.
//
// Example:
//
//	trie := NewTrie()
//	trie.Insert("flower")
//	trie.Insert("flow")
//	trie.Insert("flight")
//	fmt.Println(trie.LongestCommonPrefix())  // Output: "fl"
func (t *Trie) LongestCommonPrefix() string {
	if t.root == nil {
		return ""
	}

	var sb strings.Builder
	node := t.root

	for len(node.children) == 1 {
		var ch rune
		for k := range node.children {
			ch = k
			break
		}
		sb.WriteRune(ch)
		node = node.children[ch]
		if node.isEnd {
			break
		}
	}

	return sb.String()
}

// CountWords returns the total number of words in the Trie.
//
// Example:
//
//	trie := NewTrie()
//	trie.Insert("apple")
//	trie.Insert("app")
//	trie.Insert("apricot")
//	fmt.Println(trie.CountWords())  // Output: 3
func (t *Trie) CountWords() int {
	return t.countWordsDFS(t.root)
}

// countWordsDFS is a helper function for CountWords using depth-first search.
func (t *Trie) countWordsDFS(node *TrieNode) int {
	count := 0
	if node.isEnd {
		count++
	}

	for _, child := range node.children {
		count += t.countWordsDFS(child)
	}

	return count
}
```

## nuts.safecast.go

```go
package gonuts

import "encoding/json"

func AssignToCastStringOr(val any, fallbackVal string) string {
	castVal, ok := val.(string)
	if ok {
		return castVal
	}
	return fallbackVal
}

func AssignToCastBoolOr(val any, fallbackVal bool) bool {
	castVal, ok := val.(bool)
	if ok {
		return castVal
	}
	return fallbackVal
}

func AssignToCastInt64Or(val any, fallbackVal int64) int64 {
	// L.Debugf("trying to cast to int64 : %d", val)
	switch expType := val.(type) {
	case float64:
		return int64(expType)
	case json.Number:
		castVal, err := expType.Int64()
		if err == nil {
			return castVal
		}
	}
	return fallbackVal
}
```

## nuts.byteskilomega.go

```go
package gonuts

import (
	"math"
	"strconv"
)

func Round(val float64, roundOn float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	newVal = round / pow
	return
}

func BytesToNiceString(size int64) (formattedString string) {
	var newSize float64
	sizeString := ""
	if size > 1024*1024*1024*1024 {
		newSize = BytesToTB(size, 2)
		sizeString = strconv.FormatFloat(newSize, 'f', -1, 64) + "TB"
	} else if size > 1024*1024*1024 {
		newSize = BytesToGB(size, 2)
		sizeString = strconv.FormatFloat(newSize, 'f', -1, 64) + "GB"
	} else if size > 1024*1024 {
		newSize = BytesToMB(size, 2)
		sizeString = strconv.FormatFloat(newSize, 'f', -1, 64) + "MB"
	} else if size > 1024 {
		newSize = BytesToKB(size, 2)
		sizeString = strconv.FormatFloat(newSize, 'f', -1, 64) + "KB"
	} else if size > 99 {
		newSize = BytesToKB(size, 2)
		sizeString = strconv.FormatFloat(newSize, 'f', -1, 64) + "KB"
	} else {
		sizeString = strconv.FormatInt(size, 10) + "B"
	}
	// L.Debugf("[BytesToNiceString] from [%d] to [%s]", size, sizeString)
	return sizeString
}

func BytesToKB(size int64, digits int) (kilobytes float64) {
	newSize := Round(float64(size)/1024, .5, digits)
	return newSize
}
func BytesToMB(size int64, digits int) (megabytes float64) {
	newSize := Round(float64(size)/1024/1024, .5, digits)
	return newSize
}
func BytesToGB(size int64, digits int) (terrabytes float64) {
	newSize := Round(float64(size)/1024/1024/1024, .5, digits)
	return newSize
}
func BytesToTB(size int64, digits int) (terrabytes float64) {
	newSize := Round(float64(size)/1024/1024/1024/1024, .5, digits)
	return newSize
}
```

## nuts.statesman.go

```go
package gonuts

import (
	"sync"
)

type StateString string

const (
	// StateStringEmpty is the state string for empty
	StateStringEmpty StateString = ""
	// StateStringCreated is the state string for created
	StateStringCreated StateString = "created"
	// StateStringStarted is the state string for started
	StateStringStarted StateString = "started"
	// StateStringStopped is the state string for stopped
	StateStringStopped StateString = "stopped"
	// StateStringPaused is the state string for paused
	StateStringPaused StateString = "paused"
	// StateStringResumed is the state string for resumed
	StateStringResumed StateString = "resumed"
	// StateStringCompleted is the state string for completed
	StateStringCompleted StateString = "completed"
	// StateStringFailed is the state string for failed
	StateStringFailed StateString = "failed"
	// StateStringTerminated is the state string for terminated
	StateStringTerminated StateString = "terminated"
)

type StatesManStateChange struct {
	SystemName string
	FromState  StateString
	ToState    StateString
	Timestamp  int64
	Data       map[string]any
}

type StatesMan struct {
	Name               string
	SafetyLock         sync.Mutex
	AvailableStates    []StateString
	LoggedStates       []StatesManStateChange
	CurrentStates      map[string]StateString // map of SystemName to (latest) StateString
	StateChangeChannel chan StatesManStateChange
}

func NewStatesMan(name string) *StatesMan {
	sm := &StatesMan{
		Name:               name,
		AvailableStates:    []StateString{StateStringCreated, StateStringStarted, StateStringStopped, StateStringPaused, StateStringResumed, StateStringCompleted, StateStringFailed, StateStringTerminated},
		LoggedStates:       []StatesManStateChange{},
		CurrentStates:      map[string]StateString{},
		StateChangeChannel: make(chan StatesManStateChange, 10),
	}
	return sm
}

func (s *StatesMan) AddStateChange(systemName string, fromState StateString, toState StateString, timestamp int64, data map[string]any) {
	stateChange := StatesManStateChange{
		SystemName: systemName,
		FromState:  fromState,
		ToState:    toState,
		Timestamp:  timestamp,
		Data:       data,
	}
	s.SafetyLock.Lock()
	s.LoggedStates = append(s.LoggedStates, stateChange)
	s.CurrentStates[systemName] = toState
	s.SafetyLock.Unlock()
	s.StateChangeChannel <- stateChange
}

func (s *StatesMan) GetState(systemName string) StateString {
	s.SafetyLock.Lock()
	defer s.SafetyLock.Unlock()
	stateCopy := s.CurrentStates[systemName]
	return stateCopy
}

func (s *StatesMan) GetStateChangesForSystem(systemName string) []StatesManStateChange {
	s.SafetyLock.Lock()
	defer s.SafetyLock.Unlock()
	var changes []StatesManStateChange
	for _, change := range s.LoggedStates {
		if change.SystemName == systemName {
			changes = append(changes, change)
		}
	}
	return changes
}

func (s *StatesMan) GetFilteredStateChanges(systemName *string, fromState *StateString, toState *StateString, since_timestamp *int64, until_timestamp *int64) []StatesManStateChange {
	s.SafetyLock.Lock()
	defer s.SafetyLock.Unlock()
	var changes []StatesManStateChange
	for _, change := range s.LoggedStates {
		if systemName != nil && change.SystemName != *systemName {
			continue
		}
		if fromState != nil && change.FromState != *fromState {
			continue
		}
		if toState != nil && change.ToState != *toState {
			continue
		}
		if since_timestamp != nil && change.Timestamp < *since_timestamp {
			continue
		}
		if until_timestamp != nil && change.Timestamp > *until_timestamp {
			continue
		}
		changes = append(changes, change)
	}
	return changes
}

func (s *StatesMan) PruneStateChanges(systemName *string, since_timestamp *int64, until_timestamp *int64) {
	s.SafetyLock.Lock()
	defer s.SafetyLock.Unlock()
	var changes []StatesManStateChange
	for _, change := range s.LoggedStates {
		if systemName != nil && change.SystemName != *systemName {
			continue
		}
		if since_timestamp != nil && change.Timestamp < *since_timestamp {
			continue
		}
		if until_timestamp != nil && change.Timestamp > *until_timestamp {
			continue
		}
		changes = append(changes, change)
	}
	s.LoggedStates = changes
}

func (s *StatesMan) CheckStates(states map[string]StateString) bool {
	for systemName, state := range states {
		if !s.CheckState(systemName, state) {
			return false
		}
	}
	return true
}

func (s *StatesMan) CheckState(systemName string, state StateString) bool {
	s.SafetyLock.Lock()
	defer s.SafetyLock.Unlock()
	for _, availableState := range s.AvailableStates {
		if state == availableState {
			return true
		}
	}
	return false
}
```

