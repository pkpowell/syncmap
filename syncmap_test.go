package syncmap

import (
	"fmt"
	"testing"
)

type TestType struct {
	Field string
	Array []int
}

// type TestBool Bool

func (t *TestType) GetID() string {
	return t.Field
}

func (t *TestType) Type() string {
	return t.Field
}
func (t *TestType) IDX() string {
	return t.Field
}

func (t *TestType) Del(bool) {}

// func (t *Bool) GetID()   {}

// func (t *TestBool) Type() {}
// func (t *TestBool) IDX()  {}

// func (t *TestBool) Del(bool) {}

var (
	p  = NewPointerMap[*TestType]()
	pc = NewCollection[*TestType, *Bool, string]()
	c  = NewCollection[string, *TestType, string]()
)

func BenchmarkPointerMapAdd(b *testing.B) {
	for i := 0; i < b.N; i++ {
		p.Add(&TestType{
			Field: fmt.Sprintf("test-%d", i),
			Array: []int{i, 2, 3},
		})
	}
}
func BenchmarkCollBoolAdd(b *testing.B) {
	for i := 0; i < b.N; i++ {
		pc.Add(&TestType{
			Field: fmt.Sprintf("test-%d", i),
			Array: []int{i, 2, 3},
		}, &Bool{})
	}
}

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
func BenchmarkGet(b *testing.B) {
	b.Run("add", BenchmarkCollectionAdd)
	b.Run("get", BenchmarkCollectionGet)
}
