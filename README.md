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

### Errors

#### `ErrorPlus`

`ErrorPlus` is an enhanced error type provided by the `gonuts` package, designed to offer robust and structured error handling in Go applications. It encapsulates the original error with additional context, error codes, stack traces, and logging integration, while ensuring immutability and thread safety.

##### Key Features

- **Immutability and Thread Safety**: All methods return new instances, ensuring that errors are immutable and safe to use across goroutines.
- **Custom Error Codes**: Supports custom error codes, including HTTP status codes or application-specific codes.
- **Contextual Data**: Allows attaching arbitrary key-value pairs to the error for additional context.
- **Stack Trace Capture**: Automatically captures the stack trace at the point where the error is created.
- **Logging Integration**: Seamlessly integrates with the `gonuts` logging system (`gonuts.L`) by default but allows for custom logger injection.
- **Error Cause Chains**: Supports error wrapping and unwrapping, compatible with Go's standard error handling mechanisms.
- **JSON Serialization**: Implements custom JSON marshaling to serialize the error, including context and stack trace.
- **Custom Formatting**: Implements the `fmt.Formatter` interface for customizable error output.

##### Creating an ErrorPlus Instance

To create a new `ErrorPlus` instance, use the `NewErrorPlus` function:

```go
err := errors.New("database connection failed")
errPlus := nuts.NewErrorPlus(err, "Unable to reach the database", 500)
```

##### Adding Contextual Information

You can enrich the error with additional context using the `WithContext` method:

```go
errPlus = errPlus.WithContext("userID", 1234).
                  WithContext("operation", "fetchUser")
```

Each call to `WithContext` returns a new `ErrorPlus` instance, preserving immutability.

##### Using Helper Functions for Common HTTP Errors

`ErrorPlus` provides helper functions to create errors with common HTTP status codes:

- `NewNotFoundError(msg string, err error) *ErrorPlus`
- `NewInternalError(msg string, err error) *ErrorPlus`
- `NewUnauthorizedError(msg string, err error) *ErrorPlus`
- `NewBadRequestError(msg string, err error) *ErrorPlus`

Example:

```go
err := errors.New("resource not found")
errPlus := nuts.NewNotFoundError("User does not exist", err)
```

##### Logging the Error

By default, `ErrorPlus` uses the `gonuts.L` logger for logging:

```go
errPlus.Log()
```

You can inject a custom logger if needed:

```go
customLogger := zap.NewExample().Sugar()
nuts.SetErrorLogger(customLogger)
errPlus.Log()
```

##### Accessing Error Information

- **Error Message**: `errPlus.Error()`
- **Error Code**: `errPlus.Code()`
- **Context**: `errPlus.Context()`
- **Stack Trace**: `errPlus.StackTrace()`
- **Timestamp**: `errPlus.Timestamp()`

Example:

```go
fmt.Println("Error Code:", errPlus.Code())
fmt.Println("Context:", errPlus.Context())
fmt.Println("Stack Trace:", errPlus.StackTrace())
```

##### Compatibility with Standard Error Handling

`ErrorPlus` is fully compatible with Go's standard error handling mechanisms:

- **Unwrapping Errors**:

  ```go
  if errors.Is(errPlus, sql.ErrNoRows) {
      // Handle "no rows" error
  }
  ```

- **Type Assertions**:

  ```go
  var pathError *os.PathError
  if errors.As(errPlus, &pathError) {
      // Handle specific path-related error
  }
  ```

##### Custom Formatting

`ErrorPlus` implements the `fmt.Formatter` interface, allowing for detailed formatting:

- **Default Format** (`%v` or `%s`):

  ```go
  fmt.Printf("%v\n", errPlus)
  // Output: Unable to reach the database: database connection failed
  ```

- **Detailed Format** (`%+v`):

  ```go
  fmt.Printf("%+v\n", errPlus)
  /* Output:
  ErrorPlus:
    Msg: Unable to reach the database
    Code: 500
    Error: database connection failed
    Context: map[userID:1234 operation:fetchUser]
    StackTrace:
  main.main
      /path/to/your/app/main.go:25
  runtime.main
      /usr/local/go/src/runtime/proc.go:225
  */
  ```

##### JSON Serialization

`ErrorPlus` can be serialized to JSON, including all its fields:

