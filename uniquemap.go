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
type UniqueMapType[K MapKey, V MapValue] map[K]unique.Handle[V]

type UniqueCollection[K MapKey, V MapValue] struct {
	mtx *sync.RWMutex
	m   UniqueMapType[K, V]
}

// NewUniqueCollection creates new empty m: map[K]V
// Mid-Stack Inlined ?
// see https://dave.cheney.net/2020/05/02/mid-stack-inlining-in-go
func NewUniqueCollection[K MapKey, V MapValue]() *UniqueCollection[K, V] {
	var c UniqueCollection[K, V]
	return newUniqueCollection(&c)
}

func newUniqueCollection[K MapKey, V MapValue](c *UniqueCollection[K, V]) *UniqueCollection[K, V] {
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
func (m *UniqueCollection[K, V]) Get(key K) (val V) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	val = m.m[key].Value()
	return
}

// Get whole map
func (m *UniqueCollection[K, V]) GetAll() *MapType[K, V] {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	val := make(MapType[K, V])

	for i, u := range m.m {
		val[i] = u.Value()
	}

	return &val
}

// Overwrite map from map
func (m *UniqueCollection[K, V]) Overwrite(d MapType[K, V]) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	clear(m.m)

	for k, v := range d {
		m.m[k] = unique.Make(v)
	}
}

// merge data from map
func (m *UniqueCollection[K, V]) Merge(d MapType[K, V]) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	for k, v := range d {
		m.m[k] = unique.Make(v)
	}
}

// Add key / val to map, return true if value changed
func (m *UniqueCollection[K, V]) Add(k K, v V) (updated bool) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	new := unique.Make(v)
	updated = new != m.m[k]

	m.m[k] = new
	return
}

// Remove key from map
func (m *UniqueCollection[K, _]) Remove(key K) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	delete(m.m, key)
}

// // Mark key as deleted
// func (m *UniqueCollection[K, _]) Delete(key K) {
// 	m.mtx.Lock()
// 	defer m.mtx.Unlock()

// 	m.m[key].Del(true)
// }

// // Mark key as not deleted
// func (m *UniqueCollection[K, _]) UnDelete(key K) {
// 	m.mtx.Lock()
// 	defer m.mtx.Unlock()

// 	m.m[key].Del(false)
// }

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
