package simple

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestQueue(t *testing.T) {
	Convey("simple/queue", t, func() {
		q := New()

		So(q.IsEmpty(), ShouldBeTrue)
		data := []int{1, 2, 3, 4, 5}
		for _, datum := range data {
			q.Push(datum)
		}
		So(q.IsEmpty(), ShouldBeFalse)

		for _, datum := range data {
			So(q.Pop().(int), ShouldEqual, datum)
		}
	})
}
