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

type PointerType[T TypeType] interface {
	comparable
	GetID() string
	Type() T
}

type PointerMap[K PointerType[T], T TypeType] struct {
	mtx *sync.RWMutex
	m   map[K]struct{}
}

// NewPointerMap init pointer map with type field T
func NewPointerMap[K PointerType[T], T TypeType]() *PointerMap[K, T] {
	return &PointerMap[K, T]{
		mtx: &sync.RWMutex{},
		m:   make(map[K]struct{}),
	}
}

func (m *PointerMap[K, _]) Exists(key K) bool {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	_, ok := m.m[key]
	return ok
}

func (m *PointerMap[K, _]) Add(key K) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	m.m[key] = struct{}{}
}

func (m *PointerMap[K, _]) Remove(key K) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	delete(m.m, key)
}

func (m *PointerMap[_, _]) Length() int {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	return len(m.m)
}

// All is an iterator over the elements of s
func (m *PointerMap[K, _]) All() iter.Seq[K] {
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

// OfType is an iterator over the elements of s with type t
func (m *PointerMap[K, T]) OfType(t T) iter.Seq[K] {
	return func(yield func(K) bool) {
		m.mtx.RLock()
		defer m.mtx.RUnlock()

		for k := range m.m {
			if k.Type() == t {
				if !yield(k) {
					return
				}
			}
		}
	}
}

func (m *PointerMap[K, _]) GetByID(id string) (k K) {
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
// KeyValMap
// ///////////////////////////
type MapType[K MapKey, V MapValue[T], T TypeType] map[K]V

type MapValue[T TypeType] interface {
	comparable
	GetID() string
	Type() T
	Del(bool)
}

type MapKey interface {
	comparable
}

type KeyValMap[K MapKey, V MapValue[T], T TypeType] struct {
	mtx *sync.RWMutex
	m   MapType[K, V, T]
}

// NewKeyValMap create new empty map
func NewKeyValMap[K MapKey, V MapValue[T], T TypeType]() *KeyValMap[K, V, T] {
	return &KeyValMap[K, V, T]{
		mtx: &sync.RWMutex{},
		m:   make(MapType[K, V, T]),
	}
}

// Exists check if key exists
func (m *KeyValMap[K, _, _]) Exists(key K) bool {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	_, ok := m.m[key]
	return ok
}

// Get val with key
func (m *KeyValMap[K, V, _]) Get(key K) V {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	return m.m[key]
}

// Get val with key
func (m *KeyValMap[K, V, _]) GetP(key K, v *V) (ok bool) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	*v, ok = m.m[key]

	return
}

// Get whole map
func (m *KeyValMap[K, V, _]) GetAll() map[K]V {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	return m.m
}

// Set / Overwrite map from map
func (m *KeyValMap[K, V, _]) Set(v map[K]V) {
	m.mtx.Lock()
	m.m = v
	m.mtx.Unlock()

}

// Add key / val to map
func (m *KeyValMap[K, V, _]) Add(k K, v V) {
	m.mtx.Lock()

	m.m[k] = v

	m.mtx.Unlock()
}

// Remove key from map
func (m *KeyValMap[K, _, _]) Remove(key K) {
	m.mtx.Lock()
	delete(m.m, key)

	m.mtx.Unlock()
}

// Mark key as deleted
func (m *KeyValMap[K, _, _]) Delete(key K) {
	m.mtx.Lock()
	m.m[key].Del(true)

	m.mtx.Unlock()
}

// Mark key as not deleted
func (m *KeyValMap[K, _, _]) UnDelete(key K) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	m.m[key].Del(false)
}

// Length of map
func (m *KeyValMap[K, _, _]) Length() int {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	return len(m.m)
}

// All iterate over whole map
func (m *KeyValMap[K, V, _]) All() iter.Seq2[K, V] {
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

// All iterate over elements with type T
func (m *KeyValMap[K, V, T]) OfType(t T) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		m.mtx.RLock()
		defer m.mtx.RUnlock()

		for k, v := range m.m {
			if v.Type() == t {
				if !yield(k, v) {
					return
				}
			}
		}
	}
}
