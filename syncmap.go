package syncmap

import (
	"iter"
	"sync"
)

/////////////////////////////
// PointerMap
/////////////////////////////

type PointerType interface {
	comparable
	GetID() string
}

type PointerMap[K PointerType] struct {
	mtx *sync.RWMutex
	m   map[K]struct{}
}

// NewPointerMap init pointer map with type field T
func NewPointerMap[K PointerType]() *PointerMap[K] {
	var p PointerMap[K]
	return newPointerMap(&p)
}

func newPointerMap[K PointerType](p *PointerMap[K]) *PointerMap[K] {
	p.mtx = &sync.RWMutex{}
	p.m = make(map[K]struct{})
	return p
}

func (m *PointerMap[K]) Exists(key K) bool {
	m.mtx.RLock()
	_, ok := m.m[key]
	m.mtx.RUnlock()

	return ok
}

func (m *PointerMap[K]) Add(key K) {
	m.mtx.Lock()
	m.m[key] = struct{}{}
	m.mtx.Unlock()
}

func (m *PointerMap[K]) Remove(key K) {
	m.mtx.Lock()
	delete(m.m, key)
	m.mtx.Unlock()
}

func (m *PointerMap[_]) Length() (l int) {
	m.mtx.RLock()
	l = len(m.m)
	m.mtx.RUnlock()

	return
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

func (t *Bool) Del(bool) {}

// NewCollection creates new empty m: map[K]V
func NewCollection[K MapKey, V MapValue]() *Collection[K, V] {
	var c Collection[K, V]
	return newCollection(&c)
}

func newCollection[K MapKey, V MapValue](c *Collection[K, V]) *Collection[K, V] {
	c.mtx = &sync.RWMutex{}
	c.m = make(MapType[K, V])
	return c
}

// Exists check if key exists
func (m *Collection[K, _]) Exists(key K) (ok bool) {
	m.mtx.RLock()
	_, ok = m.m[key]
	m.mtx.RUnlock()

	return
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
	*v, ok = m.m[key]
	m.mtx.RUnlock()
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
func (m *Collection[K, _]) Length() (l int) {
	m.mtx.RLock()
	l = len(m.m)
	m.mtx.RUnlock()

	return
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
