package ring

import (
	"sync"
)

type Ring struct {
	elements []interface{}
	index    int
	lock     sync.Mutex
}

func (r *Ring) Add(ele interface{}) {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.elements[r.index%len(r.elements)] = ele
	r.index++
}

func (r *Ring) Take() []interface{} {
	r.lock.Lock()
	defer r.lock.Unlock()

	var start int
	var size int
	if r.index < len(r.elements) { // 当前未超过环的大小，从0开始读取
		size = r.index
	} else { // 当前超过环的大小，已覆盖一部分数据，需要从新的位置开始读取
		size = len(r.elements)
		// 比如环大小为4 写了1，2，3，4，5，index=5
		// 实际数据为 5, 2, 3, 4 应该从1开始读取 index % len
		start = r.index % len(r.elements)
	}

	elements := make([]interface{}, size)
	for i := 0; i < size; i++ {
		elements[i] = r.elements[(start+i)%len(r.elements)]
	}
	return elements
}

func NewRing(cap int) *Ring {
	if cap < 1 {
		panic("ring capacity must greater than 1")
	}
	return &Ring{
		elements: make([]interface{}, cap),
	}
}
