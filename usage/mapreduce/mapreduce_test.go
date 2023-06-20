package mapreduce

import (
	"errors"
	"io/ioutil"
	"log"
	"runtime"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
)

var errDummy = errors.New("dummy")

func init() {
	log.SetOutput(ioutil.Discard)
}

func TestFinish(t *testing.T) {
	defer goleak.VerifyNone(t)

	var total int32
	err := Finish(func() error {
		atomic.AddInt32(&total, 2)
		return nil
	}, func() error {
		atomic.AddInt32(&total, 3)
		return nil
	}, func() error {
		atomic.AddInt32(&total, 5)
		return nil
	})
	assert.Equal(t, int32(10), atomic.LoadInt32(&total))
	assert.Nil(t, err)
}

func TestFinishNone(t *testing.T) {
	defer goleak.VerifyNone(t)
	assert.Nil(t, Finish())
}

func TestFinishErr(t *testing.T) {
	defer goleak.VerifyNone(t)
	var total int32
	err := Finish(func() error {
		atomic.AddInt32(&total, 2)
		return nil
	}, func() error {
		atomic.AddInt32(&total, 5)
		return errDummy
	}, func() error {
		atomic.AddInt32(&total, 3)
		return nil
	})
	assert.Equal(t, errDummy, err)
}

func TestFinishVoid(t *testing.T) {
	defer goleak.VerifyNone(t)

	var total int32
	FinishVoid(func() {
		atomic.AddInt32(&total, 2)
	}, func() {
		atomic.AddInt32(&total, 3)
	}, func() {
		atomic.AddInt32(&total, 5)
	})
	assert.Equal(t, int32(10), atomic.LoadInt32(&total))
}

func TestForEach(t *testing.T) {
	const tasks = 1000

	t.Run("foreach all", func(t *testing.T) {
		defer goleak.VerifyNone(t)

		var total uint32
		ForEach(func(source chan<- interface{}) {
			for i := 0; i < tasks; i++ {
				source <- i
			}
		}, func(item interface{}) {
			atomic.AddUint32(&total, uint32(1))
		}, WithWorkers(-1))

		assert.Equal(t, tasks, int(total))
	})
	t.Run("panic", func(t *testing.T) {
		defer goleak.VerifyNone(t)
		ForEach(func(source chan<- interface{}) {
			for i := 0; i < tasks; i++ {
				source <- i
			}
		}, func(item interface{}) {
			panic("foo")
		})
	})
}

func TestMapReduce(t *testing.T) {
	// 默认数据生产函数
	generate := func(source chan<- interface{}) {
		for i := 0; i < 5; i++ {
			source <- i
		}
	}
	// 默认数据处理函数
	defaultMapper := func(item interface{}, writer Writer, cancel func(err error)) {
		v := item.(int)
		writer.Write(v * v)
	}
	// 默认数据聚合函数
	defaultReducer := func(pipe <-chan interface{}, writer Writer, cancel func(err error)) {
		var result int
		for item := range pipe {
			result += item.(int)
		}
		writer.Write(result)
	}

	tests := []struct {
		mapper      MapperFunc
		reducer     ReducerFunc
		expectErr   error
		expectValue interface{}
	}{
		{
			mapper:      defaultMapper,
			reducer:     defaultReducer,
			expectErr:   nil,
			expectValue: 30,
		},
		{
			mapper: func(item interface{}, writer Writer, cancel func(error)) {
				v := item.(int)
				if v%3 == 0 {
					cancel(errDummy)
				}
				writer.Write(v)
			},
			reducer:   defaultReducer,
			expectErr: errDummy,
		},
		{
			mapper: func(item interface{}, writer Writer, cancel func(error)) {
				v := item.(int)
				if v%3 == 0 {
					cancel(nil)
				}
				writer.Write(v * v)
			},
			reducer:   defaultReducer,
			expectErr: ErrCancelWithNil,
		},
		{
			mapper: defaultMapper,
			reducer: func(pipe <-chan interface{}, writer Writer, cancel func(err error)) {
				var result int
				for item := range pipe {
					result += item.(int)
					if result > 10 {
						cancel(errDummy)
					}
				}
				writer.Write(result)
			},
			expectErr: errDummy,
		},
	}

	for _, test := range tests {
		t.Run("map reduce test", func(t *testing.T) {
			value, err := MapReduce(generate, test.mapper, test.reducer, WithWorkers(runtime.NumCPU()))
			assert.Equal(t, test.expectErr, err)
			assert.Equal(t, test.expectValue, value)
		})
	}
}

func TestDoVoid(t *testing.T) {

}

func TestWithSource(t *testing.T) {

}

func BenchmarkMapReduce(b *testing.B) {
	for i := 0; i < b.N; i++ {

	}
}
