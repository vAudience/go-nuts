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
