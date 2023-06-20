package gpool

import (
	"errors"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

var ErrInvalidPoolCap = errors.New("invalid pool cap")
var ErrPoolAlreadyClosed = errors.New("pool already closed")

type poolStatus int64

const (
	RUNNING poolStatus = iota + 1
	STOP
)

type Task struct {
	Handler func(v ...interface{})
	Params  []interface{}
}

type Pool struct {
	capacity       uint64     // 容量
	runningWorkers uint64     // 当前运行的 worker（goroutine）数量
	status         poolStatus // 运行中或已关闭, 用于安全关闭任务池
	chTask         chan *Task // 任务队列
	sync.Mutex

	PanicHandler func(interface{})
}

func (p *Pool) SetStatus(status poolStatus) {
	p.Lock()
	defer p.Unlock()

	p.status = status
}

func (p *Pool) Capacity() uint64 {
	return p.capacity
}

func NewPool(capacity uint64) (*Pool, error) {
	if capacity == 0 {
		return nil, ErrInvalidPoolCap
	}
	return &Pool{
		capacity: capacity,
		status:   RUNNING,
		chTask:   make(chan *Task),
	}, nil
}

func (p *Pool) checkWorker() {
	p.Lock()
	defer p.Unlock()

	if p.runningWorkers == 0 && len(p.chTask) > 0 {
		p.run()
	}
}

func (p *Pool) run() {
	// 运行中的任务加一
	p.incRunning()

	go func() {
		defer func() { // worker 结束, 运行中的任务减一
			p.decRunning()
			if r := recover(); r != nil {
				if p.PanicHandler != nil {
					p.PanicHandler(r)
				} else {
					log.Printf("Worker panic: %s \n", r)
				}
			}
			p.checkWorker()
		}()

		select {
		case task, ok := <-p.chTask:
			if !ok {
				return
			}
			task.Handler(task.Params...)
		}
	}()
}

func (p *Pool) Put(task *Task) error {
	p.Lock()
	defer p.Unlock()

	if p.status == STOP {
		return ErrPoolAlreadyClosed
	}

	if p.GetRunningWorkers() < p.Capacity() {
		p.run()
	}

	if p.status == RUNNING {
		p.chTask <- task
	}
	return nil
}

func (p *Pool) Close() {
	p.SetStatus(STOP)
	for len(p.chTask) > 0 {
		time.Sleep(1e6)
	}

	close(p.chTask)
}

func (p *Pool) incRunning() {
	atomic.AddUint64(&p.runningWorkers, 1)
}

func (p *Pool) decRunning() {
	atomic.AddUint64(&p.runningWorkers, ^uint64(0))
}

func (p *Pool) GetRunningWorkers() uint64 {
	return atomic.LoadUint64(&p.runningWorkers)
}
