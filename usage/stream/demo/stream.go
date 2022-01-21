package demo

import (
	ring2 "github.com/ct-zh/golib/collection/ring"
	"github.com/ct-zh/golib/common"
	"github.com/ct-zh/golib/usage/threading"
	"sort"
	"sync"
)

const (
	defaultWorkers = 16
	minWorkers     = 1
)

type (
	rxOptions struct {
		unlimitedWorkers bool
		workers          int
	}

	Option func(opts *rxOptions)

	GenerateFunc func(source chan<- interface{})

	// KeyFunc defines the method to generate keys for the elements in a Stream
	KeyFunc func(item interface{}) interface{}

	ReduceFunc func(pipe <-chan interface{}) (interface{}, error)

	ForAllFunc func(pipe <-chan interface{})

	ForeachFunc func(item interface{})

	FilterFunc func(item interface{}) bool
	// LessFunc defines  the method to compare the elements in a Stream
	LessFunc func(a, b interface{}) bool

	MapFunc func(item interface{}) interface{}
	// WalkFunc 遍历item写到stream中
	WalkFunc func(item interface{}, pipe chan<- interface{})

	PredicateFunc func(item interface{}) bool
)

type Stream struct {
	source <-chan interface{}
}

func (s Stream) AllMatch(fn PredicateFunc) bool {
	for item := range s.source {
		if !fn(item) {
			go drain(s.source)
			return false
		}
	}
	return true
}

func (s Stream) AnyMatch(fn PredicateFunc) bool {
	for item := range s.source {
		if fn(item) {
			go drain(s.source)
			return true
		}
	}
	return false
}

func (s Stream) Count() int {
	var count int
	for range s.source {
		count++
	}
	return count
}

// Concat 拼接多条stream流
func (s Stream) Concat(streams ...Stream) Stream {
	source := make(chan interface{})

	go func() {
		g := threading.NewRoutineGroup()

		g.Run(func() {
			for item := range s.source {
				source <- item
			}
		})

		for _, itemStream := range streams {
			val := itemStream // important!
			g.Run(func() {
				for item := range val.source {
					source <- item
				}
			})
		}

		g.Wait()
		close(source)
	}()

	return Range(source)
}

// Done 清空stream所有数据
func (s Stream) Done() {
	drain(s.source)
}

// Distinct 去重
func (s Stream) Distinct(fn KeyFunc) Stream {
	source := make(chan interface{})
	threading.GoSafe(func() {
		defer close(source)

		keys := make(map[interface{}]common.PlaceholderType)
		for item := range s.source {
			key := fn(item)
			if _, ok := keys[key]; !ok {
				source <- item
				keys[key] = struct{}{}
			}
		}
	})
	return Range(source)
}

func (s Stream) ForAll(fn ForAllFunc) {
	fn(s.source)
}

// Foreach 循环读取
func (s Stream) Foreach(fn ForeachFunc) {
	for item := range s.source {
		fn(item)
	}
}

// Filter 对元素进行过滤,不符合过滤函数的将被舍弃
func (s Stream) Filter(filterFn FilterFunc, opts ...Option) Stream {
	return s.Walk(func(item interface{}, pipe chan<- interface{}) {
		if filterFn(item) {
			pipe <- item
		}
	}, opts...)
}

// Group 对元素进行分组
func (s Stream) Group(fn KeyFunc) Stream {
	groups := make(map[interface{}][]interface{}) // key是接口，value是接口数组
	for item := range s.source {
		key := fn(item)
		groups[key] = append(groups[key], item)
	}
	source := make(chan interface{})
	go func() {
		for _, group := range groups {
			source <- group
		}
		close(source)
	}()
	return Range(source)
}

// Head 获取前n个元素
func (s Stream) Head(n int64) Stream {
	if n <= 0 {
		panic("n must greater than 1")
	}
	source := make(chan interface{})

	threading.GoSafe(func() {
		for item := range s.source {
			n--
			if n >= 0 {
				source <- item
			}
			if n == 0 { // 情况一，指定n条数据已经全部写完，直接close
				close(source)
			}
		}
		if n > 0 { // 情况二，当前数据已经全部消耗，仍然不够n条，也需要close
			close(source)
		}
	})

	return Range(source)
}

