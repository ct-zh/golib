package model

import (
	"github.com/LannisterAlwaysPaysHisDebts/goLearn/src/blog/simple/conf"
)

type Post struct {
	ID      string
	Title   string
	Date    string
	Summary string
	Body    string
	File    string
	ImgFile string
	Item    string
	Author  string

	Cmts     []conf.Comment
	CmtCnt   int
	VisitCnt int
}
