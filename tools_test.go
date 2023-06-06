package silly_kits

import (
	"errors"
	"testing"
	"time"
)

func TestFind(t *testing.T) {
	idx, val, err := Find([]int{0, 1, 2, 3, 4, 5}, func(i int) bool {
		return i%2 == 1
	})
	if idx != 1 || val != 1 || err != nil {
		t.Logf("idx=%d val= %d err=%s", idx, val, err)
		t.Fail()
	}
}
func TestForEach(t *testing.T) {
	if err := ForEach([]int{1, 2, 3, 10}, func(i int) error {
		if i > 5 {
			return errors.New("unknown")
		}
		return nil
	}); err == nil {
		t.Fail()
	}
}
func TestRetry(t *testing.T) {
	var c = 0
	val, err := Retry(func() (int, error) {
		c += 1
		if c < 3 {
			return -1, errors.New("unknown")
		}
		return 1023, nil
	}, 3, time.Second*3)
	if err != nil || val != 1023 {
		t.Fail()
	}
}
