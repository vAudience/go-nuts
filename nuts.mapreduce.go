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
