package syncmap

import (
	"iter"
	"sync"
)

type PointerMap[K PointerType, V struct{}, T any] struct {
	mtx sync.RWMutex
	m   map[K]V
}

type PointerMapTypes interface {
	any
}
type PointerType interface {
	comparable
	// ObjType()
	All()
}

func NewPointerMap[K PointerType, T PointerMapTypes]() PointerMap[K, struct{}, T] {
	return PointerMap[K, struct{}, T]{
		mtx: sync.RWMutex{},
		m:   make(map[K]struct{}),
	}
}

func (m *PointerMap[K, V, T]) Exists(key K) bool {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	return m.m[key] == V{}
}

func (m *PointerMap[K, V, T]) Add(key K) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	m.m[key] = struct{}{}
}

func (m *PointerMap[K, V, T]) Remove(key K) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	delete(m.m, key)
}

func (m *PointerMap[K, V, T]) Length() int {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	return len(m.m)
}

// All is an iterator over the elements of s.
func (s *PointerMap[K, V, T]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k, v := range s.m {
			if !yield(k, v) {
				return
			}
		}
	}
}

// // All is an iterator over the elements of s.
// func (s *PointerMap[K, V, T]) ByType(t T) iter.Seq2[K, V] {
// 	return func(yield func(K, V) bool) {
// 		for k, v := range s.m {
// 			if k.ObjType(t) == t {
// 				if !yield(k, v) {
// 					return
// 				}
// 			}
// 		}
// 	}
// }

// func(k *PointerType)ObjType(){

// }
