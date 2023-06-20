package mapreduce

// mapreduce

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"

	"github.com/ct-zh/golib/common"
	"github.com/ct-zh/golib/usage/errorx"
)

const (
	defaultWorkers = 16
	minWorkers     = 1
)

var (
	ErrCancelWithNil  = errors.New("mapreduce cancelled with nil")
	ErrReduceNoOutput = errors.New("reduce not write value")
)

type (
	// GenerateFunc 数据生产函数 is used to let callers send elements into source
	GenerateFunc func(source chan<- interface{})

	// MapFunc is used to do element processing and write the output to writer
	MapFunc func(item interface{}, writer Writer)

	// ForEachFunc is used to do element processing, but no output
	ForEachFunc func(item interface{})

	// MapperFunc 数据加工
	MapperFunc func(item interface{}, writer Writer, cancel func(error))

	// ReducerFunc 数据聚合
	ReducerFunc func(pipe <-chan interface{}, writer Writer, cancel func(err error))

	VoidReducerFunc func(pipe <-chan interface{}, cancel func(err error))

	Option func(opts *mrOptions)

	mapperContext struct {
		ctx       context.Context
		mapper    MapFunc
		source    <-chan interface{}
		panicChan *onceChan
		collector chan<- interface{}
		doneChan  <-chan common.PlaceholderType
		workers   int
	}

	onceChan struct {
		channel chan interface{}
		wrote   int32
	}

	mrOptions struct {
		ctx     context.Context
		workers int
	}

	// Writer 任意提供Write的结构都可以作为Writer
	Writer interface {
		Write(v interface{})
	}
)

// Finish 并发执行 func,出现错误则抛出
func Finish(fns ...func() error) error {
	if len(fns) == 0 {
		return nil
	}

	return DoVoid(func(source chan<- interface{}) {
		for _, fn := range fns {
			source <- fn
		}
	}, func(item interface{}, writer Writer, cancel func(error)) {
		fn := item.(func() error)
		if err := fn(); err != nil {
			cancel(err)
		}
	}, func(pipe <-chan interface{}, cancel func(err error)) {
	}, WithWorkers(len(fns)))
}

// FinishVoid runs functions parallel without error output
// 并发执行函数，出现错误也不终止流程
func FinishVoid(fns ...func()) {
	if len(fns) == 0 {
		return
	}
	ForEach(func(source chan<- interface{}) {
		for _, fn := range fns {
			source <- fn
		}
	}, func(item interface{}) {
		fn, ok := item.(func())
		if ok {
			fn()
		}
	}, WithWorkers(len(fns)))
}

// ForEach maps all elements from given generate func but no output.
func ForEach(generate GenerateFunc, mapper ForEachFunc, opts ...Option) {
	options := buildOptions(opts...)
	panicChan := &onceChan{channel: make(chan interface{})}
	source := buildSource(generate, panicChan)
	collector := make(chan interface{}, options.workers)
	done := make(chan common.PlaceholderType)

	go executeMappers(mapperContext{
		ctx: options.ctx,
		mapper: func(item interface{}, writer Writer) {
			mapper(item)
		},
		source:    source,
		panicChan: panicChan,
		collector: collector,
		doneChan:  done,
		workers:   options.workers,
	})

	for {
		select {
		case v := <-panicChan.channel:
			panic(v)
		case _, ok := <-collector:
			if !ok {
				return
			}
		}
	}
}

// MapReduce maps all elements generated from give generate func,
// and reduces the output elements with given reducer
func MapReduce(generateFunc GenerateFunc,
	mapperFunc MapperFunc, reducerFunc ReducerFunc,
	option ...Option) (interface{}, error) {
	panicChan := &onceChan{channel: make(chan interface{})}
	source := buildSource(generateFunc, panicChan)
	return mapReduceWithPanicChan(source, panicChan, mapperFunc, reducerFunc, option...)
}

// DoVoid 聚合处理
func DoVoid(generateFunc GenerateFunc,
	mapperFunc MapperFunc, reducerFunc VoidReducerFunc, opts ...Option) error {
	_, err := MapReduce(generateFunc, mapperFunc, func(pipe <-chan interface{}, writer Writer, cancel func(err error)) {
		reducerFunc(pipe, cancel)
	}, opts...)
	if errors.Is(err, ErrReduceNoOutput) {
		return nil
	}
	return err
}

