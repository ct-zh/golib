package errorx

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"sync"
	"sync/atomic"
	"testing"
)

var targetErr = errors.New("test error")

func TestAtomic(t *testing.T) {
	var e AtomicError
	assert.Nil(t, e.Load())
	e.Set(targetErr)
	assert.Equal(t, targetErr, e.Load())
}

func TestChain(t *testing.T) {
	assert.Nil(t, Chain(func() error {
		return nil
	}, func() error {
		return nil
	}))
	assert.Equal(t, targetErr, Chain(func() error {
		return nil
	}, func() error {
		return targetErr
	}))
	assert.Equal(t, targetErr, Chain(func() error {
		return targetErr
	}, func() error {
		return nil
	}))
}

func BenchmarkAtomic(b *testing.B) {
	wg := sync.WaitGroup{}
	b.Run("test load", func(b *testing.B) {
		// 先开个协程无限set，再开bench一直load
		var a AtomicError
		var done int32
		go func() {
			for {
				if atomic.LoadInt32(&done) == 1 {
					break
				}
				wg.Add(1)
				go func() {
					a.Set(targetErr)
					wg.Done()
				}()
			}
		}()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = a.Load()
		}
		b.StopTimer()
		atomic.StoreInt32(&done, 1)
		wg.Wait()
	})
	b.Run("test load", func(b *testing.B) {
		// 先开个协程无限set，再开bench一直load
		var a AtomicError
		var done int32
		go func() {
			for {
				if atomic.LoadInt32(&done) == 1 {
					break
				}
				wg.Add(1)
				go func() {
					_ = a.Load()
					wg.Done()
				}()
			}
		}()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			a.Set(targetErr)
		}
		b.StopTimer()
		atomic.StoreInt32(&done, 1)
		wg.Wait()
	})

}
