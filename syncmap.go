package syncmap

import (
	"fmt"
	"iter"
	"strconv"
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

func (m *PointerMap[K]) Exists(key K) (ok bool) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	_, ok = m.m[key]

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

func (m *PointerMap[_]) Len() int {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	return len(m.m)
}

func (m *PointerMap[_]) LenStr() string {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	return strconv.Itoa(len(m.m))
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

// NewCollection creates new empty m: map[K]V
// Mid-Stack Inlined ?
// see https://dave.cheney.net/2020/05/02/mid-stack-inlining-in-go

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
	if m.mtx.TryRLock() {
		defer m.mtx.RUnlock()
		_, ok = m.m[key]
	} else {
		fmt.Println("TryRLock failed")
	}
	return
}

// Get val with key
func (m *Collection[K, V]) Get(key K) (val V) {
	if m.mtx.TryRLock() {
		defer m.mtx.RUnlock()
		val = m.m[key]
	} else {
		fmt.Println("TryRLock failed")
	}
	return
}

// Get val with key
func (m *Collection[K, V]) GetP(key K, v *V) (ok bool) {
	if m.mtx.TryRLock() {
		defer m.mtx.RUnlock()
		*v, ok = m.m[key]
	} else {
		fmt.Println("TryRLock failed")
	}
	return
}

// Get whole map
func (m *Collection[K, V]) GetAll() (val MapType[K, V]) {
	if m.mtx.TryRLock() {
		defer m.mtx.RUnlock()
		val = m.m
	} else {
		fmt.Printf("TryRLock failed for %T\n", *new(V))
	}
	return
}

// Set / Overwrite map from map
func (m *Collection[K, V]) Set(v map[K]V) {
	if m.mtx.TryLock() {
		defer m.mtx.Unlock()
		m.m = v
	} else {
		fmt.Printf("Set - TryLock failed for %T\n", *new(V))
	}
}

// Add key / val to map
func (m *Collection[K, V]) Add(k K, v V) {
	if m.mtx.TryLock() {
		defer m.mtx.Unlock()
		m.m[k] = v
	} else {
		fmt.Printf("Add - TryLock failed for %T\n", *new(V))
	}
}

// Remove key from map
func (m *Collection[K, _]) Remove(key K) {
	if m.mtx.TryLock() {
		defer m.mtx.Unlock()
		delete(m.m, key)
	} else {
		fmt.Println("Remove - TryLock failed")
	}

}

// Mark key as deleted
func (m *Collection[K, _]) Delete(key K) {
	if m.mtx.TryLock() {
		defer m.mtx.Unlock()
		m.m[key].Del(true)
	} else {
		fmt.Println("Delete - TryLock failed")
	}

}

// Mark key as not deleted
func (m *Collection[K, _]) UnDelete(key K) {
	if m.mtx.TryLock() {
		defer m.mtx.Unlock()
		m.m[key].Del(false)
	} else {
		fmt.Println("UnDelete - TryLock failed")
	}

}

// Len of map
func (m *Collection[_, _]) Len() int {
	if m.mtx.TryRLock() {
		defer m.mtx.RUnlock()
		return len(m.m)
	} else {
		fmt.Println("Len - TryRLock failed")
		return 0
	}

}

func (m *Collection[_, _]) LenStr() string {
	if m.mtx.TryRLock() {
		defer m.mtx.RUnlock()
		return strconv.Itoa(len(m.m))
	} else {
		fmt.Println("LenStr - TryRLock failed")
		return ""
	}

}

// All iterates over all elements of K
func (m *Collection[K, V]) All() iter.Seq2[K, V] {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	return func(yield func(K, V) bool) {

		for k, v := range m.m {
			if !yield(k, v) {
				return
			}
		}
	}
}
