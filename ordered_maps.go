package gomaps

import (
	"fmt"
	"sync"
)

// OrderedMap is a threadsafe generic map that preserves the insertion order of key-value pairs.
// It allows O(1) access by key and maintains a consistent order for iteration.
// Deletions is an O(n) operation due to the need to shift elements in the underlying structure.
type OrderedMap[K comparable, V any] struct {
	keys    map[K]int
	indexes map[int]K
	values  []V
	mu      *sync.RWMutex
}

// NewOrderedMap returns an empty initialised ordered map.
func NewOrderedMap[K comparable, V any]() OrderedMap[K, V] {
	return OrderedMap[K, V]{
		keys:    make(map[K]int),
		indexes: make(map[int]K),
		values:  []V{},
		mu:      &sync.RWMutex{},
	}
}

// Set adds the key value pair to the map if not present, otherwise overwrites existing values.
// If existing values must not be overwritten, use the `Push` method instead.
func (om *OrderedMap[K, V]) Set(key K, value V) {
	om.mu.Lock()
	defer om.mu.Unlock()

	idx, exists := om.keys[key]
	if exists {
		om.values[idx] = value
	} else {
		idx = len(om.values)
		om.keys[key] = idx
		om.indexes[idx] = key
		om.values = append(om.values, value)
	}
}

// Push adds the key value pair to the map or returns an error if the key value pair already exists.
// If existing values can be overwritten, use the `Set` method instead.
func (om *OrderedMap[K, V]) Push(key K, value V) error {
	om.mu.Lock()
	defer om.mu.Unlock()

	_, exists := om.keys[key]
	if exists {
		return fmt.Errorf("key %v already exists", key)
	}

	idx := len(om.values)
	om.keys[key] = idx
	om.indexes[idx] = key
	om.values = append(om.values, value)

	return nil
}

// Get returns the value mapped to by the key if it exists and a success bool.
func (om *OrderedMap[K, V]) Get(key K) (V, bool) {
	om.mu.RLock()
	defer om.mu.RUnlock()

	idx, ok := om.keys[key]
	if !ok {
		var zero V
		return zero, false
	}

	return om.values[idx], true
}

// Delete removes the value mapped to by key if it exists.
// If no key value pair exists, the function results in a no-op.
func (om *OrderedMap[K, V]) Delete(key K) {
	om.mu.Lock()
	defer om.mu.Unlock()

	idx, exists := om.keys[key]
	if !exists {
		return
	}

	delete(om.keys, key)
	delete(om.indexes, idx)
	if idx == len(om.values)-1 {
		om.values = om.values[:idx]
	} else {
		om.values = append(om.values[:idx], om.values[idx+1:]...)
	}

	for i := idx; i < len(om.values); i++ {
		k := om.indexes[i+1]
		om.keys[k] = i
		om.indexes[i] = k
		delete(om.indexes, i+1)
	}
}

// Range calls the function `fn` for each key-value pair in the OrderedMap.
// If `fn` returns false, iteration will immediately stop.
// This is a readonly method and will deadlock on attempts to modify the underlying data.
// To modify the underlying data, use the `RangeUnsafe` method.
func (om *OrderedMap[K, V]) Range(fn func(key K, value V) bool) {
	om.mu.RLock()
	defer om.mu.RUnlock()

	for i, value := range om.values {
		if !fn(om.indexes[i], value) {
			return
		}
	}
}

// RangeUnsafe calls the function `fn` for each key-value pair in the OrderedMap.
// If `fn` returns false, iteration will immediately stop.
// Allows for the updating of values during iteration.
func (om *OrderedMap[K, V]) RangeUnsafe(fn func(key K, value V) bool) {
	om.mu.RLock()

	type kv struct {
		key   K
		value V
	}

	snapshot := make([]kv, len(om.values))
	for i, v := range om.values {
		snapshot[i] = kv{key: om.indexes[i], value: v}
	}

	om.mu.RUnlock()

	for _, item := range snapshot {
		if !fn(item.key, item.value) {
			return
		}
	}
}

// Reverse sorts the OrderedMap inplace in the reverse order.
func (om *OrderedMap[K, V]) Reverse() {
	om.mu.Lock()
	defer om.mu.Unlock()

	reversed := make([]V, len(om.values))
	reverseIdx := len(om.values) - 1

	for idx := range len(om.values) {
		key := om.indexes[idx]
		valueIdx := om.keys[key]

		om.keys[key] = reverseIdx
		om.indexes[reverseIdx] = key

		reversed[reverseIdx] = om.values[valueIdx]
		reverseIdx--
	}

	om.values = reversed
}

func (om *OrderedMap[K, V]) Concat(oms ...OrderedMap[K, V]) *OrderedMap[K, V] {
	result := *om

	for _, m := range oms {
		m.Range(func(key K, value V) bool {
			result.Set(key, value)
			return true
		})
	}

	return &result
}
