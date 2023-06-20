package mapreduce

import (
	"context"

	"github.com/ct-zh/golib/common"
)

type guardedWriter struct {
	ctx     context.Context
	channel chan<- interface{}
	done    <-chan common.PlaceholderType
}

func newGuardedWriter(ctx context.Context, channel chan<- interface{}, done <-chan common.PlaceholderType) *guardedWriter {
	return &guardedWriter{ctx: ctx, channel: channel, done: done}
}

func (w guardedWriter) Write(v interface{}) {
	select {
	case <-w.ctx.Done():
		return
	case <-w.done:
		return
	default:
		w.channel <- v
	}
}
