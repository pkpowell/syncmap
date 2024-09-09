package syncmap

import (
	"iter"
	"sync"
)

type TypeType interface {
	comparable
}

/////////////////////////////
// PointerMap
/////////////////////////////

type PointerType[T TypeType] interface {
	comparable
	ID() string
	Type() T
}

type PointerBoolType interface {
	struct{}
}

type PointerMap[K PointerType[T], V PointerBoolType, T TypeType] struct {
	mtx sync.RWMutex
	m   map[K]V
}

func NewPointerMap[K PointerType[T], T TypeType]() PointerMap[K, struct{}, T] {
	return PointerMap[K, struct{}, T]{
		mtx: sync.RWMutex{},
		m:   make(map[K]struct{}),
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

// OfType is an iterator over the elements of s with type t.
func (s *PointerMap[K, V, T]) OfType(t T) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k, v := range s.m {
			if k.Type() == t {
				if !yield(k, v) {
					return
				}
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

/////////////////////////////
// KeyValMap
/////////////////////////////

type MapValue[T TypeType] interface {
	comparable
	ID() string
	Type() T
}

type MapKey interface {
	comparable
}

type KeyValMap[K MapKey, V MapValue[T], T TypeType] struct {
	mtx sync.RWMutex
	m   map[K]V
}

func (m *KeyValMap[K, V, T]) Get(key K) V {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	return m.m[key]
}

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

func (s *KeyValMap[K, V, T]) OfType(t T) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k, v := range s.m {
			if v.Type() == t {
				if !yield(k, v) {
					return
				}
			}
		}
	}
}
