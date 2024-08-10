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
