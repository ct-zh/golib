package simple

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func Test(t *testing.T) {
	Convey("tree test", t, func() {
		// 创建一个tree
		root := NewNode(5)
		root.Left = NewNode(3)
		root.Right = NewNode(10)
		root.Left.Right = NewNode(4)
		root.Right.Left = NewNode(8)

		So(root.GetVal(), ShouldEqual, 5)
		root.SetVal(6)
		So(root.GetVal(), ShouldEqual, 6)

		maxInt := 0
		root.TraverseFn(func(node2 *node) {
			if node2.GetVal() > maxInt {
				maxInt = node2.GetVal()
			}
		})
		So(maxInt, ShouldEqual, 10)

		minInt := 100
		ch := root.TraverseWithChannel()
		for v := range ch {
			if v.GetVal() < minInt {
				minInt = v.GetVal()
			}
		}
		So(minInt, ShouldEqual, 3)

	})
}
