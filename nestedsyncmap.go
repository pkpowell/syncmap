package syncmap

import (
	"strconv"
	"sync"
)

// ///////////////////////////
// Nested Collection
// ///////////////////////////

type NestedCollection[K MapKey, I MapKey, V MapValue] struct {
	mtx *sync.RWMutex
	m   map[K]map[I]V
}

// NewNestedCollection creates new empty m: map[K]map[I]V
func NewNestedCollection[K MapKey, I MapKey, V MapValue]() *NestedCollection[K, I, V] {
	var c NestedCollection[K, I, V]
	return newNestedCollection(&c)
}

func newNestedCollection[K MapKey, I MapKey, V MapValue](c *NestedCollection[K, I, V]) *NestedCollection[K, I, V] {
	c.mtx = &sync.RWMutex{}
	c.m = make(map[K]map[I]V)
	return c
}

// Exists check if key exists
func (m *NestedCollection[K, I, _]) Exists(k K) (ok bool) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	_, ok = m.m[k]
	return
}

// Get val with key
func (m *NestedCollection[K, I, V]) Get(k K, i I) (val V, ok bool) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	val, ok = m.m[k][i]
	return
}

// Get val with key and write to v
func (m *NestedCollection[K, I, V]) GetP(k K, i I, v *V) (ok bool) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	*v, ok = m.m[k][i]
	return
}

// Get whole map - use ToMap
func (m *NestedCollection[K, I, V]) GetAll() (val *map[K]map[I]V) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	val = &m.m
	return
}

// Get whole map - replaces GetAll
func (m *NestedCollection[K, I, V]) ToMap() (val map[K]map[I]V) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	val = m.m
	return
}

// Set / Overwrite map from map
func (m *NestedCollection[K, I, V]) Set(v map[K]map[I]V) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	m.m = v
}

// Add key / val to map
func (m *NestedCollection[K, I, V]) Add(k K, i I, v V) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	m.m[k][i] = v
}

// Add key / val to map, returns 'updated'. Doesn't really work....
func (m *NestedCollection[K, I, V]) AddCompare(k K, i I, v V) (updated bool) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	updated = m.m[k][i] == v

	m.m[k][i] = v
	return
}

// Remove key from map
func (m *NestedCollection[K, I, _]) Remove(k *K, i *I) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	if k == nil {
		return
	}

	if i == nil {
		delete(m.m, *k)
		return
	}

	delete(m.m[*k], *i)
}

// Mark key as deleted
func (m *NestedCollection[K, I, _]) Delete(k *K, i *I) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	if k == nil {
		return
	}

	m.m[*k][*i].Del(true)
}

// Mark key as not deleted
func (m *NestedCollection[K, I, _]) UnDelete(k *K, i *I) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	if k == nil {
		return
	}
	m.m[*k][*i].Del(false)
}

// Len of map
func (m *NestedCollection[_, _, _]) Len() int {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	return len(m.m)
}

func (m *NestedCollection[_, _, _]) LenStr() string {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	return strconv.Itoa(len(m.m))
}

// // All iterates over all elements of K
// // Replaced by Iter()
// func (m *NestedCollection[K, I, V]) AllOuter() iter.Seq2[K, map[I]V] {
// 	return func(yield func(K, I, V) bool) {
// 		m.mtx.RLock()
// 		defer m.mtx.RUnlock()

// 		for k, v := range m.m {
// 			if !yield(k, v) {
// 				return
// 			}
// 		}
// 	}
// }

// // Iter iterates over all elements of K - replaces .All()
// func (m *NestedCollection[K, I, V]) Iter() iter.Seq2[K, V] {
// 	return func(yield func(K, V) bool) {
// 		m.mtx.RLock()
// 		defer m.mtx.RUnlock()

// 		for k, v := range m.m {
// 			if !yield(k, v) {
// 				return
// 			}
// 		}
// 	}
// }
