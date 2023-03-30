package threading

import (
	"context"
	"fmt"
	"runtime"
	"sync"

	"github.com/ct-zh/golib/usage/tls"
)

// A Group is a collection of goroutines working on subtasks that are part of the same overall task.
// Group是处理子任务的goroutine的集合，这些子任务是同一总体任务的一部分。

// go function will run in tls.For with panic recover and report.
// go函数将在tls.For中运行，并进行紧急恢复和报告。

// group can use GOMAXPROCS function to set max goroutine to work.
// Group可以使用GOMAXPROCS函数将最大goroutine设置为工作。

type Group struct {
	err     error
	wg      sync.WaitGroup
	errOnce sync.Once

	workerOnce sync.Once
	ch         chan call
	chs        []call
}

type call struct {
	ctx context.Context
	f   func(ctx context.Context) error
}

func (g *Group) do(ctx context.Context, f func(ctx2 context.Context) error) {
	var err error
	defer func() {
		if r := recover(); r != nil {
			buf := make([]byte, 64<<10)
			buf = buf[:runtime.Stack(buf, false)]
			err = fmt.Errorf("errgroup: panic recovered: %s\n%s", r, buf)
			// error log
			//logging.For(ctx).Errorf("service panic: %s", err)
			// crash log
			//logging.CrashLog(err)
			// stdout log
			//log.Printf("service panic: %s\n", err)
			// stderr log
			//_, _ = fmt.Fprintf(os.Stderr, "service panic: %s\n", err)
			// metrics 监控
			//metrics.Meter("", 1, "name", "go.group")
		}
		if err != nil {
			g.errOnce.Do(func() {
				g.err = err
			})
		}
	}()
	tls.For(ctx, func() {
		err = f(ctx)
	})
}

// GOMAXPROCS set max goroutine to work.
// GOMAXPROCS设置最大goroutine为工作。
func (g *Group) GOMAXPROCS(n int) *Group {
	if n <= 0 {
		n = 1
	}
	g.workerOnce.Do(func() {
		g.ch = make(chan call, n)
		for i := 0; i < n; i++ {
			go func() {
				for c := range g.ch {
					g.do(c.ctx, c.f)
				}
			}()
		}
	})
	return g
}

// Go calls the given function in a new goroutine within tls.For.
// Go在tls.For中的新goroutine中调用给定的函数。
// Go will recover if any panic occurs.
// Go将恢复任何紧急情况。
// The _ctx will passed to f.
// The First error will be returned by Wait.
func (g *Group) Go(ctx context.Context, f func(ctx context.Context) error) {
	g.wg.Add(1)
	if g.ch != nil {
		select {
		case g.ch <- call{ctx, f}:
		default:

		}
	}
	go g.do(ctx, f)
}

// Wait blocks until all function calls from the Go method have returned, then returns the first non-nil error (if any) from them.
// 直到Go方法的所有函数调用都返回，然后从中返回第一个非零错误（如果有）
func (g *Group) Wait() error {
	if g.ch != nil {
		for _, f := range g.chs {
			g.ch <- f
		}
	}
	g.wg.Wait()
	if g.ch != nil {
		close(g.ch) // let all receiver exit
	}
	return g.err
}
