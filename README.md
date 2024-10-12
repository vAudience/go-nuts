# go-nuts

`go-nuts` is a versatile Go package that provides a collection of utility functions and types to simplify common programming tasks. It offers functionalities ranging from ID generation to logging, state management, concurrency helpers, and more.
go-"nuts" stands for "n-utilities".

## Installation

To install `go-nuts`, use the following command:

```bash
go get github.com/vaudience/go-nuts
```

we alias the package to `nuts` in our code.

```go
import (
    nuts "github.com/vaudience/go-nuts"
)
```

## Functionality Overview

### ID Generation

#### `NanoID(prefix string) string`

Generates a unique ID with a given prefix. It uses a cryptographically secure random number generator to create a unique identifier, then prepends the given prefix to it.

Example:

```go
id := nuts.NanoID("user")
fmt.Println(id) // Output: user_6ByTSYmGzT2c
```

#### `NID(prefix string, length int) string`

Generates a unique ID with a specified length and optional prefix. It creates a unique identifier of the given length using a predefined alphabet. If a prefix is provided, it's prepended to the generated ID.

Example:

```go
id := nuts.NID("doc", 8)
fmt.Println(id) // Output: doc_r3tM9wK1
```

#### `GenerateRandomString(chars []rune, length int) string`

Creates a random string of a given length using the provided character set.

Example:

```go
chars := []rune("abcdefghijklmnopqrstuvwxyz")
randomStr := nuts.GenerateRandomString(chars, 10)
fmt.Println(randomStr) // Output: (a random 10-character string using the given alphabet)
```

### Logging

#### `L *zap.SugaredLogger`

`L` is a pre-initialized logger using the `zap` logging library. It provides a convenient way to log messages at various levels throughout your application.

Example:

```go
nuts.L.Info("This is an info message")
nuts.L.Debug("This is a debug message")
nuts.L.Error("This is an error message")
```

#### `Init_Logger(targetLevel zapcore.Level, instanceId string, log2file bool, logfilePath string) *zap.SugaredLogger`

Initializes a new logger with the specified configuration.

#### `SetLoglevel(loglevel string, instanceId string, log2file bool, logfilePath string)`

Sets the log level for the logger. Available log levels are "DEBUG", "INFO", "WARN", "ERROR", "FATAL", and "PANIC".

Example:

```go
nuts.SetLoglevel("DEBUG", "myapp", true, "/var/log/myapp/")
```

#### `SyslogTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder)`

A custom time encoder for the logger that formats time in syslog format.

#### `GetPrettyJson(object any) string`

Converts an object to a pretty-printed JSON string.

### Concurrent Data Structures

#### `ConcurrentMap[K comparable, V any]`

A thread-safe map implementation that supports concurrent read and write operations.

Example:

```go
cm := nuts.NewConcurrentMap[string, int](16)
cm.Set("key", 42)
value, exists := cm.Get("key")
```

Methods include:

- `Set(key K, value V)`
- `Get(key K) (V, bool)`
- `Delete(key K)`
- `Len() int`
- `Clear()`
- `Keys() []K`
- `Values() []V`
- `Range(f func(K, V) bool)`
- `GetOrSet(key K, value V) (V, bool)`
- `SetIfAbsent(key K, value V) bool`

### State Management

#### `StatesMan`

A flexible state machine manager that supports complex state transitions, timed transitions, and hooks.

Example:

```go
sm := nuts.NewStatesMan("TrafficLight")
sm.AddState("Red", "Red Light", []nuts.SMAction{func(ctx interface{}) { fmt.Println("Red light on") }}, nil)
sm.SetInitialState("Red")
sm.AddTransition("Red", "Green", "Next", nil)
```

Methods include:

- `AddState(id StateID, name string, entryActions, exitActions []SMAction)`
- `AddChildState(parentID, childID StateID) error`
- `SetInitialState(id StateID) error`
- `AddTransition(from, to StateID, event EventID, condition SMCondition, actions ...SMAction)`
- `AddTimedTransition(from, to StateID, duration time.Duration, actions ...SMAction)`
- `TriggerEvent(event EventID)`
- `AddPreHook(hook SMAction)`
- `AddPostHook(hook SMAction)`
- `Run()`
- `GetCurrentState() StateID`
- `Export() (string, error)`
- `Import(jsonStr string) error`
- `GenerateDOT() string`

### JSON Operations

#### `RemoveJsonFields(obj any, fieldsToRemove []string) (string, error)`

Removes specified fields from a JSON object.

Example:

```go
jsonStr, err := nuts.RemoveJsonFields(myObj, []string{"sensitive_field", "internal_id"})
```

#### `SelectJsonFields(obj any, fieldsToSelect []string) (string, error)`

Selects only specified fields from a JSON object.

Example:

```go
jsonStr, err := nuts.SelectJsonFields(myObj, []string{"name", "email", "age"})
```

### Time Operations

#### `TimeFromUnixTimestamp(timestamp int64) time.Time`

Converts a Unix timestamp to a Go `time.Time`.

#### `TimeFromJSTimestamp(timestamp int64) time.Time`

Converts a JavaScript timestamp (milliseconds since epoch) to a Go `time.Time`.

#### `TimeToJSTimestamp(t time.Time) int64`

Converts a Go `time.Time` to a JavaScript timestamp.

### Error Handling

#### `CircuitBreaker`

Implements the Circuit Breaker pattern for fault tolerance in distributed systems.

Example:

```go
cb := nuts.NewCircuitBreaker(5, 10*time.Second, 2)
err := cb.Execute(func() error {
    return someRiskyOperation()
})
```

