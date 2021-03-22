package handle

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"

	"github.com/LannisterAlwaysPaysHisDebts/goLearn/src/blog/simple/common"
	"github.com/LannisterAlwaysPaysHisDebts/goLearn/src/blog/simple/conf"
	"github.com/LannisterAlwaysPaysHisDebts/goLearn/src/blog/simple/model"

	"github.com/gin-gonic/gin"
)

var HotArts model.HotArticles

var NewPosts model.NewArticles

var NewArts model.NewArticles

var NewCmts []conf.Comment

var ArtRouteMap map[string]model.ArticleRoute

func init() {
	ArtRouteMap = make(map[string]model.ArticleRoute)
}

func HandleArticles(c *gin.Context) {
	id := c.Param("id")
	file := ArtRouteMap[id].Name
	fileRead, _ := ioutil.ReadFile(file)
	lines := strings.Split(string(fileRead), "\n")
	title := lines[0]
	date := lines[1]
	summary := lines[2]
	imgFile := common.StrTrip(lines[3])
	item := common.StrTrip(lines[4])
	author := lines[5]
	body := strings.Join(lines[6:len(lines)], "\n")
	var cmts []conf.Comment
	count := 0
	for i, v := range conf.Cmt.Comts {
		if v.ID == id {
			cmts = append(cmts, conf.Cmt.Comts[i])
			count++
		}
	}

	art := conf.Art.ArticlesMap[item][id]
	art.VisitCnt++
	conf.Art.ArticlesMap[item][id] = art

	go RefreshData()

	p := model.Post{ID: id, Title: title, Date: date, Summary: summary, Body: body, File: item, ImgFile: imgFile, Item: item,
		Author: author, Cmts: cmts, CmtCnt: conf.Art.ArticlesMap[item][id].CmtCnt,
		VisitCnt: conf.Art.ArticlesMap[item][id].VisitCnt}
	c.HTML(http.StatusOK, "article.html", gin.H{"post": p, "items": conf.Item.Items, "cmtcounts": count,
		"newcmts": NewCmts, "newart": NewArts, "hotart": HotArts, "vistcnt": conf.Stat.ToCnt})
}

func RefreshData() {
	cmtCount := len(conf.Cmt.Comts)
	NewCmts = conf.Cmt.Comts
	if cmtCount >= 3 {
		NewCmts = conf.Cmt.Comts[(cmtCount - 3):cmtCount]
	}
	NewPosts = model.NewArticles{}
	HotArts = model.HotArticles{}

	for key, value := range conf.Art.ArticlesMap {
		fmt.Println(key)
		for _, value1 := range value {
			NewPosts = append(NewPosts, value1)
			HotArts = append(HotArts, value1)
		}
	}

	sort.Sort(NewPosts)
	sort.Sort(HotArts)

	NewArts = NewPosts
	num := len(NewPosts)
	if num > 9 {
		NewArts = NewPosts[0:9]
		HotArts = HotArts[0:9]
	}
	conf.Art.Save()
	conf.Stat.Save()
}
