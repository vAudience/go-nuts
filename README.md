# go-nuts

`go-nuts` is a versatile Go library that provides a collection of helpful utilities and data structures to simplify common programming tasks. It offers a wide range of functionalities from string manipulation to concurrent data structures.

## Installation

To install `go-nuts`, use the following command:

```bash
go get github.com/vaudience/go-nuts
```

## Usage

Import the package in your Go code:

```go
import "github.com/vaudience/go-nuts"
```

## Features

### Utility Functions

1. **StringSliceContains**: Checks if a string slice contains a specific string.
2. **StringSliceIndexOf**: Finds the index of a string in a slice.
3. **StringSliceRemoveString**: Removes a specific string from a slice.
4. **StringSliceRemoveIndex**: Removes an element at a specific index from a slice.
5. **StringSlicesHaveSameContent**: Compares two string slices for equality.
6. **RemoveJsonFields**: Removes specified fields from a JSON object.
7. **SelectJsonFields**: Selects only specified fields from a JSON object.
8. **TimeFromUnixTimestamp**: Converts Unix timestamp to Go time.Time.
9. **TimeFromJSTimestamp**: Converts JavaScript timestamp to Go time.Time.
10. **TimeToJSTimestamp**: Converts Go time.Time to JavaScript timestamp.
11. **NormalizePassword**: Converts a password string to a byte slice.
12. **GeneratePassword**: Generates a hashed password.
13. **ComparePasswords**: Compares a hashed password with an input password.
14. **BytesToNiceString**: Converts bytes to a human-readable string (KB, MB, GB, etc.).
15. **BytesToKB/MB/GB/TB**: Converts bytes to kilobytes, megabytes, gigabytes, or terabytes.
16. **NanoID**: Generates a unique ID with a prefix.
17. **NID**: Generates a unique ID with a specified length and optional prefix.
18. **GenerateRandomString**: Generates a random string from a given set of characters.

### Data Structures

1. **ConcurrentMap**: A thread-safe map implementation.
2. **Set**: A generic set data structure.
3. **Trie**: An efficient tree-like data structure for string operations.
4. **CircuitBreaker**: Implements the Circuit Breaker pattern for fault tolerance.
5. **RateLimiter**: Implements a token bucket rate limiter.
6. **EventEmitter**: A flexible publish-subscribe event system.
7. **StatesMan**: A state management system with logging and querying capabilities.

### Algorithms and Patterns

1. **Debounce**: Implements debounce functionality for function calls.
2. **Retry**: Retries an operation with exponential backoff.
3. **ParallelSliceMap**: Applies a function to each element of a slice concurrently.
4. **ConcurrentMapReduce**: Performs a map-reduce operation concurrently on an input slice.

### Code Generation

1. **GenerateEnum**: Generates Go code for type-safe enums.
2. **WriteEnumToFile**: Generates enum code and writes it to a file.

### Logging and Debugging

1. **Init_Logger**: Initializes a logger with custom configuration.
2. **SetLoglevel**: Sets the log level for the logger.
3. **PrintMemoryUsage**: Prints current memory usage statistics.

### JSON and Configuration

1. **JSONPathExtractor**: Extracts values from JSON data using a path-like syntax.
2. **LoadConfigFromYAML**: Loads configuration from a YAML file.
3. **GenerateMarkdownFromFiles**: Generates a markdown file containing code snippets from files in a directory.

## Examples

Here are a few examples of how to use some of the functions in `go-nuts`:

```go
// Using ConcurrentMap
cm := gonuts.NewConcurrentMap[string, int](16)
cm.Set("foo", 42)
value, exists := cm.Get("foo")
if exists {
    fmt.Printf("Value: %d\n", value) // Output: Value: 42
}

// Using RateLimiter
limiter := gonuts.NewRateLimiter(10, 100)  // 10 tokens per second, max 100 tokens
if limiter.Allow() {
    // Perform rate-limited operation
    fmt.Println("Operation allowed")
} else {
    fmt.Println("Operation throttled")
}

// Using Retry
err := gonuts.Retry(context.Background(), 5, time.Second, time.Minute, func() error {
    return someUnreliableOperation()
})
if err != nil {
    log.Printf("Operation failed after retries: %v", err)
}

// Using JSONPathExtractor
jsonData := `{"name": "John", "age": 30, "city": "New York"}`
extractor, _ := gonuts.NewJSONPathExtractor(jsonData)
name, _ := extractor.ExtractString("name")
fmt.Println(name) // Output: John

// Using Trie
trie := gonuts.NewTrie()
trie.Insert("apple")
trie.Insert("app")
fmt.Println(trie.Search("apple")) // Output: true
fmt.Println(trie.StartsWith("app")) // Output: true
```

## Contributing

Contributions to `go-nuts` are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the [MIT License](LICENSE).
