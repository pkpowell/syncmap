package syncmap

import (
	"iter"
	"strconv"
	"sync"
)

// ///////////////////////////
// Collection
// ///////////////////////////

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
	m   map[K]V
}

// NewCollection creates new empty m: map[K]V
// Mid-Stack Inlined ?
// see https://dave.cheney.net/2020/05/02/mid-stack-inlining-in-go

func NewCollection[K MapKey, V MapValue]() (c *Collection[K, V]) {
	// var c Collection[K, V]
	return newCollection(c)
}

func newCollection[K MapKey, V MapValue](c *Collection[K, V]) *Collection[K, V] {
	c.mtx = &sync.RWMutex{}
	c.m = make(map[K]V)
	return c
}

// Exists check if key exists
func (c *Collection[K, _]) Exists(key K) (ok bool) {
	c.mtx.RLock()
	defer c.mtx.RUnlock()

	_, ok = c.m[key]
	return ok
}

// Get val with key
func (c *Collection[K, V]) Get(key K) (val V, ok bool) {
	c.mtx.RLock()
	defer c.mtx.RUnlock()

	val, ok = c.m[key]
	return val, ok
}

// Get val with key and write to v
func (c *Collection[K, V]) GetP(key K, val *V) (ok bool) {
	c.mtx.RLock()
	defer c.mtx.RUnlock()

	*val, ok = c.m[key]
	return ok
}

// // Get whole map - use ToMap
// func (c *Collection[K, V]) GetAll() (val *map[K]V) {
// 	c.mtx.RLock()
// 	defer c.mtx.RUnlock()

// 	val = &c.m
// 	return val
// }

// Get whole map - replaces GetAll
func (c *Collection[K, V]) ToMap() map[K]V {
	c.mtx.RLock()
	defer c.mtx.RUnlock()

	return c.m
	// val = c.m
	// return val
}

// Set / Overwrite map from map
func (c *Collection[K, V]) Set(v map[K]V) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	c.m = v
}

// Add key / val to map
func (c *Collection[K, V]) Add(k K, v V) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	c.m[k] = v
}

// Add key / val to map, returns 'updated'. Doesn't really work....
func (c *Collection[K, V]) AddCompare(k K, v V) (updated bool) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	updated = c.m[k] == v

	c.m[k] = v
	return
}

// Remove key from map
func (c *Collection[K, _]) Remove(key K) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	delete(c.m, key)
}

// Mark key as deleted
func (c *Collection[K, _]) Delete(key K) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	c.m[key].Del(true)
}

// Mark key as not deleted
func (c *Collection[K, _]) UnDelete(key K) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	c.m[key].Del(false)
}

// Len of map
func (c *Collection[_, _]) Len() int {
	c.mtx.RLock()
	defer c.mtx.RUnlock()

	return len(c.m)
}

func (c *Collection[_, _]) LenStr() string {
	c.mtx.RLock()
	defer c.mtx.RUnlock()

	return strconv.Itoa(len(c.m))
}

// All iterates over all elements of K
// Replaced by Iter()
func (c *Collection[K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		c.mtx.RLock()
		defer c.mtx.RUnlock()

		for k, v := range c.m {
			if !yield(k, v) {
				return
			}
		}
	}
}

// Iter iterates over all elements of K - replaces .All()
func (c *Collection[K, V]) Iter() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		c.mtx.RLock()
		defer c.mtx.RUnlock()

		for k, v := range c.m {
			if !yield(k, v) {
				return
			}
		}
	}
}
