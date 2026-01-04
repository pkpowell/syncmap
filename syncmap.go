package syncmap

import (
	"iter"
	"strconv"
	"sync"
)

// ///////////////////////////
// Collection
// ///////////////////////////
// type MapType[K MapKey, V MapValue] map[K]V

type MapValue interface {
	comparable
	GetID() string
	Del(bool)
}

type MapKey interface {
	comparable
}

type Collection[K MapKey, V MapValue] struct {
	mtx *sync.RWMutex
	m   map[K]V
}

// NewCollection creates new empty m: map[K]V
// Mid-Stack Inlined ?
// see https://dave.cheney.net/2020/05/02/mid-stack-inlining-in-go

func NewCollection[K MapKey, V MapValue]() *Collection[K, V] {
	var c Collection[K, V]
	return newCollection(&c)
}

func newCollection[K MapKey, V MapValue](c *Collection[K, V]) *Collection[K, V] {
	c.mtx = &sync.RWMutex{}
	c.m = make(map[K]V)
	return c
}

// Exists check if key exists
func (m *Collection[K, _]) Exists(key K) (ok bool) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	_, ok = m.m[key]
	return
}

// Get val with key
func (m *Collection[K, V]) Get(key K) (val V, ok bool) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	val, ok = m.m[key]
	return
}

// Get val with key and write to v
func (m *Collection[K, V]) GetP(key K, v *V) (ok bool) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	*v, ok = m.m[key]
	return
}

// Get whole map - use ToMap
func (m *Collection[K, V]) GetAll() (val *map[K]V) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	val = &m.m
	return
}

// Get whole map - replaces GetAll
func (m *Collection[K, V]) ToMap() (val map[K]V) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	val = m.m
	return
}

// Set / Overwrite map from map
func (m *Collection[K, V]) Set(v map[K]V) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	m.m = v
}

// Add key / val to map
func (m *Collection[K, V]) Add(k K, v V) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	m.m[k] = v
}

// Add key / val to map, returns 'updated'. Doesn't really work....
func (m *Collection[K, V]) AddCompare(k K, v V) (updated bool) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	updated = m.m[k] == v

	m.m[k] = v
	return
}

// Remove key from map
func (m *Collection[K, _]) Remove(key K) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	delete(m.m, key)
}

// Mark key as deleted
func (m *Collection[K, _]) Delete(key K) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	m.m[key].Del(true)
}

// Mark key as not deleted
func (m *Collection[K, _]) UnDelete(key K) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	m.m[key].Del(false)
}

// Len of map
func (m *Collection[_, _]) Len() int {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	return len(m.m)
}

func (m *Collection[_, _]) LenStr() string {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	return strconv.Itoa(len(m.m))
}

// All iterates over all elements of K
// Replaced by Iter()
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

// Iter iterates over all elements of K - replaces .All()
func (m *Collection[K, V]) Iter() iter.Seq2[K, V] {
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
