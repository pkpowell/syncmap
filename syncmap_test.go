package syncmap

import (
	"fmt"
	"testing"
)

type TestType struct {
	Field string
	Array []int
}

type TestBool Bool

func (t *TestType) GetID() string {
	return t.Field
}

func (t *TestType) FilterType() string {
	return t.Field
}
func (t *TestType) IDX() string {
	return t.Field
}

func (t *TestType) Del(bool) {}

func (t *TestBool) FilterType() {}
func (t *TestBool) IDX()        {}

// func (t *TestBool) Type() {}
// func (t *TestBool) IDX()  {}

// func (t *TestBool) Del(bool) {}

var (
	p = NewPointerMap[*TestType]()
	// pc = NewCollection[*TestType, *TestBool, string]()
	c = NewCollection[string, *TestType, string]()
)

func BenchmarkPointerMapAdd(b *testing.B) {
	for i := 0; i < b.N; i++ {
		p.Add(&TestType{
			Field: fmt.Sprintf("test-%d", i),
			Array: []int{i, 2, 3},
		})
	}
}

// func BenchmarkCollBoolAdd(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		pc.Add(&TestType{
// 			Field: fmt.Sprintf("test-%d", i),
// 			Array: []int{i, 2, 3},
// 		}, &TestBool{})
// 	}
// }

func BenchmarkCollectionAdd(b *testing.B) {
	for i := 0; i < b.N; i++ {
		c.Add(fmt.Sprintf("test-%d", i), &TestType{
			Field: fmt.Sprintf("test-%d", i),
			Array: []int{i, 2, 3},
		})
	}
}
func BenchmarkCollectionGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		c.Get(fmt.Sprintf("test-%d", i))
	}
}
func BenchmarkCollectionGetAll(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = c.All()
	}
}
func BenchmarkGet(b *testing.B) {
	b.Run("add", BenchmarkCollectionAdd)
	b.Run("get", BenchmarkCollectionGet)
	b.Run("getall", BenchmarkCollectionGetAll)
}
