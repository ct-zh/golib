package tracktime

import (
	"testing"
	"time"
)

func TestDo(t *testing.T) {
	defer Do("test_do")()

	time.Sleep(time.Second * 2)
}
