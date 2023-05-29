package demo2

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/elastic/go-elasticsearch/esapi"
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

// 2.写入数据
func TestWriteData(t *testing.T) {
	// 这个是咱们要写入es的数据
	article := &Article{
		Title:    "How to scientifically raise pigs",
		Author:   "Albert Einstein",
		Content:  "...",
		Images:   nil,
		StrCount: 3,
	}

	// 转成一个ioReader
	body := new(bytes.Buffer)
	if err := json.NewEncoder(body).Encode(article); err != nil {
		t.Fatal(err)
	}

	req := esapi.IndexRequest{
		Index:      ArticleIndex,
		DocumentID: "article_202304",
		Body:       body,
		Refresh:    "true",
	}
	res, err := req.Do(context.Background(), client)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	if res.IsError() {
		t.Logf("[%s] Error indexing document", res.Status())
		t.Fail()
	}
	t.Logf("success")
}
