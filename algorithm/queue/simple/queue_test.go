package simple

import (
	"fmt"
	"sync"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func tQueue(q Queue) {
	So(q.IsEmpty(), ShouldBeTrue)
	data := []int{1, 2, 3, 4, 5}
	for _, datum := range data {
		q.RPush(datum)
	}
	So(q.IsEmpty(), ShouldBeFalse)

	for _, datum := range data {
		So(q.LPop().(int), ShouldEqual, datum)
	}

	So(q.IsEmpty(), ShouldBeTrue)
	for _, datum := range data {
		q.LPush(datum)
	}
	So(q.IsEmpty(), ShouldBeFalse)

	for _, datum := range data {
		So(q.RPop().(int), ShouldEqual, datum)
	}
	So(q.IsEmpty(), ShouldBeTrue)
}

func TestQueue(t *testing.T) {
	Convey("queue", t, func() {
		Convey("simple/simpleQueue", func() {
			q := New()
			tQueue(q)
		})
		Convey("syncQueue", func() {
			q2 := NewSyncQueue()
			tQueue(q2)
		})
	})
}

// 并发测试; list是线程不安全的
func TestParallel(t *testing.T) {
	//q := New()
	q := NewSyncQueue()

	// 总共20 * 500 = 10000 次push
	wg := sync.WaitGroup{}
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for i := 0; i < 500; i++ {
				q.RPush(fmt.Sprintf("%d:%d", id, i))
			}
		}(i)
	}

	wg.Wait()
	t.Logf("len: %d", q.Len()) // 结果len小于10000次
}
