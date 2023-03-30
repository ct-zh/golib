package tls

import (
	"context"
	"sync"
)

type (
	contextKey struct{}
	tlsStorage = sync.Map
)

var (
	contextKeyActive = contextKey{}
	tlsStore         []tlsStorage
)

const (
	maxShardsCount = 127
)

// GoID get an uniq id, not as the same goroutine id
func GoID() uint64 {
	return uint64(uintptr(G()))
}

func setLocal(goid uint64, local map[interface{}]interface{}) {
	shardIndex := goid % maxShardsCount
	tlsStore[shardIndex].Store(goid, local)
}

func deleteLocal(goid uint64) {
	shardIndex := goid % maxShardsCount
	tlsStore[shardIndex].Delete(goid)
}

func For(ctx context.Context, f func()) func() {
	return func() {
		local := make(map[interface{}]interface{})
		local[contextKeyActive] = ctx
		goid := GoID()
		setLocal(goid, local)
		defer deleteLocal(goid)
		f()
	}
}
