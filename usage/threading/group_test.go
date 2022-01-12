package threading

import (
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGoroutineGroup_Run(t *testing.T) {
	var count int32

	g := NewRoutineGroup()
	for i := 0; i < 5; i++ {
		g.Run(func() {
			atomic.AddInt32(&count, 1)
		})
	}

	g.Wait()
	assert.Equal(t, int32(5), count)
}

func TestGoroutineGroup_RunSafe(t *testing.T) {
	var count int32
	g := NewRoutineGroup()
	var once sync.Once
	for i := 0; i < 5; i++ {
		g.RunSafe(func() {
			once.Do(func() {
				panic("")
			})
			atomic.AddInt32(&count, 1)
		})
	}
	g.Wait()
	assert.Equal(t, int32(4), count)
}
