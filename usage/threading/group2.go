package threading

import (
	"context"
	"fmt"
	"runtime"
	"sync"
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

}