```go
jsonData, err := json.Marshal(errPlus)
if err != nil {
    // Handle serialization error
}
fmt.Println(string(jsonData))
```

Sample JSON output:

```json
{
  "msg": "Unable to reach the database",
  "code": 500,
  "context": {
    "userID": 1234,
    "operation": "fetchUser"
  },
  "stackTrace": [
    "main.main\n\t/path/to/your/app/main.go:25",
    "runtime.main\n\t/usr/local/go/src/runtime/proc.go:225"
  ],
  "timestamp": "2023-10-10T14:48:00Z",
  "error": "database connection failed"
}
```

##### Full API Reference

###### Creation Methods

- **`NewErrorPlus(err error, msg string, code int) *ErrorPlus`**

  Creates a new `ErrorPlus` instance with the given error, message, and code.

- **Helper Functions for Common Errors**

  - `NewNotFoundError(msg string, err error) *ErrorPlus`
  - `NewInternalError(msg string, err error) *ErrorPlus`
  - `NewUnauthorizedError(msg string, err error) *ErrorPlus`
  - `NewBadRequestError(msg string, err error) *ErrorPlus`

###### Modifier Methods (Immutable)

- **`WithMsg(msg string) *ErrorPlus`**

  Returns a new `ErrorPlus` instance with the updated message.

- **`WithCode(code int) *ErrorPlus`**

  Returns a new `ErrorPlus` instance with the updated code.

- **`WithContext(key string, value interface{}) *ErrorPlus`**

  Returns a new `ErrorPlus` instance with additional context.

- **`WithValues(msg string, code int) *ErrorPlus`**

  Returns a new `ErrorPlus` instance with updated message and code.

###### Accessor Methods

- **`Error() string`**

  Implements the `error` interface.

- **`Unwrap() error`**

  Returns the underlying error.

- **`Is(target error) bool`**

  Checks if the target error matches the underlying error.

- **`As(target interface{}) bool`**

  Attempts to map the `ErrorPlus` to a target error type.

- **`Code() int`**

  Retrieves the error code.

- **`Context() map[string]interface{}`**

  Retrieves a copy of the context map.

- **`StackTrace() []string`**

  Retrieves the captured stack trace.

- **`Timestamp() time.Time`**

  Retrieves the timestamp when the error was created.

###### Logging with ErrorPlus

- **`Log()`**

  Logs the error using the configured logger.

- **`SetErrorLogger(logger *zap.SugaredLogger)`**

  Sets a custom logger for all `ErrorPlus` instances.

###### Formatting and Serialization

- **`Format(f fmt.State, c rune)`**

  Implements the `fmt.Formatter` interface for custom formatting.

- **`MarshalJSON() ([]byte, error)`**

  Implements the `json.Marshaler` interface for JSON serialization.

##### Usage Example

```go
package main

import (
    "database/sql"
    "errors"
    "fmt"

    nuts "github.com/vaudience/go-nuts"
)

func main() {
    // Simulate a database error
    baseErr := sql.ErrNoRows

    // Create an ErrorPlus instance
    errPlus := nuts.NewNotFoundError("User not found in the database", baseErr).
        WithContext("userID", 1234).
        WithContext("operation", "fetchUser")

    // Log the error
    errPlus.Log()

    // Check if the error is a sql.ErrNoRows
    if errors.Is(errPlus, sql.ErrNoRows) {
        fmt.Println("Handle the 'no rows' error")
    }

    // Print the error with detailed formatting
    fmt.Printf("%+v\n", errPlus)
}
```

##### Notes and Best Practices

- **Immutability**: All modifier methods return new instances to ensure thread safety and prevent side effects.
- **Contextual Information**: Use the `WithContext` method to add valuable debugging information.
- **Error Chaining**: Utilize `Unwrap`, `Is`, and `As` to interact with the underlying error.
- **Logging**: Leverage the `Log` method to integrate error logging seamlessly.
- **JSON Serialization**: Be cautious when deserializing errors, as the original error type may not be fully reconstructed.

##### Integration with Other `gonuts` Features

`ErrorPlus` is designed to work harmoniously with other utilities provided by the `gonuts` package. For instance, you can use `ErrorPlus` alongside the `EventEmitter` for emitting errors or within `StateMan` for state transition errors.

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
