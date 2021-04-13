package simple

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"sync"
	"testing"
)

func TestQueue(t *testing.T) {
	Convey("simple/simpleQueue", t, func() {
		q := New()

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
	})
}

// 并发测试; list是线程不安全的
func TestParallel(t *testing.T) {
	q := New()

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