// WithSource 直接从channel读取数据
func mapReduceWithPanicChan(source <-chan interface{}, panicChan *onceChan,
	mapper MapperFunc, reducer ReducerFunc, opts ...Option) (interface{}, error) {
	options := buildOptions(opts...)
	// output is used to write the final result
	output := make(chan interface{})
	defer func() {
		// reducer can only write once, if more, panic
		for range output {
			panic("more than one element written in reducer")
		}
	}()

	// collector is used to collect data from mapper, and consume in reducer
	collector := make(chan interface{}, options.workers)
	// if done is closed, all mappers and reducer should stop processing
	done := make(chan common.PlaceholderType)
	writer := newGuardedWriter(options.ctx, output, done)
	var closeOnce sync.Once
	// use atomic.Value to avoid data race
	var retErr errorx.AtomicError
	finish := func() {
		closeOnce.Do(func() {
			close(done)
			close(output)
		})
	}
	cancel := once(func(err error) {
		if err != nil {
			retErr.Set(err)
		} else {
			retErr.Set(ErrCancelWithNil)
		}

		drain(source)
		finish()
	})

	go func() {
		defer func() {
			drain(collector)
			if r := recover(); r != nil {
				panicChan.write(r)
			}
			finish()
		}()

		reducer(collector, writer, cancel)
	}()

	go executeMappers(mapperContext{
		ctx: options.ctx,
		mapper: func(item interface{}, w Writer) {
			mapper(item, w, cancel)
		},
		source:    source,
		panicChan: panicChan,
		collector: collector,
		doneChan:  done,
		workers:   options.workers,
	})

	select {
	case <-options.ctx.Done():
		cancel(context.DeadlineExceeded)
		return nil, context.DeadlineExceeded
	case v := <-panicChan.channel:
		// drain output here, otherwise for loop panic in defer
		drain(output)
		panic(v)
	case v, ok := <-output:
		if err := retErr.Load(); err != nil {
			return nil, err
		} else if ok {
			return v, nil
		} else {
			return nil, ErrReduceNoOutput
		}
	}
}

func WithWorkers(workers int) Option {
	return func(opts *mrOptions) {
		if workers < minWorkers {
			opts.workers = minWorkers
		} else {
			opts.workers = workers
		}
	}
}

func buildOptions(option ...Option) *mrOptions {
	mrOpt := newMrOptions()
	for _, opt := range option {
		opt(mrOpt)
	}
	return mrOpt
}

func buildSource(generateFunc GenerateFunc, panicChan *onceChan) chan interface{} {
	source := make(chan interface{})
	go func() {
		defer func() {
			if r := recover(); r != nil {
				panicChan.write(r)
			}
			close(source)
		}()

		generateFunc(source)
	}()
	return source
}

func drain(ch <-chan interface{}) {
	for range ch {
	}
}

func executeMappers(mCtx mapperContext) {
	var wg sync.WaitGroup
	defer func() {
		wg.Wait()
		close(mCtx.collector)
		drain(mCtx.source)
	}()

	var failed int32
	pool := make(chan common.PlaceholderType, mCtx.workers)
	writer := newGuardedWriter(mCtx.ctx, mCtx.collector, mCtx.doneChan)
	for atomic.LoadInt32(&failed) == 0 {
		select {
		case <-mCtx.ctx.Done():
			return
		case <-mCtx.doneChan:
			return
		case pool <- common.PlaceholderType{}:
			item, ok := <-mCtx.source
			if !ok {
				<-pool
				return
			}
			wg.Add(1)
			go func() {
				defer func() {
					if r := recover(); r != nil {
						atomic.AddInt32(&failed, 1)
						mCtx.panicChan.write(r)
					}
					wg.Done()
					<-pool
				}()

				mCtx.mapper(item, writer)
			}()
		}
	}
}

func once(fn func(error)) func(err error) {
	once := new(sync.Once)
	return func(err error) {
		once.Do(func() {
			fn(err)
		})
	}
}

func newMrOptions() *mrOptions {
	return &mrOptions{
		ctx:     context.Background(),
		workers: defaultWorkers,
	}
}

func (o *onceChan) write(val interface{}) {
	if atomic.AddInt32(&o.wrote, 1) > 1 {
		return
	}
	o.channel <- val
}
