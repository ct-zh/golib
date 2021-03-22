package conf

import "log"

type Articles struct {
	ID       string
	Item     string
	Title    string
	Date     string
	Summary  string
	File     string
	ImgFile  string // 文章前的图片展示
	Author   string
	CmtCnt   int //评论数量
	VisitCnt int //浏览数量
}

type ArtCfg struct {
	BaseConf
	Ver         int
	ArticlesMap map[string]map[string]Articles
}

var Art ArtCfg

func init() {
	log.Println("ArtCfg init ENTER")
	Art.BaseConf.conf = &Art
	Art.ArticlesMap = make(map[string]map[string]Articles)
}
