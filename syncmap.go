package syncmap

import (
	"iter"
	"sync"
)

type PointerMap[K PointerType[T], V PointerMapBoolType, T comparable] struct {
	mtx sync.RWMutex
	m   map[K]V
}

type PointerMapBoolType interface {
	comparable
	struct{}
	// ID() string
	// Type() T
	// All()
}

type PointerType[T comparable] interface {
	comparable
	ID() string
	Type() T
	All()
}

func NewPointerMap[K PointerType[T], V PointerMapBoolType, T comparable]() PointerMap[K, V, T] {
	return PointerMap[K, V, T]{
		mtx: sync.RWMutex{},
		m:   make(map[K]V),
	}
}

func (m *PointerMap[K, V, _]) Exists(key K) bool {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	return m.m[key] == V{}
}

func (m *PointerMap[K, V, _]) Add(key K) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	m.m[key] = V{}
}

func (m *PointerMap[K, V, _]) Remove(key K) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	delete(m.m, key)
}

func (m *PointerMap[K, V, _]) Length() int {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	return len(m.m)
}

func (m *PointerMap[K, V, T]) OfType(t T) (k K) {
	for k := range m.m {
		if k.Type() == t {
			return k
		}
	}
	return k
}

// All is an iterator over the elements of s.
func (s *PointerMap[K, V, _]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k, v := range s.m {
			if !yield(k, v) {
				return
			}
		}
	}
}

func (m *PointerMap[K, V, _]) GetByID(id string) (k K) {
	for k := range m.m {
		if k.ID() == id {
			return k
		}
	}
	return k
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

type MapValue[T comparable] interface {
	comparable
	ID() string
	Type() T
	Delete(string)
	Add(string)
}

type MapKey[T comparable] interface {
	comparable
	ID() string
	Type() T
	Delete(string)
	Add(string)
}

type KeyValMap[K MapKey[T], V MapValue[T], T comparable] struct {
	mtx sync.RWMutex
	m   map[K]V
}

func (m *KeyValMap[K, V, T]) Get(key K) V {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	return m.m[key]
}

// func (m *PointerMap[K, V]) Set(key K, value V) {
// 	m.mtx.Lock()
// 	m.m[key] = value
// 	m.mtx.Unlock()
// }

func (m *KeyValMap[K, V, _]) Set(key K, value V) {
	m.mtx.Lock()
	m.m[key] = value
	m.mtx.Unlock()
}

func (m *KeyValMap[K, V, _]) Del(key K) {
	m.mtx.Lock()
	delete(m.m, key)
	m.mtx.Unlock()
}

func (m *KeyValMap[K, V, _]) Length() int {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	return len(m.m)
}

func (s *KeyValMap[K, V, _]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k, v := range s.m {
			if !yield(k, v) {
				return
			}
		}
	}
}

func (m *KeyValMap[K, V, T]) OfType(t T) (k K) {
	for k := range m.m {
		if k.Type() == t {
			return k
		}
	}
	return k
}
