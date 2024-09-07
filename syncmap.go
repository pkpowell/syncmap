package syncmap

import (
	"iter"
	"sync"
)

type PointerMapType interface {
	comparable
	ID() string
	Length() int
	Type() any
}

type PointerMap[K PointerMapType, V struct{}] struct {
	mtx sync.RWMutex
	m   map[K]V
}

func NewPointerMap[K PointerMapType]() PointerMap[K, struct{}] {
	return PointerMap[K, struct{}]{
		mtx: sync.RWMutex{},
		m:   make(map[K]struct{}),
	}
}

func (m *PointerMap[K, V]) Exists(key K) bool {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	return m.m[key] == V{}
}

func (m *PointerMap[K, V]) Add(key K) {
	m.mtx.Lock()
	defer m.mtx.RUnlock()

	m.m[key] = struct{}{}
}

func (m *PointerMap[K, V]) Remove(key K) {
	m.mtx.Lock()
	defer m.mtx.RUnlock()

	delete(m.m, key)
}

func (m *PointerMap[K, V]) GetByID(id string) (k K) {
	for k := range m.m {
		if k.ID() == id {
			return k
		}
	}
	return k
}

func (m *PointerMap[K, V]) Type(t string) (k K) {
	for k := range m.m {
		if k.Type() == t {
			return k
		}
	}
	return k
}

func (m *PointerMap[K, V]) Length() int {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	return len(m.m)
}

// func (m *PointerMap[K, V]) Length() int {
// 	m.mtx.Lock()
// 	defer m.mtx.Unlock()

// 	return len(m.m)
// }

// All is an iterator over the elements of s.
func (s *PointerMap[K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k, v := range s.m {
			if !yield(k, v) {
				return
			}
		}
	}
}
