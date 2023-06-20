package gpool

import "sync"

type poolstatus int

const (
	Run poolstatus = iota + 1
	Stop
)

type Handler struct {
	Fn   func(...interface{})
	Args []interface{}
}

type Pool2 struct {
	cap     uint
	workers uint
	status  poolstatus

	ch chan *Handler
	mu sync.Mutex
}

func (p *Pool2) run() {
	p.workers++

	defer func() {
		p.workers--
	}()

	select {
	case handler, ok := <-p.ch:
		if !ok {
			return
		}
		handler.Fn(handler.Args)
	}
}

func (p *Pool2) Put(handler *Handler) {
	p.ch <- handler
}
