package silly_kits

import (
	"testing"
)

func simpleIterFunc(v int) ([]int, error) {
	if v > 1 {
		return nil, Done
	}
	return []int{1, 2, 3}, nil
}

type Number interface {
	~int | ~uint | ~float32
}

func testIterFunc[T Number](t *testing.T, it Iterator[T], expected []T) {
	cursor := 0
	for {
		if val, err := it.Next(); err != nil {
			break
		} else if val != expected[cursor] {
			t.Fail()
		}
		cursor += 1
	}
	if len(expected) != cursor {
		t.Logf("curor is %d but expected length is %d", cursor, len(expected))
		t.Fail()
	}
}
func TestLazy(t *testing.T) {
	it := Lazy(simpleIterFunc)
	testIterFunc(t, it, []int{1, 2, 3, 1, 2, 3})
}
func TestWrapper(t *testing.T) {
	it := Lazy(simpleIterFunc)
	w := Wrapper(it, func(t int) (float32, error) {
		return float32(t), nil
	})
	testIterFunc(t, w, []float32{1, 2, 3, 1, 2, 3})
}
func TestSimpleIter(t *testing.T) {
	ret := []int{1, 2, 3, 4}
	it := SimpleIter(ret)
	testIterFunc(t, it, ret)
}
func TestSimpleChain(t *testing.T) {
	it := SimpleChain(SimpleIter([]int{1, 2, 3, 4}), SimpleIter([]int{5, 6, 7, 8, 9}))
	testIterFunc(t, it, []int{1, 2, 3, 4, 5, 6, 7, 8, 9})
}

func TestChain(t *testing.T) {
	c := Chain(SimpleIter([]Iterator[int]{SimpleIter([]int{1, 2, 3, 4}), SimpleIter([]int{5, 6, 7, 8, 9})}))
	testIterFunc(t, c, []int{1, 2, 3, 4, 5, 6, 7, 8, 9})
}
func TestEmpty(t *testing.T) {
	testIterFunc(t, Empty[int](), []int{})
}
