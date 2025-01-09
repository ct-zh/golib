package ring

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRing(t *testing.T) {
	assert.Panics(t, func() {
		NewRing(0)
	})
}

func TestRingLess(t *testing.T) {
	r := NewRing(5)
	for i := 0; i < 3; i++ {
		r.Add(i)
	}
	assert.Equal(t, []interface{}{0, 1, 2}, r.Take())
}

func TestRingMore(t *testing.T) {
	r := NewRing(5)
	for i := 0; i < 12; i++ {
		r.Add(i)
	}
	assert.Equal(t, []interface{}{7, 8, 9, 10, 11}, r.Take())
}

func TestAdd(t *testing.T) {
	r := NewRing(5051)
	wg := sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for j := 0; j <= i; j++ {
				r.Add(j)
			}
		}(i)
	}
	wg.Wait()
	assert.Equal(t, 5050, len(r.Take()))
}

func BenchmarkAdd(b *testing.B) {
	r := NewRing(500)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for i := 0; i < b.N; i++ {
				r.Add(i)
			}
		}
	})
}
