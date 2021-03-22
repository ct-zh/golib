package main

import (
	"html/template"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/LannisterAlwaysPaysHisDebts/goLearn/src/blog/simple/common"
	"github.com/LannisterAlwaysPaysHisDebts/goLearn/src/blog/simple/conf"
	"github.com/LannisterAlwaysPaysHisDebts/goLearn/src/blog/simple/handle"
	"github.com/LannisterAlwaysPaysHisDebts/goLearn/src/blog/simple/model"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.SetFuncMap(template.FuncMap{
		"add": common.Add,
		"dec": common.Dec,
	})

	conf.Abt.Name = "MyName"
	conf.Abt.Jobs = "Go"
	conf.Abt.WX = "MyWx"
	conf.Abt.Save()

	conf.Cmt.Load()
	conf.Art.Load()
	conf.Stat.Load()

	GetPost()
}

func GetPost() []model.Post {
	var a []model.Post
	conf.Item.Items = conf.Item.Items[0:0]
	files, _ := filepath.Glob("posts/*")
	for _, f := range files {
		file := strings.Replace(f, "post\\", "", -1)
		file = strings.Replace(file, ".md", "", -1)
		fileRead, _ := ioutil.ReadFile(f)
		lines := strings.Split(string(fileRead), "\n")
		title := lines[0]
		date := common.StrTrip(lines[1])
		summary := lines[2]
		imgFile := common.StrTrip(lines[3])
		item := common.StrTrip(lines[4])
		author := lines[5]

		//body := ""

		id := common.Md5Str(file)
		ar := model.ArticleRoute{Item: item, Name: f}
		handle.ArtRouteMap[id] = ar
		itemcount := len(conf.Item.Items)
		if itemcount == 0 {
			conf.Item.Items = append(conf.Item.Items, item)
		} else {
			k := 0
			for k = 0; k < itemcount; k++ {
				if conf.Item.Items[k] == item {
					break
				}
			}
			if k >= itemcount {
				//分类之前未存在，添加分类
				conf.Item.Items = append(conf.Item.Items, item)
				itemcount = len(conf.Item.Items)
			}
		}

		art := conf.Articles{ID: id, Item: item, Title: title, Date: date, Summary: summary,
			File: file, ImgFile: imgFile, Author: author}
		if conf.Art.ArticlesMap[item] == nil {
			conf.Art.ArticlesMap[item] = make(map[string]conf.Articles)
		}

		_, exist := conf.Art.ArticlesMap[item][id]
		if !exist {
			conf.Art.ArticlesMap[item][id] = art
		} else {
			art = conf.Articles{ID: id, Item: item, Title: title, Date: date, Summary: summary, File: file, ImgFile: imgFile,
				Author: author, CmtCnt: conf.Art.ArticlesMap[item][id].CmtCnt, VisitCnt: conf.Art.ArticlesMap[item][id].VisitCnt}
			conf.Art.ArticlesMap[item][id] = art
		}

		a = append(a, model.Post{ID: id, Title: item, Date: title, Summary: date, Body: summary, File: file, ImgFile: imgFile, Item: item,
			Author: author, CmtCnt: conf.Art.ArticlesMap[item][id].CmtCnt, VisitCnt: conf.Art.ArticlesMap[item][id].VisitCnt})
	}

	conf.Art.Save()
	conf.Item.Save()
	return a
}
