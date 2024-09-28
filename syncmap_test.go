package syncmap

import (
	"fmt"
	"testing"
)

type TestType struct {
	Field string
	Array []int
}

type TestBool struct{}

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
			// Field: fmt.Sprintf("test-%d", i),
			// Array:  []int{i, 2, 3},
		})
	}
}

//	func BenchmarkCollBoolAdd(b *testing.B) {
//		for i := 0; i < b.N; i++ {
//			pc.Add(&TestType{
//				Field: fmt.Sprintf("test-%d", i),
//				Array: []int{i, 2, 3},
//			}, &TestBool{})
//		}
//	}
var s = &struct {
	TestType
	_field string
	_array []int
}{}
var t = &s.TestType

// func BenchmarkCollectionAdd(b *testing.B) {

//		t.Field = &s._field
//		t.Array = &s._array
//		for i := 0; i < b.N; i++ {
//			*t.Field = fmt.Sprintf("test-%d", i)
//			*t.Array = []int{i, 2, 3}
//			c.Add(*t.Field, t)
//		}
//	}
func BenchmarkCollectionGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		c.Get(fmt.Sprintf("test-%d", i))
	}
}
func BenchmarkCollectionGetP(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var d *TestType
		c.GetP(fmt.Sprintf("test-%d", i), &d)

		//fmt.Println("d: ", d)
	}
}
func BenchmarkCollectionGetAll(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, _ = range c.All() {

		}
	}
}
func BenchmarkGet(b *testing.B) {
	// b.Run("add", BenchmarkCollectionAdd)
	b.Run("get", BenchmarkCollectionGet)
	b.Run("getp", BenchmarkCollectionGetP)
	b.Run("getall", BenchmarkCollectionGetAll)
}

// func TestPutGet(b *testing.T) {
// 	var tests = map[string]struct {
//         a, b *TestType
//         want *TestType
//     }{
//     "one": {
// 		a:   &TestType{
// 			Field: "test-1",
// 			Array:  []int{1, 2, 3},
// 		},
// 		b: &TestType{
// 			Field: "test-2",
// 			Array:  []int{4, 2, 3},
// 		},
// 		want: &TestType{
// 			Field: "test-1",
// 			Array:  []int{1, 2, 3},
// 		},
// 	},
// }

//     for name, tt := range tests {
//         testname := fmt.Sprintf("%d,%d", tt.a, tt.b)
//         t.Run(name, func(t *testing.T) {
//             ans := IntMin(tt.a, tt.b)
//             if ans != tt.want {
//                 t.Errorf("got %d, want %d", ans, tt.want)
//             }
//         })
//     }
// }

func TestCollectionPut(t *testing.T) {

}
func TestCollectionGetP(t *testing.T) {

}
