package model

import "github.com/LannisterAlwaysPaysHisDebts/goLearn/src/blog/simple/conf"

type ArticleRoute struct {
	Item string
	Name string
}

type NewArticles []conf.Articles

func (a NewArticles) Len() int {
	return len(a)
}

func (a NewArticles) Less(i, j int) bool {
	return a[i].Date > a[j].Date
}

func (a NewArticles) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

type HotArticles []conf.Articles

func (h HotArticles) Len() int {
	return len(h)
}

func (h HotArticles) Less(i, j int) bool {
	return h[i].VisitCnt > h[j].VisitCnt
}

func (h HotArticles) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}