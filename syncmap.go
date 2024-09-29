package syncmap

import (
	"iter"
	"sync"
)

type TypeType interface {
	comparable
}

/////////////////////////////
// PointerMap
/////////////////////////////

type PointerType interface {
	comparable
	GetID() string
	// FilterType() T
}

type PointerMap[K PointerType] struct {
	mtx *sync.RWMutex
	m   map[K]struct{}
}

// NewPointerMap init pointer map with type field T
func NewPointerMap[K PointerType]() *PointerMap[K] {
	return &PointerMap[K]{
		mtx: &sync.RWMutex{},
		m:   make(map[K]struct{}),
	}
}

func (m *PointerMap[K]) Exists(key K) bool {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	_, ok := m.m[key]
	return ok
}

func (m *PointerMap[K]) Add(key K) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	m.m[key] = struct{}{}
}

func (m *PointerMap[K]) Remove(key K) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	delete(m.m, key)
}

func (m *PointerMap[_]) Length() int {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	return len(m.m)
}

// All iterates over the elements of K
func (m *PointerMap[K]) All() iter.Seq[K] {
	return func(yield func(K) bool) {
		m.mtx.RLock()
		defer m.mtx.RUnlock()

		for k := range m.m {
			if !yield(k) {
				return
			}
		}
	}
}

// // OfType is an iterator over the elements of s with type t
// func (m *PointerMap[K, T]) OfType(t T) iter.Seq[K] {
// 	return func(yield func(K) bool) {
// 		m.mtx.RLock()
// 		defer m.mtx.RUnlock()

// 		for k := range m.m {
// 			if k.FilterType() == t {
// 				if !yield(k) {
// 					return
// 				}
// 			}
// 		}
// 	}
// }

func (m *PointerMap[K]) GetByID(id string) (k K) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	for k := range m.m {
		if k.GetID() == id {
			return k
		}
	}
	return k
}

// ///////////////////////////
// Collection
// ///////////////////////////
type MapType[K MapKey, V MapValue] map[K]V

type MapValue interface {
	comparable
	GetID() string
	// FilterType() T
	Del(bool)
}

type MapKey interface {
	comparable
}

type Collection[K MapKey, V MapValue] struct {
	mtx *sync.RWMutex
	m   MapType[K, V]
}

type Bool struct{}

func (t *Bool) GetID() string {
	return ""
}

// func (t *Bool) FilterType() string {
// 	return ""
// }

func (t *Bool) Del(bool) {}

// NewCollection create new empty m: map[K]V
func NewCollection[K MapKey, V MapValue]() *Collection[K, V] {
	return &Collection[K, V]{
		mtx: &sync.RWMutex{},
		m:   make(MapType[K, V]),
	}
}

// Exists check if key exists
func (m *Collection[K, _]) Exists(key K) bool {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	_, ok := m.m[key]
	return ok
}

// Get val with key
func (m *Collection[K, V]) Get(key K) V {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	return m.m[key]
}

// Get val with key
func (m *Collection[K, V]) GetP(key K, v *V) (ok bool) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	*v, ok = m.m[key]

	return
}

// Get whole map
func (m *Collection[K, V]) GetAll() map[K]V {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	return m.m
}

// Set / Overwrite map from map
func (m *Collection[K, V]) Set(v map[K]V) {
	m.mtx.Lock()
	m.m = v
	m.mtx.Unlock()
}

// Add key / val to map
func (m *Collection[K, V]) Add(k K, v V) {
	m.mtx.Lock()
	m.m[k] = v
	m.mtx.Unlock()
}

// Remove key from map
func (m *Collection[K, _]) Remove(key K) {
	m.mtx.Lock()
	delete(m.m, key)

	m.mtx.Unlock()
}

// Mark key as deleted
func (m *Collection[K, _]) Delete(key K) {
	m.mtx.Lock()
	m.m[key].Del(true)

	m.mtx.Unlock()
}

// Mark key as not deleted
func (m *Collection[K, _]) UnDelete(key K) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	m.m[key].Del(false)
}

// Length of map
func (m *Collection[K, _]) Length() int {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	return len(m.m)
}

// All iterates over all elements of K
func (m *Collection[K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		m.mtx.RLock()
		defer m.mtx.RUnlock()

		for k, v := range m.m {
			if !yield(k, v) {
				return
			}
		}
	}
}

// // All iterates over elements with type T
// func (m *Collection[K, V, T]) OfType(t T) iter.Seq2[K, V] {
// 	return func(yield func(K, V) bool) {
// 		m.mtx.RLock()
// 		defer m.mtx.RUnlock()

// 		for k, v := range m.m {
// 			if v.FilterType() == t {
// 				if !yield(k, v) {
// 					return
// 				}
// 			}
// 		}
// 	}
// }
