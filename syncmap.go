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
	FilterType() T
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

// All iterates over the elements of K
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
			if k.FilterType() == t {
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
// Collection
// ///////////////////////////
type MapType[K MapKey, V MapValue[T], T TypeType] map[K]V

type MapValue[T TypeType] interface {
	comparable
	GetID() string
	FilterType() T
	Del(bool)
}

type MapKey interface {
	comparable
}

type Collection[K MapKey, V MapValue[T], T TypeType] struct {
	mtx *sync.RWMutex
	m   MapType[K, V, T]
}

type Bool struct{}

func (t *Bool) GetID() string {
	return ""
}

func (t *Bool) FilterType() string {
	return ""
}

func (t *Bool) Del(bool) {}

// NewCollection create new empty m: map[K]V
func NewCollection[K MapKey, V MapValue[T], T TypeType]() *Collection[K, V, T] {
	return &Collection[K, V, T]{
		mtx: &sync.RWMutex{},
		m:   make(MapType[K, V, T]),
	}
}

// Exists check if key exists
func (m *Collection[K, _, _]) Exists(key K) bool {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	_, ok := m.m[key]
	return ok
}

// Get val with key
func (m *Collection[K, V, _]) Get(key K) V {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	return m.m[key]
}

// Get val with key
func (m *Collection[K, V, _]) GetP(key K, v *V) (ok bool) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	*v, ok = m.m[key]

	return
}

// Get whole map
func (m *Collection[K, V, _]) GetAll() map[K]V {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	return m.m
}

// Set / Overwrite map from map
func (m *Collection[K, V, _]) Set(v map[K]V) {
	m.mtx.Lock()
	m.m = v
	m.mtx.Unlock()
}

// Add key / val to map
func (m *Collection[K, V, _]) Add(k K, v V) {
	m.mtx.Lock()
	m.m[k] = v
	m.mtx.Unlock()
}

// Remove key from map
func (m *Collection[K, _, _]) Remove(key K) {
	m.mtx.Lock()
	delete(m.m, key)

	m.mtx.Unlock()
}

// Mark key as deleted
func (m *Collection[K, _, _]) Delete(key K) {
	m.mtx.Lock()
	m.m[key].Del(true)

	m.mtx.Unlock()
}

// Mark key as not deleted
func (m *Collection[K, _, _]) UnDelete(key K) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	m.m[key].Del(false)
}

// Length of map
func (m *Collection[K, _, _]) Length() int {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	return len(m.m)
}

// All iterates over all elements of K
func (m *Collection[K, V, _]) All() iter.Seq2[K, V] {
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

// All iterates over elements with type T
func (m *Collection[K, V, T]) OfType(t T) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		m.mtx.RLock()
		defer m.mtx.RUnlock()

		for k, v := range m.m {
			if v.FilterType() == t {
				if !yield(k, v) {
					return
				}
			}
		}
	}
}

// NewCollection create new empty m: map[K]V
func NewCollection2[K MapKey, V MapValue[T], T TypeType]() *Collection[K, V, T] {
	c := &struct {
		Collection[K, V, T]
		_mtx sync.RWMutex
		_m   MapType[K, V, T]
	}{
		_mtx: sync.RWMutex{},
		_m:   make(MapType[K, V, T]),
	}

	x := &c.Collection
	x.mtx = &c._mtx
	x.m = c._m
	return x
	// return &Collection[K, V, T]{
	// 	mtx: &sync.RWMutex{},
	// 	m:   make(MapType[K, V, T]),
	// }
}
