package demo

import (
	"io/ioutil"
	"log"
	"runtime"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStream_AllMatch(t *testing.T) {
	runCheckedTest(t, func(t *testing.T) {
		result := Just(1, 2, 3, 4, 5).AllMatch(func(item interface{}) bool {
			return item.(int) < 10
		})
		assert.True(t, result)

		result2 := Just(1, 2, 3, 4, 5).AllMatch(func(item interface{}) bool {
			return item.(int)%2 == 0
		})
		assert.False(t, result2)
	})
}

func TestStream_AnyMatch(t *testing.T) {
	runCheckedTest(t, func(t *testing.T) {
		result1 := Just(1, 2, 3).AnyMatch(func(item interface{}) bool {
			return item.(int) == 2
		})
		assert.True(t, result1)
		result2 := Just(1, 2, 3).AnyMatch(func(item interface{}) bool {
			return item.(int) == 5
		})
		assert.False(t, result2)
	})
}

func TestStream_Count(t *testing.T) {
	runCheckedTest(t, func(t *testing.T) {
		tests := []struct {
			name   string
			data   []interface{}
			result int
		}{
			{
				name:   "empty elements with nil",
				result: 0,
			},
			{
				name:   "empty elements",
				data:   []interface{}{},
				result: 0,
			},
			{
				name:   "some elements",
				data:   []interface{}{1, 2, 3},
				result: 3,
			},
			{
				name:   "1 elements",
				data:   []interface{}{1},
				result: 1,
			},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				result := Just(test.data...).Count()
				assert.Equal(t, test.result, result)
			})
		}
	})
}

func TestStream_Concat(t *testing.T) {
	runCheckedTest(t, func(t *testing.T) {
		var result []int
		s0 := Just(0, 0, 0, 0)
		s1 := Just(1, 3, 5, 7)
		Just(2, 4, 6, 8).Concat(s0, s1).Foreach(func(item interface{}) {
			result = append(result, item.(int))
		})
		t.Logf("%+v", result)
	})
}

func TestStream_Distinct(t *testing.T) {
	runCheckedTest(t, func(t *testing.T) {
		var result int
		_, err := Just(4, 1, 3, 2, 3, 4).Distinct(func(item interface{}) interface{} {
			return item
		}).Reduce(func(pipe <-chan interface{}) (interface{}, error) {
			for item := range pipe {
				result += item.(int)
			}
			return result, nil
		})
		assert.Nil(t, err)
		assert.Equal(t, 10, result)
	})

	runCheckedTest(t, func(t *testing.T) {
		Just(1, 2, 3, 3, 4, 4, 5, 5).Distinct(func(item interface{}) interface{} {
			uid := item.(int)
			if uid > 3 {
				return 4
			}
			return item
		}).Foreach(func(item interface{}) {
			t.Logf("%+v ", item)
		})
	})
}

func TestStream_ForAll(t *testing.T) {
	runCheckedTest(t, func(t *testing.T) {
		var result []int
		Just(1, 2, 3).ForAll(func(pipe <-chan interface{}) {
			for item := range pipe {
				result = append(result, item.(int))
			}
		})
		assert.Equal(t, []int{1, 2, 3}, result)
	})
}

func TestStream_Filter(t *testing.T) {
	runCheckedTest(t, func(t *testing.T) {
		result, err := Just(1, 2, 3, 4).Filter(func(item interface{}) bool {
			return item.(int)%2 == 0
		}).Reduce(func(pipe <-chan interface{}) (interface{}, error) {
			result := 0
			for item := range pipe {
				result += item.(int)
			}
			return result, nil
		})
		assert.Nil(t, err)
		assert.Equal(t, 6, result)
	})
}

func TestStream_Group(t *testing.T) {
	runCheckedTest(t, func(t *testing.T) {
		var result [][]interface{}
		Just(1, 2, 3, 4, 5, 6).Group(func(item interface{}) interface{} {
			if item.(int)%2 == 0 {
				return "singular"
			} else {
				return "complex"
			}
		}).Foreach(func(item interface{}) {
			result = append(result, item.([]interface{}))
		})
		assert.Equal(t, 2, len(result))
		assert.Equal(t, [][]interface{}{
			{1, 3, 5}, {2, 4, 6},
		}, result)
	})
}

func TestStream_Head(t *testing.T) {
	runCheckedTest(t, func(t *testing.T) {
		tests := []struct {
			name   string
			head   int64
			data   []interface{}
			answer int
		}{
			{
				name:   "multiple elements",
				data:   []interface{}{1, 2, 3, 4, 5},
				head:   2,
				answer: 3,
			},
		}

		for _, test := range tests {
			var result int
			Just(test.data...).Head(test.head).Foreach(func(item interface{}) {
				result += item.(int)
			})
			assert.Equal(t, test.answer, result)
		}
	})
}

func TestStream_Map(t *testing.T) {
	runCheckedTest(t, func(t *testing.T) {
		log.SetOutput(ioutil.Discard)

		tests := []struct {
			mapper MapFunc
			expect int
		}{
			{
				mapper: func(item interface{}) interface{} {
					return item.(int) * item.(int)
				},
				expect: 30,
			},
			{
				mapper: func(item interface{}) interface{} {
					v := item.(int)
					if v%2 == 0 {
						return 0
					}
					return v * v
				},
				expect: 10,
			},
			{
				mapper: func(item interface{}) interface{} {
					v := item.(int)
					if v%2 == 0 {
						panic(v)
					}
					return v * v
				},
				expect: 10,
			},
		}

		for i, test := range tests {
			name := "test " + strconv.Itoa(i)
			t.Run(name, func(t *testing.T) {
				var result int
				var workers int
				if i%2 == 0 {
					workers = 0
				} else {
					workers = runtime.NumCPU()
				}

				_, err := From(func(source chan<- interface{}) {
					for i := 0; i < 5; i++ {
						source <- i
					}
				}).Map(test.mapper, WithWorkers(workers)).Reduce(func(pipe <-chan interface{}) (interface{}, error) {
					for val := range pipe {
						result += val.(int)

					}
					return result, nil
				})
				assert.Nil(t, err)
				assert.Equal(t, test.expect, result)
			})
		}
	})
}

func runCheckedTest(t *testing.T, fn func(t *testing.T)) {
	goroutines := runtime.NumGoroutine()
	fn(t)
	// let scheduler schedule first
	time.Sleep(time.Millisecond)
	assert.True(t, runtime.NumGoroutine() <= goroutines)
}
