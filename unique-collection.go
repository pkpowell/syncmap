package syncmap

import (
	"iter"
	"strconv"
	"sync"
	"unique"
)

// ///////////////////////////
// Unique Collection
// ///////////////////////////
type UniqueMapType[K UniqueMapKey, V UniqueMapValue] map[K]unique.Handle[V]

type UniqueMapValue interface {
	comparable
	GetID() string
	Del(bool)
}

type UniqueMapKey interface {
	comparable
}

type UniqueCollection[K UniqueMapKey, V UniqueMapValue] struct {
	mtx *sync.RWMutex
	m   UniqueMapType[K, V]
}

// NewCollection creates new empty m: map[K]V
// Mid-Stack Inlined ?
// see https://dave.cheney.net/2020/05/02/mid-stack-inlining-in-go

func NewUniqueCollection[K UniqueMapKey, V UniqueMapValue]() *UniqueCollection[K, V] {
	var c UniqueCollection[K, V]
	return newUniqueCollection(&c)
}

func newUniqueCollection[K UniqueMapKey, V UniqueMapValue](c *UniqueCollection[K, V]) *UniqueCollection[K, V] {
	c.mtx = &sync.RWMutex{}
	c.m = make(UniqueMapType[K, V])
	return c
}

// Exists check if key exists
func (m *UniqueCollection[K, _]) Exists(key K) (ok bool) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	_, ok = m.m[key]
	return
}

// Get val with key
func (m *UniqueCollection[K, V]) Get(key K) (val unique.Handle[V]) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	val = m.m[key]
	return
}

// Get val with key
func (m *UniqueCollection[K, V]) GetP(key K, v *unique.Handle[V]) (ok bool) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	*v, ok = m.m[key]
	return
}

// Get whole map
func (m *UniqueCollection[K, V]) GetAll() (val *UniqueMapType[K, V]) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	val = &m.m
	return
}

// Set / Overwrite map from map
func (m *UniqueCollection[K, V]) Set(v map[K]unique.Handle[V]) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	m.m = v
}

// Add key / val to map
func (m *UniqueCollection[K, V]) Add(k K, v unique.Handle[V]) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	m.m[k] = v
}

// Remove key from map
func (m *UniqueCollection[K, _]) Remove(key K) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	delete(m.m, key)
}

// Mark key as deleted
func (m *UniqueCollection[K, _]) Delete(key K) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	m.m[key].Del(true)
}

// Mark key as not deleted
func (m *UniqueCollection[K, _]) UnDelete(key K) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	m.m[key].Del(false)
}

// Len of map
func (m *UniqueCollection[_, _]) Len() int {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	return len(m.m)
}

func (m *UniqueCollection[_, _]) LenStr() string {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	return strconv.Itoa(len(m.m))
}

// All iterates over all elements of K
func (m *UniqueCollection[K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		m.mtx.RLock()
		defer m.mtx.RUnlock()

		for k, v := range m.m {
			if !yield(k, v.Value()) {
				return
			}
		}
	}
}
