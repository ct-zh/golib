package simple

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSyncSet_Add(t *testing.T) {
	Convey("add", t, func() {
		s := SyncSet{}
		key := "aaaa"

		s.Add(key)
		So(s.Has(key), ShouldBeTrue)
		So(s.Has("bbbb"), ShouldBeFalse)
	})
}

func TestSyncSet_Delete(t *testing.T) {
	Convey("delete", t, func() {
		set := SyncSet{}
		set.Add("aaaa")
		So(set.Has("aaaa"), ShouldBeTrue)
		set.Delete("aaaa")
		So(set.Has("aaaa"), ShouldBeFalse)
	})
}

// 两个set的基准测试
func BenchmarkSet(b *testing.B) {
	b.Run("set 性能测试", func(b *testing.B) {
		s := NewSet()
		for i := 0; i < b.N; i++ {
			s.Add("aaa")
			if s.Has("aaa") {
				s.Delete("aaa")
			}
		}
	})
	b.Run("sync set 性能测试", func(b *testing.B) {
		s := SyncSet{}
		for i := 0; i < b.N; i++ {
			s.Add("aaa")
			if s.Has("aaa") {
				s.Delete("aaa")
			}
		}
	})
}