func (s Stream) NoneMatch(fn PredicateFunc) bool {
	for item := range s.source {
		if fn(item) {
			go drain(s.source)
			return false
		}
	}
	return true
}

// Map 对元素做映射转换
func (s Stream) Map(fn MapFunc, opt ...Option) Stream {
	return s.Walk(func(item interface{}, pipe chan<- interface{}) {
		pipe <- fn(item)
	}, opt...)
}

// Reduce 执行对应函数，返回自定义结果
func (s Stream) Reduce(fn ReduceFunc) (interface{}, error) {
	return fn(s.source)
}

func (s Stream) Reverse() Stream {
	var items []interface{}
	for item := range s.source {
		items = append(items, item)
	}
	for i := len(items)/2 - 1; i >= 0; i-- {
		app := len(items) - 1 - i
		items[i], items[app] = items[app], items[i]
	}

	return Just(items...)
}

func (s Stream) Sort(fn LessFunc) Stream {
	var items []interface{}
	for item := range s.source {
		items = append(items, item)
	}

	sort.Slice(items, func(i, j int) bool {
		return fn(i, j)
	})
	return Just(items...)
}

func (s Stream) Tail(n int64) Stream {
	if n < 1 {
		panic("n must be greater than 1")
	}

	source := make(chan interface{})
	go func() {
		ring := ring2.NewRing(int(n))
		for item := range s.source {
			ring.Add(item)
		}
		for _, val := range ring.Take() {
			source <- val
		}
		close(source)
	}()
	return Range(source)
}

// Walk 遍历元素
func (s Stream) Walk(fn WalkFunc, opts ...Option) Stream {
	option := buildOptions(opts...)
	if option.unlimitedWorkers {
		return s.walkUnlimited(fn, option)
	}
	return s.walkLimited(fn, option)
}

// walkLimited 限制协程数量进行遍历;
func (s Stream) walkLimited(fn WalkFunc, option *rxOptions) Stream {
	// 如何控制并发数量？ 不无限制的开协程 => mini协程池
	pipe := make(chan interface{}, option.workers)
	go func() {
		var wg sync.WaitGroup
		pool := make(chan common.PlaceholderType, option.workers)

		for item := range s.source {
			val := item
			pool <- common.Placeholder
			wg.Add(1)
			threading.GoSafe(func() {
				defer func() {
					wg.Done()
					<-pool
				}()
				fn(val, pipe)
			})
		}
		wg.Wait()
		close(pipe)
	}()
	return Range(pipe)
}

// walkUnlimited 未限制协程数量进行遍历
func (s Stream) walkUnlimited(fn WalkFunc, option *rxOptions) Stream {
	pipe := make(chan interface{}, option.workers)
	go func() {
		var wg sync.WaitGroup
		for item := range s.source {
			// important! 并发陷阱
			val := item
			wg.Add(1)
			threading.GoSafe(func() {
				defer wg.Done()
				fn(val, pipe)
			})
		}
		wg.Wait()
		close(pipe)
	}()
	return Range(pipe)
}

func Range(source <-chan interface{}) Stream {
	return Stream{source: source}
}

func Just(items ...interface{}) Stream {
	source := make(chan interface{}, len(items))
	for _, item := range items {
		source <- item
	}
	close(source)
	return Range(source)
}

func From(generate GenerateFunc) Stream {
	source := make(chan interface{})
	threading.GoSafe(func() {
		defer close(source)
		generate(source)
	})
	return Range(source)
}

func WithWorkers(workers int) Option {
	return func(opts *rxOptions) {
		if workers < minWorkers {
			opts.workers = minWorkers
		} else {
			opts.workers = workers
		}
	}
}

func drain(channel <-chan interface{}) {
	for range channel {
	}
}

func buildOptions(opts ...Option) *rxOptions {
	options := newOptions()
	for _, opt := range opts {
		opt(options)
	}
	return options
}

func newOptions() *rxOptions {
	return &rxOptions{workers: defaultWorkers}
}
