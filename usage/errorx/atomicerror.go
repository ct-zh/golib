package errorx

import "sync/atomic"

type AtomicError struct {
	err atomic.Value
}

func (e *AtomicError) Set(err error) {
	if err != nil {
		e.err.Store(err)
	}
}

func (e *AtomicError) Load() error {
	if err := e.err.Load(); err != nil {
		return err.(error)
	}
	return nil
}

func Chain(fn ...func() error) error {
	for _, f := range fn {
		if err := f(); err != nil {
			return err
		}
	}
	return nil
}
