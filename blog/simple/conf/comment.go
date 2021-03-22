package conf

import "log"

type Comment struct {
	ID string
	Item string
	Title string
	User string
	Text string
	Email string
	Times string
	URL string
}

type ComtCfg struct {
	BaseConf
	Comts []Comment
}

var Cmt ComtCfg

func init()  {
	log.Println("ComtCfg init ENTER")
	Cmt.BaseConf.conf = &Cmt
	Cmt.Comts = make([]Comment, 0)
}