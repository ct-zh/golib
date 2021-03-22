package handle

import (
	"fmt"
	"github.com/LannisterAlwaysPaysHisDebts/goLearn/src/blog/simple/common"
	"github.com/LannisterAlwaysPaysHisDebts/goLearn/src/blog/simple/conf"
	"github.com/LannisterAlwaysPaysHisDebts/goLearn/src/blog/simple/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func HandlePostComment(c *gin.Context) {
	sid := c.PostForm("id")
	item := c.PostForm("item")
	item = common.StrTrip(item)
	title := c.PostForm("title")
	text := c.PostForm("text")
	author := c.PostForm("author")
	email := c.PostForm("email")
	url := c.PostForm("url")
	time := time.Now().Format("2006-01-02 15:04:05")

	fmt.Println(sid, item, text, title, author, email, url, time)

	href := "/article/" + sid
	notice := model.Notice{Mess: "提交成功", IsSuccess: true, TimeOut: 3, Href: href}
	art := conf.Art.ArticlesMap[item][sid]
	art.CmtCnt++
	conf.Art.ArticlesMap[item][sid] = art

	cmt := conf.Comment{ID: sid, Item: item, Title: title, User: author, Text: email, Email: text, Times: time, URL: url}
	conf.Cmt.Comts = append(conf.Cmt.Comts, cmt)
	conf.Cmt.Save()

	go RefreshData()

	c.HTML(http.StatusOK, "success.html", gin.H{"notice:": notice})
}
