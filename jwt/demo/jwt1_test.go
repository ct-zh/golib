package demo

import (
	"strconv"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

const SECRET string = "fb697a356a4363170061f0ae9da5828b"

func Test(t *testing.T) {
	Convey("jwt hs256 test", t, func() {
		Convey("token测试", func() {
			uid := 10

			// 创建token
			token, err := CreateToken(strconv.Itoa(uid), SECRET, time.Now().Add(time.Minute*15).Unix())
			So(err, ShouldBeNil)
			So(token, ShouldNotBeEmpty)

			t.Logf("CreateToken: %s \n", token)

			// 验证token
			res, err := ParseToken(token, SECRET)
			So(err, ShouldBeNil)
			So(res, ShouldNotBeEmpty)
			So(res, ShouldEqual, strconv.Itoa(uid))

			t.Logf("parse result: %+v \n", res)

			// 同一个uid创建第二个token
			token2, err := CreateToken(strconv.Itoa(uid), SECRET, time.Now().Add(time.Minute*15).Unix())
			So(err, ShouldBeNil)
			So(token2, ShouldNotBeEmpty)

			t.Logf("CreateToken2: %s \n", token2)

			// 用第一个token解析
			res1, err := ParseToken(token, SECRET)
			t.Logf("parse token1 again: result: %+v err: %+v \n", res1, err)

			res3, err := ParseToken(token, SECRET)
			So(err, ShouldBeNil)
			So(res3, ShouldNotBeEmpty)
			So(res3, ShouldEqual, strconv.Itoa(uid))

			t.Logf("parse token2: result: %+v \n", res3)

		})

		Convey("token过期测试", func() {
			uid := 11

			token, err := CreateToken(strconv.Itoa(uid), SECRET, time.Now().Add(time.Second*15).Unix())

			So(err, ShouldBeNil)
			So(token, ShouldNotBeEmpty)

			t.Logf("CreateToken: %s \n", token)

			// 验证token
			res, err := ParseToken(token, SECRET)
			So(err, ShouldBeNil)
			So(res, ShouldNotBeEmpty)
			So(res, ShouldEqual, strconv.Itoa(uid))

			t.Logf("parse result: %+v \n", res)
			t.Log("waiting for token to expire...")

			// 验证过期token
			time.Sleep(time.Second * 16)
			res2, err := ParseToken(token, SECRET)
			So(err, ShouldNotBeNil)
			t.Logf("parse result: %+v err: %+v \n", res2, err)
		})
	})
}
