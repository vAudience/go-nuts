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

// NewConcurrentMap creates a new ConcurrentMap with the specified number of shards.
//
// Example:
//
//	cm := NewConcurrentMap[string, int](16)
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
	hash := fnv32(fmt.Sprintf("%v", key))
	return cm.shards[hash%uint32(cm.numShards)]
}

// Set adds a key-value pair to the map or updates the value if the key already exists.
//
// Example:
//
//	cm.Set("key", 42)
func (cm *ConcurrentMap[K, V]) Set(key K, value V) {
	shard := cm.getShard(key)
	shard.mu.Lock()
	defer shard.mu.Unlock()
	shard.items[key] = value
}

// Get retrieves a value from the map.
//
// Example:
//
//	value, exists := cm.Get("key")
//	if exists {
//	    fmt.Println(value)
//	}
func (cm *ConcurrentMap[K, V]) Get(key K) (V, bool) {
	shard := cm.getShard(key)
	shard.mu.RLock()
	defer shard.mu.RUnlock()
	val, ok := shard.items[key]
	return val, ok
}

// Delete removes a key-value pair from the map.
//
// Example:
//
//	cm.Delete("key")
func (cm *ConcurrentMap[K, V]) Delete(key K) {
	shard := cm.getShard(key)
	shard.mu.Lock()
	defer shard.mu.Unlock()
	delete(shard.items, key)
}

// Len returns the total number of items in the map.
//
// Example:
//
//	count := cm.Len()
//	fmt.Printf("Map contains %d items\n", count)
func (cm *ConcurrentMap[K, V]) Len() int {
	count := 0
	for _, shard := range cm.shards {
		shard.mu.RLock()
		count += len(shard.items)
		shard.mu.RUnlock()
	}
	return count
}

// Clear removes all items from the map.
//
// Example:
//
//	cm.Clear()
func (cm *ConcurrentMap[K, V]) Clear() {
	for _, shard := range cm.shards {
		shard.mu.Lock()
		shard.items = make(map[K]V)
		shard.mu.Unlock()
	}
}

// Keys returns a slice of all keys in the map.
//
// Example:
//
//	keys := cm.Keys()
//	for _, key := range keys {
//	    fmt.Println(key)
//	}
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

// Values returns a slice of all values in the map.
//
// Example:
//
//	values := cm.Values()
//	for _, value := range values {
//	    fmt.Println(value)
//	}
func (cm *ConcurrentMap[K, V]) Values() []V {
	values := make([]V, 0, cm.Len())
	for _, shard := range cm.shards {
		shard.mu.RLock()
		for _, value := range shard.items {
			values = append(values, value)
		}
		shard.mu.RUnlock()
	}
	return values
}

// Range calls the given function for each key-value pair in the map.
//
// Example:
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

// GetOrSet returns the existing value for the key if present.
// Otherwise, it sets and returns the given value.
//
// Example:
//
//	value, loaded := cm.GetOrSet("key", 42)
//	if loaded {
//	    fmt.Println("Key already existed")
//	} else {
//	    fmt.Println("New key-value pair added")
//	}
func (cm *ConcurrentMap[K, V]) GetOrSet(key K, value V) (V, bool) {
	shard := cm.getShard(key)
	shard.mu.Lock()
	defer shard.mu.Unlock()
	if val, ok := shard.items[key]; ok {
		return val, true
	}
	shard.items[key] = value
	return value, false
}

// SetIfAbsent sets the value for a key only if it is not already present.
//
// Example:
//
//	added := cm.SetIfAbsent("key", 42)
//	if added {
//	    fmt.Println("Value was set")
//	} else {
//	    fmt.Println("Key already existed")
//	}
func (cm *ConcurrentMap[K, V]) SetIfAbsent(key K, value V) bool {
	shard := cm.getShard(key)
	shard.mu.Lock()
	defer shard.mu.Unlock()
	if _, ok := shard.items[key]; !ok {
		shard.items[key] = value
		return true
	}
	return false
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
