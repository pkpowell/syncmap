package syncmap

import (
	"fmt"
	"testing"
)

type TestType struct {
	Field string
	Array []int
}

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

var (
	p  = NewPointerMap[*TestType]()
	pc = NewCollection[struct{}, *TestType, string]()
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

func BenchmarkCollectionAdd(b *testing.B) {
	for i := 0; i < b.N; i++ {
		c.Add(fmt.Sprintf("test-%d", i), &TestType{
			Field: fmt.Sprintf("test-%d", i),
			Array: []int{i, 2, 3},
		})
	}
}
