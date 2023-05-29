package demo2

import (
	"bytes"
	"context"
	"encoding/json"
	jsoniter "github.com/json-iterator/go"
	"testing"

	"github.com/elastic/go-elasticsearch/esapi"
	. "github.com/smartystreets/goconvey/convey"
)

// es基本操作

// ## 写入数据

// 定义index名称，相当于指定数据库名，这个可以写在配置文件里面
const ArticleIndex = "user_article"

// 1.定义一个结构体
type Article struct {
	Title    string   `json:"title,omitempty"`
	Author   string   `json:"author,omitempty"`
	Content  string   `json:"content,omitempty"`
	Images   []string `json:"images,omitempty"`
	StrCount int64    `json:"str_count,omitempty"`
}

// 2.增、删、查
func TestWriteData(t *testing.T) {
	// 这个是咱们要写入es的数据
	article := &Article{
		Title:    "How to scientifically raise pigs",
		Author:   "Albert Einstein",
		Content:  "...",
		Images:   nil,
		StrCount: 3,
	}

	Convey("TestMyData", t, func() {
		Convey("TestWriteData", func() {
			// 转成一个ioReader
			body := new(bytes.Buffer)
			err := json.NewEncoder(body).Encode(article)
			So(err, ShouldBeNil)

			req := esapi.IndexRequest{
				Index:      ArticleIndex,
				DocumentID: "article_202304",
				Body:       body,
				Refresh:    "true",
			}

			res, err := req.Do(context.Background(), client)
			So(err, ShouldBeNil)

			defer res.Body.Close()
			So(res.IsError(), ShouldBeFalse)

			result := map[string]interface{}{}
			err = jsoniter.NewDecoder(res.Body).Decode(&result)
			So(err, ShouldBeNil)

			Printf("res = %+v", result)
		})

		Convey("TestReadData", func() {
			buf := new(bytes.Buffer)
			query := map[string]interface{}{
				"query": map[string]interface{}{
					"term": map[string]interface{}{
						"title": "How to scientifically raise pigs",
					},
				},
			}
			encoder := jsoniter.NewEncoder(buf)
			err := encoder.Encode(&query)
			So(err, ShouldBeNil)

			res, err := client.Search(
				client.Search.WithContext(context.Background()),
				client.Search.WithIndex(ArticleIndex),
				client.Search.WithBody(buf),
				client.Search.WithTrackTotalHits(true),
				client.Search.WithPretty(),
			)
			So(err, ShouldBeNil)
			defer res.Body.Close()

			So(res.IsError(), ShouldBeFalse)

			r := new(Article)
			err = jsoniter.NewDecoder(res.Body).Decode(&r)
			So(err, ShouldBeNil)

			So(r.Title, ShouldEqual, "How to scientifically raise pigs")
			So(r, ShouldEqual, article)
		})

	})

}
