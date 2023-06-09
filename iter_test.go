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
func compareSlice[T Number](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for i, va := range a {
		if b[i] != va {
			return false
		}
	}
	return true
}
func TestMap(t *testing.T) {
	if ret, err := Map([]int{1, 2, 3}, func(t int) (int, error) {
		return t * 2, nil
	}); err != nil {
		t.Fail()
	} else if !compareSlice(ret, []int{2, 4, 6}) {
		t.Fail()
	}
}

func TestApply(t *testing.T) {
	x := []int{1}
	if ret, err := Apply(x, func(val []int) ([]int, error) {
		val[0] *= 2
		return val, nil
	}, func(val []int) ([]int, error) {
		val[0] += 1
		return val, nil
	}); err != nil {
		t.Log("apply error", err)
		t.Fail()
	} else if ret[0] != 3 {
		t.Log(ret)
		t.Fail()
	}
}
func TestWithFilter(t *testing.T) {
	it := WithFilter(SimpleIter([]int{1, 2, 3, 4}), func(i int) bool {
		return i%2 == 1
	})
	testIterFunc(t, it, []int{1, 3})
}
func TestFilter(t *testing.T) {
	ret := Filter([]int{1, 2, 3, 4}, func(i int) bool {
		return i%2 == 1
	})
	if !compareSlice(ret, []int{1, 3}) {
		t.Fail()
	}
}
func TestAny(t *testing.T) {
	if !Any([]int{1, 2, 4, 5, 1024}, func(i int) bool {
		return i > 100
	}) {
		t.Fail()
	}
}
func TestAll(t *testing.T) {
	if !All([]int{1, 3, 5, 7}, func(i int) bool {
		return i%2 == 1
	}) {
		t.Fail()
	}
}
func TestSillyRange(t *testing.T) {
	it := SillyRange[int](func() (int, error) {
		return 0, nil
	}, func(a int) (int, error) {
		if a >= 4 {
			return -1, Done
		}
		a += 1
		return a, nil
	})
	testIterFunc(t, it, []int{0, 1, 2, 3, 4})
}
func TestReduce(t *testing.T) {
	it := SillyRange[int](func() (int, error) {
		return 1, nil
	}, func(a int) (int, error) {
		if a >= 100 {
			return -1, Done
		}
		a += 1
		return a, nil
	})
	ret, err := Reduce[int](it, func(a, b int) (int, error) {
		return a + b, nil
	}, nil)
	if ret != 5050 || err != nil {
		t.Logf("ret= %d err = %s", ret, err)
		t.Fail()
	}
}
func TestFindIter(t *testing.T) {
	ret, err := FindIter(SimpleIter([]int{1, 2, 3, 4, 5}), func(i int) (bool, error) {
		return i > 1 && i%2 == 1, nil
	})
	if ret != 3 && err != nil {
		t.Fail()
	}
}
func TestNewEnumerate(t *testing.T) {
	it := NewEnumerate(SimpleIter([]int{1, 2, 3, 4}))
	pairs := [][]int{{0, 1}, {1, 2}, {2, 3}, {3, 4}}
	for _, pair := range pairs {
		c, v, err := it.Next()
		if err != nil && c != pair[0] && v != pair[1] {
			t.Fail()
		}
	}
	if _, _, err := it.Next(); err != Done {
		t.Fail()
	}
}
