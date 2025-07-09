package syncmap

import (
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
