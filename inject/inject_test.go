package inject

import (
	"reflect"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type arg1 interface{}

func add(a int, b arg1) int {
	return a + b.(int)
}

func TestNew(t *testing.T) {
	Convey("inject包New函数测试: 注入函数", t, func() {
		inj := New()
		inj.Map(3)
		inj.MapTo(4, (*arg1)(nil))
		res, err := inj.Invoke(add)
		So(err, ShouldBeNil)
		So(len(res), ShouldEqual, 1)
		So(res[0].Int(), ShouldEqual, 7)
	})
}

type User struct {
	Name string `inject`
	Age  int    `inject`
}

func TestApply(t *testing.T) {
	Convey("inject包Apply函数测试: 注入struct", t, func() {
		u1 := User{}

		var name1 string = "李元芳"
		var age1 int = 18

		inj := New()
		inj.Map(age1)
		inj.Map(name1)
		inj.Apply(&u1)
		So(u1.Name, ShouldEqual, name1)
		So(u1.Age, ShouldEqual, age1)
	})
}

func TestSetParent(t *testing.T) {
	Convey("setParent ", t, func() {
		inj1 := New()
		inj1.Map("刘德华")
		inj1.Map(18)

		So(inj1.Get(reflect.TypeOf("a")).IsValid(), ShouldBeTrue)
		So(inj1.Get(reflect.TypeOf(1)).IsValid(), ShouldBeTrue)

		So(inj1.Get(reflect.TypeOf([]byte("string"))).IsValid(), ShouldBeFalse)

		inj2 := New()
		inj2.Map([]byte("test"))
		inj1.SetParent(inj2)
		So(inj1.Get(reflect.TypeOf([]byte("string"))).IsValid(), ShouldBeTrue)
	})
}
