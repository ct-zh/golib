package simple

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSet_Add(t *testing.T) {
	Convey("add", t, func() {
		set := NewSet()
		set.Add("aaaa")
		So(set.Has("aaaa"), ShouldBeTrue)
		So(set.Has("bbbb"), ShouldBeFalse)
	})
}

func TestSet_Delete(t *testing.T) {
	Convey("delete", t, func() {
		set := NewSet()
		set.Add("aaaa")
		So(set.Has("aaaa"), ShouldBeTrue)
		set.Delete("aaaa")
		So(set.Has("aaaa"), ShouldBeFalse)
	})
}