Methods include:

- `Execute(f func() error) error`
- `State() CircuitBreakerState`
- `LastError() error`

### Password Handling

#### `NormalizePassword(p string) []byte`

Converts a password string to a byte slice.

#### `GeneratePassword(p string) string`

Generates a hashed password using bcrypt.

#### `ComparePasswords(hashedPwd string, inputPwd string) bool`

Compares a hashed password with an input password.

### Byte Size Formatting

#### `BytesToNiceString(size int64) string`

Converts a byte size to a human-readable string (e.g., KB, MB, GB).

#### `BytesToKB/MB/GB/TB(size int64, digits int) float64`

Converts bytes to kilobytes, megabytes, gigabytes, or terabytes.

### Event Emitting

#### `EventEmitter`

A flexible publish-subscribe event system with named listeners.

Example:

```go
emitter := nuts.NewEventEmitter()
emitter.On("userLoggedIn", "logLoginTime", func(username string) {
    fmt.Printf("User logged in: %s at %v\n", username, time.Now())
})
emitter.Emit("userLoggedIn", "JohnDoe")
```

Methods include:

- `On(event, name string, fn interface{}) (string, error)`
- `Off(event, name string) error`
- `Emit(event string, args ...interface{}) error`
- `EmitConcurrent(event string, args ...interface{}) error`
- `Once(event, name string, fn interface{}) (string, error)`
- `ListenerCount(event string) int`
- `ListenerNames(event string) []string`
- `Events() []string`

### Trie Data Structure

#### `Trie`

An efficient tree-like data structure for string operations.

Example:

```go
trie := nuts.NewTrie()
trie.Insert("apple")
fmt.Println(trie.Search("apple"))  // Output: true
fmt.Println(trie.StartsWith("app"))  // Output: true
```

Methods include:

- `Insert(word string)`
- `InsertWithValue(word string, value interface{})`
- `BulkInsert(words []string)`
- `Search(word string) bool`
- `StartsWith(prefix string) bool`
- `AutoComplete(prefix string, limit int) []string`
- `WildcardSearch(pattern string) []string`
- `LongestCommonPrefix() string`

### Rate Limiting

#### `RateLimiter`

Implements a token bucket rate limiter.

Example:

```go
limiter := nuts.NewRateLimiter(10, 100)  // 10 tokens per second, max 100 tokens
if limiter.Allow() {
    // Perform rate-limited operation
}
```

Methods include:

- `Allow() bool`
- `AllowN(n float64) bool`
- `Wait(ctx context.Context) error`
- `WaitN(ctx context.Context, n float64) error`

### Retrying Operations

#### `Retry(ctx context.Context, attempts int, initialDelay, maxDelay time.Duration, f func() error) error`

Attempts to execute the given function with exponential backoff.

Example:

```go
err := nuts.Retry(ctx, 5, time.Second, time.Minute, func() error {
    return someUnreliableOperation()
})
```

### URL Building

#### `URLBuilder`

Provides a fluent interface for constructing URLs.

Example:

```go
builder, _ := nuts.NewURLBuilder("https://api.example.com")
url := builder.AddPath("v1").AddPath("users").AddQuery("page", "1").Build()
```

Methods include:

- `SetScheme(scheme string) *URLBuilder`
- `SetHost(host string) *URLBuilder`
- `SetPort(port string) *URLBuilder`
- `SetCredentials(username, password string) *URLBuilder`
- `AddPath(segment string) *URLBuilder`
- `SetPath(path string) *URLBuilder`
- `AddQuery(key, value string) *URLBuilder`
- `SetQuery(key, value string) *URLBuilder`
- `RemoveQuery(key string) *URLBuilder`
- `SetFragment(fragment string) *URLBuilder`
- `Build() string`
- `BuildURL() (*url.URL, error)`
- `Clone() *URLBuilder`

### Parallel Processing

#### `ParallelSliceMap[T, R any](ctx context.Context, input []T, mapFunc MapFunc[T, R]) ([]R, error)`

Applies a function to each element of a slice concurrently.

#### `ConcurrentMapReduce[T, R any](ctx context.Context, input []T, mapFunc MapFunc[T, R], reduceFunc ReduceFunc[R], initialValue R) (R, error)`

Performs a map-reduce operation concurrently on the input slice.

### Enum Generation

#### `GenerateEnum(def EnumDefinition) (string, error)`

Generates Go code for type-safe enums based on the provided definition.

#### `WriteEnumToFile(def EnumDefinition, filename string) error`

Generates enum code and writes it to a file.

### Miscellaneous Utilities

#### `Debounce(fn any, duration time.Duration, callback func(int)) func(...any)`

Creates a debounced version of a function that delays its execution.

#### `Interval(call func() bool, duration time.Duration, runImmediately bool) *GoInterval`

Creates a new interval that runs a function on a regular interval.

#### `PrintMemoryUsage() bool`

Prints current memory usage statistics.

#### `Set[T comparable]`

A generic set data structure.

#### `JSONPathExtractor`

Extracts values from JSON data using a path-like syntax.

Example:

```go
extractor, _ := nuts.NewJSONPathExtractor(jsonData)
value, _ := extractor.Extract("address.city")
```

### Version Management

#### `Init()`

Initializes version data from a `version.json` file.

#### `GetVersion() string`

Returns the current version.

#### `GetGitCommit() string`

Returns the current Git commit hash.

#### `GetGitBranch() string`

Returns the current Git branch name.

#### `GetVersionData() VersionData`

Returns the complete version data structure.
